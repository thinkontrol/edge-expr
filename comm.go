package edgeexpr

import (
	"time"
)

type PushValue struct {
	// Key       string     `json:"key"`
	Value     any        `json:"value"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

type Command struct {
	CommandID string         `json:"command_id"`
	Command   string         `json:"command"`
	Payload   map[string]any `json:"payload,omitempty"`
	Timestamp *time.Time     `json:"timestamp,omitempty"`
}

type CommandResponse struct {
	CommandID string         `json:"command_id"`
	Message   string         `json:"message,omitempty"`
	Success   bool           `json:"success"`
	Payload   map[string]any `json:"payload,omitempty"`
	Timestamp *time.Time     `json:"timestamp,omitempty"`
}
