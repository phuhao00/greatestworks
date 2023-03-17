package aop

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"greatestworks/aop/protos"
)

// CheckDeployment checks that a deployment is well-formed.
func CheckDeployment(d *protos.Deployment) error {
	if d == nil {
		return fmt.Errorf("nil deployment")
	}
	appName := d.App.Name
	if appName == "" {
		appName = filepath.Base(d.App.Binary)
	}
	if _, err := uuid.Parse(d.Id); err != nil {
		return fmt.Errorf("invalid deployment id for %s: %w", appName, err)
	}
	return nil
}

// DeploymentID returns the identifier that uniquely identifies a particular deployment.
func DeploymentID(d *protos.Deployment) uuid.UUID {
	id, err := uuid.Parse(d.Id)
	if err != nil {
		panic(fmt.Sprintf("bad UUID in internal proto: %v", err))
	}
	return id
}
