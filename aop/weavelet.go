package aop

import (
	"fmt"

	"github.com/google/uuid"
	"greatestworks/aop/protos"
)

// CheckWeaveletInfo checks that weavelet information is well-formed.
func CheckWeaveletInfo(w *protos.WeaveletInfo) error {
	if w == nil {
		return fmt.Errorf("WeaveletInfo: nil")
	}
	if w.App == "" {
		return fmt.Errorf("WeaveletInfo: missing app name")
	}
	if _, err := uuid.Parse(w.DeploymentId); err != nil {
		return fmt.Errorf("WeaveletInfo: invalid deployment id: %w", err)
	}
	if w.Group == nil {
		return fmt.Errorf("WeaveletInfo: nil colocation group")
	}
	if w.Group.Name == "" {
		return fmt.Errorf("WeaveletInfo: missing colocation group name")
	}
	if w.GroupId == "" {
		return fmt.Errorf("WeaveletInfo: missing colocation group replica id")
	}
	if _, err := uuid.Parse(w.Id); err != nil {
		return fmt.Errorf("WeaveletInfo: invalid weavelet id: %w", err)
	}
	if w.SingleProcess && !w.SingleMachine {
		return fmt.Errorf("WeaveletInfo: single process but not single machine")
	}
	return nil
}
