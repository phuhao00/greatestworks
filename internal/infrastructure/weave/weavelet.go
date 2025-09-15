// Package weave provides Service Weaver integration functionality
package weave

import (
	"fmt"

	"github.com/google/uuid"
)

// WeaveletInfo represents weavelet information for Service Weaver
type WeaveletInfo struct {
	App          string            `json:"app"`
	DeploymentId string            `json:"deployment_id"`
	Group        *ColocationGroup  `json:"group"`
	GroupId      string            `json:"group_id"`
	Id           string            `json:"id"`
	SingleProcess bool             `json:"single_process"`
	SingleMachine bool             `json:"single_machine"`
	Labels       map[string]string `json:"labels,omitempty"`
}

// ColocationGroup represents a colocation group in Service Weaver
type ColocationGroup struct {
	Name       string   `json:"name"`
	Components []string `json:"components"`
}

// WeaveletConfig holds configuration for weavelet integration
type WeaveletConfig struct {
	Enabled       bool              `yaml:"enabled" json:"enabled"`
	App           string            `yaml:"app" json:"app"`
	DeploymentId  string            `yaml:"deployment_id" json:"deployment_id"`
	GroupName     string            `yaml:"group_name" json:"group_name"`
	Components    []string          `yaml:"components" json:"components"`
	SingleProcess bool              `yaml:"single_process" json:"single_process"`
	SingleMachine bool              `yaml:"single_machine" json:"single_machine"`
	Labels        map[string]string `yaml:"labels" json:"labels"`
}

// DefaultWeaveletConfig returns default weavelet configuration
func DefaultWeaveletConfig() WeaveletConfig {
	return WeaveletConfig{
		Enabled:       false,
		App:           "greatestworks",
		DeploymentId:  uuid.New().String(),
		GroupName:     "main",
		Components:    []string{"player", "battle", "inventory"},
		SingleProcess: true,
		SingleMachine: true,
		Labels:        make(map[string]string),
	}
}

// WeaveletManager manages weavelet lifecycle and operations
type WeaveletManager struct {
	config *WeaveletConfig
	info   *WeaveletInfo
}

// NewWeaveletManager creates a new weavelet manager
func NewWeaveletManager(config *WeaveletConfig) *WeaveletManager {
	if config == nil {
		defaultConfig := DefaultWeaveletConfig()
		config = &defaultConfig
	}
	
	return &WeaveletManager{
		config: config,
	}
}

// Initialize initializes the weavelet manager
func (wm *WeaveletManager) Initialize() error {
	if !wm.config.Enabled {
		return nil
	}
	
	// Create weavelet info
	wm.info = &WeaveletInfo{
		App:          wm.config.App,
		DeploymentId: wm.config.DeploymentId,
		Group: &ColocationGroup{
			Name:       wm.config.GroupName,
			Components: wm.config.Components,
		},
		GroupId:      uuid.New().String(),
		Id:           uuid.New().String(),
		SingleProcess: wm.config.SingleProcess,
		SingleMachine: wm.config.SingleMachine,
		Labels:       wm.config.Labels,
	}
	
	// Validate weavelet info
	if err := CheckWeaveletInfo(wm.info); err != nil {
		return fmt.Errorf("weavelet validation failed: %w", err)
	}
	
	return nil
}

// GetInfo returns the weavelet information
func (wm *WeaveletManager) GetInfo() *WeaveletInfo {
	return wm.info
}

// IsEnabled returns whether weavelet is enabled
func (wm *WeaveletManager) IsEnabled() bool {
	return wm.config.Enabled
}

// UpdateConfig updates the weavelet configuration
func (wm *WeaveletManager) UpdateConfig(config *WeaveletConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	
	wm.config = config
	return wm.Initialize()
}

// AddLabel adds a label to the weavelet
func (wm *WeaveletManager) AddLabel(key, value string) {
	if wm.config.Labels == nil {
		wm.config.Labels = make(map[string]string)
	}
	wm.config.Labels[key] = value
	
	if wm.info != nil {
		if wm.info.Labels == nil {
			wm.info.Labels = make(map[string]string)
		}
		wm.info.Labels[key] = value
	}
}

// RemoveLabel removes a label from the weavelet
func (wm *WeaveletManager) RemoveLabel(key string) {
	if wm.config.Labels != nil {
		delete(wm.config.Labels, key)
	}
	
	if wm.info != nil && wm.info.Labels != nil {
		delete(wm.info.Labels, key)
	}
}

// GetLabels returns all labels
func (wm *WeaveletManager) GetLabels() map[string]string {
	if wm.info != nil && wm.info.Labels != nil {
		return wm.info.Labels
	}
	return wm.config.Labels
}

// CheckWeaveletInfo checks that weavelet information is well-formed.
func CheckWeaveletInfo(w *WeaveletInfo) error {
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

// ValidateColocationGroup validates a colocation group
func ValidateColocationGroup(group *ColocationGroup) error {
	if group == nil {
		return fmt.Errorf("colocation group cannot be nil")
	}
	if group.Name == "" {
		return fmt.Errorf("colocation group name cannot be empty")
	}
	if len(group.Components) == 0 {
		return fmt.Errorf("colocation group must have at least one component")
	}
	
	// Check for duplicate components
	componentSet := make(map[string]bool)
	for _, component := range group.Components {
		if component == "" {
			return fmt.Errorf("component name cannot be empty")
		}
		if componentSet[component] {
			return fmt.Errorf("duplicate component: %s", component)
		}
		componentSet[component] = true
	}
	
	return nil
}

// CreateWeaveletInfo creates a new WeaveletInfo with validation
func CreateWeaveletInfo(config *WeaveletConfig) (*WeaveletInfo, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	
	info := &WeaveletInfo{
		App:          config.App,
		DeploymentId: config.DeploymentId,
		Group: &ColocationGroup{
			Name:       config.GroupName,
			Components: config.Components,
		},
		GroupId:      uuid.New().String(),
		Id:           uuid.New().String(),
		SingleProcess: config.SingleProcess,
		SingleMachine: config.SingleMachine,
		Labels:       make(map[string]string),
	}
	
	// Copy labels
	for k, v := range config.Labels {
		info.Labels[k] = v
	}
	
	// Validate the created info
	if err := CheckWeaveletInfo(info); err != nil {
		return nil, err
	}
	
	return info, nil
}