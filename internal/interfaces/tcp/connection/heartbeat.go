package connection

import (
	"context"
	"fmt"
	"sync"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// HeartbeatManager 心跳管理器
type HeartbeatManager struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
	logger   logging.Logger
	interval time.Duration
	timeout  time.Duration
}

// NewHeartbeatManager 创建心跳管理器
func NewHeartbeatManager(logger logging.Logger, interval, timeout time.Duration) *HeartbeatManager {
	return &HeartbeatManager{
		sessions: make(map[string]*Session),
		logger:   logger,
		interval: interval,
		timeout:  timeout,
	}
}

// AddSession 添加会话
func (hm *HeartbeatManager) AddSession(session *Session) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	hm.sessions[session.ID] = session
	hm.logger.Info("会话已添加到心跳管理", map[string]interface{}{
		"session_id": session.ID,
	})
}

// RemoveSession 移除会话
func (hm *HeartbeatManager) RemoveSession(sessionID string) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	delete(hm.sessions, sessionID)
	hm.logger.Info("会话已从心跳管理移除", map[string]interface{}{
		"session_id": sessionID,
	})
}

// Start 启动心跳检查
func (hm *HeartbeatManager) Start(ctx context.Context) {
	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	hm.logger.Info("心跳管理器启动", map[string]interface{}{
		"interval": hm.interval,
		"timeout":  hm.timeout,
	})

	for {
		select {
		case <-ctx.Done():
			hm.logger.Info("心跳管理器停止")
			return
		case <-ticker.C:
			hm.checkHeartbeats()
		}
	}
}

// checkHeartbeats 检查心跳
func (hm *HeartbeatManager) checkHeartbeats() {
	hm.mutex.RLock()
	sessions := make([]*Session, 0, len(hm.sessions))
	for _, session := range hm.sessions {
		sessions = append(sessions, session)
	}
	hm.mutex.RUnlock()

	now := time.Now()
	timeoutCount := 0

	for _, session := range sessions {
		if now.Sub(session.LastActivity) > hm.timeout {
			hm.logger.Info("会话心跳超时", map[string]interface{}{
				"session_id":    session.ID,
				"last_activity": session.LastActivity,
				"timeout":       hm.timeout,
			})

			// 关闭超时会话
			session.Close()
			hm.RemoveSession(session.ID)
			timeoutCount++
		}
	}

	if timeoutCount > 0 {
		hm.logger.Info("心跳检查完成", map[string]interface{}{
			"timeout_count": timeoutCount,
			"active_count":  len(sessions) - timeoutCount,
		})
	}
}

// SendHeartbeat 发送心跳
func (hm *HeartbeatManager) SendHeartbeat(sessionID string) error {
	hm.mutex.RLock()
	session, exists := hm.sessions[sessionID]
	hm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("会话不存在: %s", sessionID)
	}

	// 更新最后活跃时间
	session.LastActivity = time.Now()

	hm.logger.Info("心跳已发送", map[string]interface{}{
		"session_id": sessionID,
	})

	return nil
}

// GetActiveSessionCount 获取活跃会话数量
func (hm *HeartbeatManager) GetActiveSessionCount() int {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	return len(hm.sessions)
}
