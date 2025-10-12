package simclient

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"greatestworks/internal/infrastructure/logging"
	tcpProtocol "greatestworks/internal/interfaces/tcp/protocol"
)

var messageIDCounter uint32

func init() {
	rand.Seed(time.Now().UnixNano())
	messageIDCounter = uint32(rand.Int31())
}

func nextMessageID() uint32 {
	return atomic.AddUint32(&messageIDCounter, 1)
}

// SimulatorClient coordinates HTTP authentication and TCP gateway traffic for a single virtual player.
type SimulatorClient struct {
	id         int
	cfg        *Config
	logger     logging.Logger
	httpClient *http.Client
	playerName string
	playerID   uint64
	seq        uint32
}

// NewSimulatorClient constructs a simulator client with per-player logging context.
func NewSimulatorClient(id int, cfg *Config, baseLogger logging.Logger) *SimulatorClient {
	cfg.Normalize()

	httpTimeout := cfg.Auth.Timeout.AsDuration()
	if httpTimeout <= 0 {
		httpTimeout = 5 * time.Second
	}

	name := fmt.Sprintf("%s_%06d", cfg.Scenario.PlayerPrefix, id)
	playerID := hashToUint64(name)

	logger := baseLogger.WithFields(logging.Fields{
		"client_id":    id,
		"player":       name,
		"player_id":    playerID,
		"scenario":     cfg.Scenario.Name,
		"gateway":      fmt.Sprintf("%s:%d", cfg.Gateway.Host, cfg.Gateway.Port),
		"auth_enabled": cfg.Auth.Enabled,
	})

	return &SimulatorClient{
		id:         id,
		cfg:        cfg,
		logger:     logger,
		httpClient: &http.Client{Timeout: httpTimeout},
		playerName: name,
		playerID:   playerID,
	}
}

// PlayerName returns the human readable identifier used for hashing.
func (c *SimulatorClient) PlayerName() string {
	return c.playerName
}

// PlayerID returns the hashed numeric player identifier encoded in gateway headers.
func (c *SimulatorClient) PlayerID() uint64 {
	return c.playerID
}

// Login authenticates the player via the auth service and returns a bearer token.
func (c *SimulatorClient) Login(ctx context.Context) (string, error) {
	if !c.cfg.Auth.Enabled {
		c.logger.Debug("auth skipped because auth.enabled=false")
		return "", nil
	}

	payload := map[string]string{
		"username": c.cfg.Auth.Username,
		"password": c.cfg.Auth.Password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal login payload: %w", err)
	}

	base := strings.TrimRight(c.cfg.Auth.BaseURL, "/")
	path := c.cfg.Auth.LoginPath
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	endpoint := base + path

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	latency := time.Since(start)
	if err != nil {
		return "", fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return "", fmt.Errorf("auth request returned %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("decode auth response: %w", err)
	}

	if loginResp.Token == "" {
		return "", fmt.Errorf("auth response missing token")
	}

	c.logger.Info("authenticated successfully", logging.Fields{
		"latency_ms": latency.Milliseconds(),
	})

	return loginResp.Token, nil
}

// ConnectGateway establishes a TCP connection to the gateway service.
func (c *SimulatorClient) ConnectGateway(ctx context.Context) (net.Conn, error) {
	address := fmt.Sprintf("%s:%d", c.cfg.Gateway.Host, c.cfg.Gateway.Port)
	dialer := &net.Dialer{Timeout: c.cfg.Gateway.ConnectTimeout.AsDuration()}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("dial gateway %s: %w", address, err)
	}

	if timeout := c.cfg.Gateway.ReadTimeout.AsDuration(); timeout > 0 {
		_ = conn.SetReadDeadline(time.Now().Add(timeout))
	}
	if timeout := c.cfg.Gateway.WriteTimeout.AsDuration(); timeout > 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(timeout))
	}

	c.logger.Info("connected to gateway", logging.Fields{
		"remote": conn.RemoteAddr().String(),
	})
	return conn, nil
}

// SendGatewayMessage writes a framed header-only message to the gateway.
func (c *SimulatorClient) SendGatewayMessage(conn net.Conn, msgType uint32, flags uint16) (uint32, error) {
	messageID := nextMessageID()
	seq := atomic.AddUint32(&c.seq, 1)
	header := buildHeader(messageID, msgType, flags, c.playerID, time.Now().Unix(), seq)

	if err := conn.SetWriteDeadline(time.Now().Add(c.cfg.Gateway.WriteTimeout.AsDuration())); err != nil {
		c.logger.Warn("failed to set write deadline", logging.Fields{"error": err})
	}

	if _, err := conn.Write(header); err != nil {
		return 0, fmt.Errorf("write message to gateway: %w", err)
	}

	return messageID, nil
}

// TryRead attempts to consume a header-sized response, ignoring timeouts to avoid blocking.
func (c *SimulatorClient) TryRead(conn net.Conn) (bool, error) {
	timeout := c.cfg.Gateway.ReadTimeout.AsDuration()
	if timeout <= 0 {
		timeout = 250 * time.Millisecond
	}

	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return false, err
	}

	buf := make([]byte, tcpProtocol.MessageHeaderSize)
	n, err := conn.Read(buf)
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			return false, nil
		}
		if err == io.EOF {
			return false, io.EOF
		}
		return false, fmt.Errorf("read gateway response: %w", err)
	}
	if n == 0 {
		return false, nil
	}
	return true, nil
}

func buildHeader(messageID, messageType uint32, flags uint16, playerID uint64, timestamp int64, sequence uint32) []byte {
	header := make([]byte, tcpProtocol.MessageHeaderSize)

	binary.BigEndian.PutUint32(header[0:], tcpProtocol.MessageMagic)
	binary.BigEndian.PutUint32(header[4:], messageID)
	binary.BigEndian.PutUint32(header[8:], messageType)
	binary.BigEndian.PutUint16(header[12:], flags)
	binary.BigEndian.PutUint64(header[14:], playerID)
	binary.BigEndian.PutUint64(header[22:], uint64(timestamp))
	// Sequence is defined as uint32 but wire format currently uses two bytes; keep lower 16 bits for compatibility.
	binary.BigEndian.PutUint16(header[30:], uint16(sequence&0xFFFF))
	return header
}

func hashToUint64(value string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(value))
	return h.Sum64()
}
