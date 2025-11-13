package simclient

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"greatestworks/internal/infrastructure/logging"
	tcpProtocol "greatestworks/internal/interfaces/tcp/protocol"
)

// E2EScenario 端到端场景：登录→移动→技能释放→登出，验证 AOI 广播
type E2EScenario struct {
	cfg    ScenarioConfig
	logger logging.Logger
}

// NewE2EScenario 创建端到端测试场景
func NewE2EScenario(cfg ScenarioConfig, logger logging.Logger) *E2EScenario {
	if cfg.ActionInterval.AsDuration() <= 0 {
		cfg.ActionInterval = NewDuration(1 * time.Second)
	}
	return &E2EScenario{cfg: cfg, logger: logger}
}

// Name 返回场景名称
func (s *E2EScenario) Name() string {
	if s.cfg.Name != "" {
		return s.cfg.Name
	}
	return "e2e-workflow"
}

// Execute 执行端到端场景
func (s *E2EScenario) Execute(ctx context.Context, client *SimulatorClient) (*ScenarioResult, error) {
	result := &ScenarioResult{ScenarioName: s.Name(), StartedAt: time.Now()}
	defer func() {
		result.CompletedAt = time.Now()
	}()

	// 步骤1：认证（如果启用）
	token, err := authenticateAndRecord(ctx, client, result, s.logger)
	if err != nil && s.cfg.StopOnError {
		return result, err
	}
	_ = token // 可选：后续可用于附加到 TCP 消息

	// 步骤2：连接网关
	start := time.Now()
	conn, err := client.ConnectGateway(ctx)
	connectLatency := time.Since(start)
	fields := map[string]interface{}{}
	if err == nil && conn != nil {
		fields["remote"] = conn.RemoteAddr().String()
	}
	result.Record("gateway.connect", connectLatency, err, fields)
	if err != nil {
		return result, err
	}
	defer conn.Close()

	// 步骤3：发送 TCP 登录包（PlayerLogin 消息）
	if err := s.sendLogin(result, client, conn); err != nil {
		if s.cfg.StopOnError {
			return result, err
		}
	}

	// 短暂等待登录响应
	time.Sleep(100 * time.Millisecond)
	s.tryReadResponse(result, client, conn, "login.response")

	// 步骤4：发送移动包（模拟移动到新位置）
	if err := s.sendMove(result, client, conn, 100.0, 50.0, 10.0); err != nil {
		if s.cfg.StopOnError {
			return result, err
		}
	}

	time.Sleep(100 * time.Millisecond)
	s.tryReadResponse(result, client, conn, "move.response")

	// 步骤5：发送技能释放包
	if err := s.sendSkillCast(result, client, conn, 1001, 2001); err != nil {
		if s.cfg.StopOnError {
			return result, err
		}
	}

	time.Sleep(100 * time.Millisecond)
	s.tryReadResponse(result, client, conn, "skill.response")

	// 步骤6：再次移动（验证多次操作）
	if err := s.sendMove(result, client, conn, 120.0, 55.0, 10.0); err != nil {
		if s.cfg.StopOnError {
			return result, err
		}
	}

	time.Sleep(100 * time.Millisecond)
	s.tryReadResponse(result, client, conn, "move2.response")

	// 步骤7：发送登出包
	if err := s.sendLogout(result, client, conn); err != nil {
		if s.cfg.StopOnError {
			return result, err
		}
	}

	time.Sleep(100 * time.Millisecond)
	s.tryReadResponse(result, client, conn, "logout.response")

	return result, nil
}

// sendLogin 发送 PlayerLogin 消息（带 JSON payload）
func (s *E2EScenario) sendLogin(result *ScenarioResult, client *SimulatorClient, conn net.Conn) error {
	payload := map[string]interface{}{
		"player_id": fmt.Sprintf("%d", client.PlayerID()),
		"map_id":    1,
	}
	return s.sendMessageWithPayload(result, client, conn, "gateway.msg.login", tcpProtocol.MsgPlayerLogin, payload)
}

