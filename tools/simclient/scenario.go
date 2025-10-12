package simclient

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"greatestworks/internal/infrastructure/logging"
	tcpProtocol "greatestworks/internal/interfaces/tcp/protocol"
)

// Scenario represents an executable workflow for a simulated player.
type Scenario interface {
	Name() string
	Execute(ctx context.Context, client *SimulatorClient) (*ScenarioResult, error)
}

// ScenarioResult captures per-action latencies and errors for analysis.
type ScenarioResult struct {
	ScenarioName string
	StartedAt    time.Time
	CompletedAt  time.Time
	Actions      []ActionRecord
}

// ActionRecord stores the outcome of a single simulated step.
type ActionRecord struct {
	Name     string
	Duration time.Duration
	Err      error
	Data     map[string]interface{}
}

// Record appends an action result to the scenario.
func (r *ScenarioResult) Record(name string, duration time.Duration, err error, fields map[string]interface{}) {
	r.Actions = append(r.Actions, ActionRecord{
		Name:     name,
		Duration: duration,
		Err:      err,
		Data:     fields,
	})
}

// Success reports whether every recorded action completed successfully.
func (r *ScenarioResult) Success() bool {
	for _, action := range r.Actions {
		if action.Err != nil {
			return false
		}
	}
	return true
}

// Errors returns a flattened slice of non-nil action errors.
func (r *ScenarioResult) Errors() []error {
	var errs []error
	for _, action := range r.Actions {
		if action.Err != nil {
			errs = append(errs, action.Err)
		}
	}
	return errs
}

func authenticateAndRecord(ctx context.Context, client *SimulatorClient, result *ScenarioResult, logger logging.Logger) (string, error) {
	start := time.Now()
	token, err := client.Login(ctx)
	result.Record("auth.login", time.Since(start), err, nil)
	if err == nil && token != "" {
		logger.Debug("authentication token acquired", logging.Fields{"token_length": len(token)})
	}
	return token, err
}

// BasicScenario drives a minimal end-to-end flow against the gateway.
type BasicScenario struct {
	cfg    ScenarioConfig
	logger logging.Logger
}

// NewBasicScenario creates the default scenario implementation.
func NewBasicScenario(cfg ScenarioConfig, logger logging.Logger) *BasicScenario {
	if cfg.ActionInterval.AsDuration() <= 0 {
		cfg.ActionInterval = NewDuration(1 * time.Second)
	}
	if cfg.Duration.AsDuration() <= 0 {
		cfg.Duration = NewDuration(10 * time.Second)
	}
	return &BasicScenario{cfg: cfg, logger: logger}
}

// Name returns the configured scenario name.
func (s *BasicScenario) Name() string {
	return s.cfg.Name
}

// Execute runs the scenario for the supplied client.
func (s *BasicScenario) Execute(ctx context.Context, client *SimulatorClient) (*ScenarioResult, error) {
	result := &ScenarioResult{ScenarioName: s.cfg.Name, StartedAt: time.Now()}
	defer func() {
		result.CompletedAt = time.Now()
	}()

	if _, err := authenticateAndRecord(ctx, client, result, s.logger); err != nil {
		if s.cfg.StopOnError {
			return result, err
		}
		s.logger.Warn("authentication failed but continuing", logging.Fields{"error": err})
	}

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

	// Login handshake.
	if err := s.sendMessage(result, client, conn, "gateway.msg.login", tcpProtocol.MsgPlayerLogin); err != nil {
		if s.cfg.StopOnError {
			return result, err
		}
	}

	duration := s.cfg.Duration.AsDuration()
	interval := s.cfg.ActionInterval.AsDuration()
	if interval <= 0 {
		interval = time.Second
	}

	deadline := time.Now().Add(duration)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for iteration := 0; ; iteration++ {
		if duration > 0 && time.Now().After(deadline) {
			break
		}

		select {
		case <-ctx.Done():
			err := ctx.Err()
			result.Record("scenario.cancelled", 0, err, nil)
			return result, err
		case <-ticker.C:
			if err := s.sendMessage(result, client, conn, "gateway.msg.move", tcpProtocol.MsgPlayerMove); err != nil && s.cfg.StopOnError {
				return result, err
			}
			if err := s.sendMessage(result, client, conn, "gateway.msg.heartbeat", tcpProtocol.MsgHeartbeat); err != nil && s.cfg.StopOnError {
				return result, err
			}

			if ok, readErr := client.TryRead(conn); readErr != nil {
				if errors.Is(readErr, io.EOF) {
					result.Record("gateway.connection.closed", 0, readErr, nil)
					return result, readErr
				}
				result.Record("gateway.read", 0, readErr, nil)
				if s.cfg.StopOnError {
					return result, readErr
				}
			} else if ok {
				result.Record("gateway.read", 0, nil, map[string]interface{}{"bytes": tcpProtocol.MessageHeaderSize})
			}
		}

		if duration <= 0 && iteration >= 1 {
			break
		}
	}

	if err := s.sendMessage(result, client, conn, "gateway.msg.logout", tcpProtocol.MsgPlayerLogout); err != nil && s.cfg.StopOnError {
		return result, err
	}

	return result, nil
}

func (s *BasicScenario) sendMessage(result *ScenarioResult, client *SimulatorClient, conn net.Conn, action string, messageType uint32) error {
	start := time.Now()
	messageID, err := client.SendGatewayMessage(conn, messageType, tcpProtocol.FlagRequest)
	fields := map[string]interface{}{
		"message_type": messageType,
		"message_id":   messageID,
	}
	result.Record(action, time.Since(start), err, fields)
	return err
}
