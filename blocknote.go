package malak

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

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

type BlockContents []BlockContent

func (bc BlockContents) Text() string {
	return ""
}

func (bc BlockContents) HTML() string {
	var html strings.Builder
	for _, block := range bc {
		html.WriteString(convertBlockToHTML(block))
	}
	return html.String()
}

func convertBlockToHTML(block BlockContent) string {
	var html strings.Builder
	style := getStyleString(block.Props)

	switch block.Type {
	case "heading":
		level, _ := block.Props["level"].(int)
		html.WriteString(fmt.Sprintf("<h%d style='%s'>%s</h%d>", int(level), style, getContent(block.Content), level))
	case "paragraph":

		extra := []string{
			"font-size:100%;",
			"line-height: 24px;",
			"margin: 16px;",
		}

		html.WriteString(fmt.Sprintf("<p style='%s'>%s</p>", strings.Join(extra, style), getContent(block.Content)))
		html.WriteString(`<hr style="width:100%;border:none;border-top:1px solid #eaeaea;border-color:#e6ebf1;margin:20px 0" />`)
	case "numberedListItem":
		html.WriteString(fmt.Sprintf("<li style='%s'>%s</li>", style, getContent(block.Content)))
	case "checkListItem":
		checked := ""
		if isChecked, ok := block.Props["checked"].(bool); ok && isChecked {
			checked = " checked"
		}
		html.WriteString(fmt.Sprintf("<li style='%s'><input type='checkbox'%s>%s</li>", style, checked, getContent(block.Content)))
	case "image":
		url, _ := block.Props["url"].(string)
		name, _ := block.Props["name"].(string)
		html.WriteString(fmt.Sprintf("<img src='%s' alt='%s' style='%s max-width:100px;'>", url, name, style))
		if caption, ok := block.Props["caption"].(string); ok && caption != "" {
			html.WriteString(fmt.Sprintf("<figcaption>%s</figcaption>", template.HTMLEscapeString(caption)))
		}
	}

	// Handle children
	if len(block.Children) > 0 {
		html.WriteString(getContent(block.Children))
	}

	return html.String()
}

func getContent(items []BlockNoteItem) string {
	var result strings.Builder
	for _, item := range items {
		style := getInlineStyle(item.Styles)
		if style != "" {
			result.WriteString(fmt.Sprintf("<span style='%s'>%s</span>", style, html.EscapeString(item.Text)))
		} else {
			result.WriteString(html.EscapeString(item.Text))
		}
	}
	return result.String()
}

func getStyleString(props map[string]any) string {
	var styles []string
	if textColor, ok := props["textColor"].(string); ok && textColor != "default" {
		styles = append(styles, fmt.Sprintf("color:%s;", textColor))
	}
	if bgColor, ok := props["backgroundColor"].(string); ok && bgColor != "default" {
		styles = append(styles, fmt.Sprintf("background-color:%s;", bgColor))
	}
	if textAlign, ok := props["textAlignment"].(string); ok && textAlign != "left" {
		styles = append(styles, fmt.Sprintf("text-align:%s;", textAlign))
	}
	return strings.Join(styles, " ")
}

func getInlineStyle(styles map[string]any) string {
	var styleStrings []string
	for key, value := range styles {
		styleStrings = append(styleStrings, fmt.Sprintf("%s:%v;", key, value))
	}
	return strings.Join(styleStrings, " ")
}
