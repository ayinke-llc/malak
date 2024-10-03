package malak

type BlockContent struct {
	ID       string          `json:"id" validate:"required"`
	Type     string          `json:"type" validate:"required"`
	Content  []BlockNoteItem `json:"content" validate:"required"`
	Children []BlockNoteItem `json:"children" validate:"required"`
	Props    map[string]any  `json:"props" validate:"required"`
}

type BlockNoteItem struct {
	Type   string         `json:"type,omitempty" validate:"required"`
	Text   string         `json:"text,omitempty" validate:"required"`
	Styles map[string]any `json:"styles" validate:"required"`
}
