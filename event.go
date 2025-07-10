package edgeexpr

type Event struct {
	Key        string `json:"key"`
	Expression string `json:"expression"`
	Category   string `json:"category,omitempty"` // Optional category for the event
	Level      int    `json:"level,omitempty"`    // Optional level for the event, e.g., 1 for critical, 2 for warning, etc.
	Message    string `json:"message,omitempty"`  // Optional message for the event
}
