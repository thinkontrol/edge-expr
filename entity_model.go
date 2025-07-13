package edgeexpr

type Field struct {
	Key        string `json:"key"`
	Expression string `json:"expression"`
	AsTag      bool   `json:"as_tag,omitempty"` // Optional flag to indicate if the field should be treated as a tag
}

type Event struct {
	Key        string `json:"key"`
	Expression string `json:"expression"`
	Category   string `json:"category,omitempty"` // Optional category for the event
	Level      int    `json:"level,omitempty"`    // Optional level for the event, e.g., 1 for critical, 2 for warning, etc.
	Message    string `json:"message,omitempty"`  // Optional message for the event
}

type EntityModel struct {
	Fields map[string]*Field `json:"fields"` // map of field name to Field struct
	Events map[string]*Event `json:"events"` // map of event name to Event struct
}