// sendMove 发送 PlayerMove 消息
func (s *E2EScenario) sendMove(result *ScenarioResult, client *SimulatorClient, conn net.Conn, x, y, z float64) error {
	payload := map[string]interface{}{
		"position": map[string]interface{}{
			"x": x,
			"y": y,
			"z": z,
		},
	}
	return s.sendMessageWithPayload(result, client, conn, "gateway.msg.move", tcpProtocol.MsgPlayerMove, payload)
}

// sendSkillCast 发送技能释放消息
func (s *E2EScenario) sendSkillCast(result *ScenarioResult, client *SimulatorClient, conn net.Conn, skillID, targetID int32) error {
	payload := map[string]interface{}{
		"skill_id":  skillID,
		"target_id": targetID,
	}
	return s.sendMessageWithPayload(result, client, conn, "gateway.msg.skill", tcpProtocol.MsgBattleSkill, payload)
}

// sendLogout 发送登出消息
func (s *E2EScenario) sendLogout(result *ScenarioResult, client *SimulatorClient, conn net.Conn) error {
	return s.sendMessageWithPayload(result, client, conn, "gateway.msg.logout", tcpProtocol.MsgPlayerLogout, nil)
}

// sendMessageWithPayload 发送带 JSON payload 的消息
func (s *E2EScenario) sendMessageWithPayload(
	result *ScenarioResult,
	client *SimulatorClient,
	conn net.Conn,
	action string,
	messageType uint32,
	payload interface{},
) error {
	start := time.Now()

	var payloadBytes []byte
	var err error
	if payload != nil {
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			result.Record(action, time.Since(start), fmt.Errorf("marshal payload: %w", err), nil)
			return err
		}
	}

	messageID := nextMessageID()
	seq := client.seq + 1
	client.seq = seq

	// 构造完整消息：Header + Payload
	header := make([]byte, tcpProtocol.MessageHeaderSize)
	binary.BigEndian.PutUint32(header[0:], tcpProtocol.MessageMagic)
	binary.BigEndian.PutUint32(header[4:], messageID)
	binary.BigEndian.PutUint32(header[8:], messageType)
	binary.BigEndian.PutUint16(header[12:], tcpProtocol.FlagRequest)
	binary.BigEndian.PutUint64(header[14:], client.PlayerID())
	binary.BigEndian.PutUint64(header[22:], uint64(time.Now().Unix()))
	binary.BigEndian.PutUint16(header[30:], uint16(seq&0xFFFF))

	// 发送 header + payload
	if err := conn.SetWriteDeadline(time.Now().Add(client.cfg.Gateway.WriteTimeout.AsDuration())); err != nil {
		s.logger.Warn("failed to set write deadline", logging.Fields{"error": err})
	}

	if _, err := conn.Write(header); err != nil {
		result.Record(action, time.Since(start), fmt.Errorf("write header: %w", err), nil)
		return err
	}

	if len(payloadBytes) > 0 {
		if _, err := conn.Write(payloadBytes); err != nil {
			result.Record(action, time.Since(start), fmt.Errorf("write payload: %w", err), nil)
			return err
		}
	}

	fields := map[string]interface{}{
		"message_type":   messageType,
		"message_id":     messageID,
		"payload_length": len(payloadBytes),
	}
	result.Record(action, time.Since(start), nil, fields)
	return nil
}

// tryReadResponse 尝试读取响应（不阻塞太久）
func (s *E2EScenario) tryReadResponse(result *ScenarioResult, client *SimulatorClient, conn net.Conn, action string) {
	timeout := 200 * time.Millisecond
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return
	}

	buf := make([]byte, 4096)
	start := time.Now()
	n, err := conn.Read(buf)
	duration := time.Since(start)

	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			// 超时不记录为错误（可能服务器不立即响应）
			result.Record(action, duration, nil, map[string]interface{}{"timeout": true})
			return
		}
		if errors.Is(err, io.EOF) {
			result.Record(action, duration, err, map[string]interface{}{"eof": true})
			return
		}
		result.Record(action, duration, err, nil)
		return
	}

	result.Record(action, duration, nil, map[string]interface{}{"bytes": n})
}
