package edgeexpr

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Event struct {
	State     bool       `json:"state"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}
type PushEvent struct {
	Events    map[string]Event  `json:"events"`
	Tags      map[string]string `json:"tags"`
	Timestamp *time.Time        `json:"timestamp,omitempty"`
}
type Field struct {
	Value     interface{} `json:"value"`
	Timestamp *time.Time  `json:"timestamp,omitempty"`
}
type PushField struct {
	Fields    map[string]Field  `json:"fields"`
	Tags      map[string]string `json:"tags"`
	Timestamp *time.Time        `json:"timestamp,omitempty"`
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

func (p *PushField) String() string {
	var lines []string

	// 添加时间戳信息
	if p.Timestamp != nil {
		lines = append(lines, fmt.Sprintf("Timestamp: %s", p.Timestamp.Format("2006-01-02 15:04:05.000")))
	}

	// 添加Tags信息
	if len(p.Tags) > 0 {
		lines = append(lines, "Tags:")
		tagKeys := make([]string, 0, len(p.Tags))
		for key := range p.Tags {
			tagKeys = append(tagKeys, key)
		}
		sort.Strings(tagKeys)
		for _, key := range tagKeys {
			lines = append(lines, fmt.Sprintf("  %s: %s", key, p.Tags[key]))
		}
	}

	// 添加Fields信息
	if len(p.Fields) > 0 {
		lines = append(lines, "Fields:")
		fieldKeys := make([]string, 0, len(p.Fields))
		for key := range p.Fields {
			fieldKeys = append(fieldKeys, key)
		}
		sort.Strings(fieldKeys)
		for _, key := range fieldKeys {
			field := p.Fields[key]
			if field.Timestamp != nil {
				lines = append(lines, fmt.Sprintf("  %s: %v (timestamp: %s)",
					key, field.Value, field.Timestamp.Format("2006-01-02 15:04:05.000")))
			} else {
				lines = append(lines, fmt.Sprintf("  %s: %v", key, field.Value))
			}
		}
	}

	return strings.Join(lines, "\n")
}

func (p *PushEvent) String() string {
	var lines []string

	// 添加时间戳信息
	if p.Timestamp != nil {
		lines = append(lines, fmt.Sprintf("Timestamp: %s", p.Timestamp.Format("2006-01-02 15:04:05.000")))
	}

	// 添加Tags信息
	if len(p.Tags) > 0 {
		lines = append(lines, "Tags:")
		tagKeys := make([]string, 0, len(p.Tags))
		for key := range p.Tags {
			tagKeys = append(tagKeys, key)
		}
		sort.Strings(tagKeys)
		for _, key := range tagKeys {
			lines = append(lines, fmt.Sprintf("  %s: %s", key, p.Tags[key]))
		}
	}

	// 添加Events信息
	if len(p.Events) > 0 {
		lines = append(lines, "Events:")
		eventKeys := make([]string, 0, len(p.Events))
		for key := range p.Events {
			eventKeys = append(eventKeys, key)
		}
		sort.Strings(eventKeys)
		for _, key := range eventKeys {
			event := p.Events[key]
			if event.Timestamp != nil {
				lines = append(lines, fmt.Sprintf("  %s: %t (timestamp: %s)",
					key, event.State, event.Timestamp.Format("2006-01-02 15:04:05.000")))
			} else {
				lines = append(lines, fmt.Sprintf("  %s: %t", key, event.State))
			}
		}
	}

	return strings.Join(lines, "\n")
}
