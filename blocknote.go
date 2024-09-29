package malak

type BlockContent struct {
	Type    string          `json:"type"`
	Content []BlockNoteItem `json:"content"`
	Props   map[string]any  `json:"props,omitempty"`
}

type BlockNoteItem struct {
	Type  string         `json:"type"`
	Text  string         `json:"text,omitempty"`
	Attrs map[string]any `json:"attrs,omitempty"`
}
