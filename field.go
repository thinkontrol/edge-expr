package edgeexpr

type Field struct {
	Key        string `json:"key"`
	Expression string `json:"expression"`
	AsTag      bool   `json:"as_tag,omitempty"` // Optional flag to indicate if the field should be treated as a tag
}
