package edgeexpr

import (
	"time"
)

type PushValue struct {
	Key       string     `json:"key" mapstructure:"key"`
	Value     any        `json:"value" mapstructure:"value"`
	Timestamp *time.Time `json:"timestamp,omitempty" mapstructure:"timestamp"`
}

type Command struct {
	CommandID string         `json:"command_id" mapstructure:"command_id"`
	Command   string         `json:"command" mapstructure:"command"`
	Payload   map[string]any `json:"payload,omitempty" mapstructure:"payload"`
	Timestamp *time.Time     `json:"timestamp,omitempty" mapstructure:"timestamp"`
}

type CommandResponse struct {
	CommandID string         `json:"command_id" mapstructure:"command_id"`
	Message   string         `json:"message,omitempty" mapstructure:"message"`
	Success   bool           `json:"success" mapstructure:"success"`
	Payload   map[string]any `json:"payload,omitempty" mapstructure:"payload"`
	Timestamp *time.Time     `json:"timestamp,omitempty" mapstructure:"timestamp"`
}
