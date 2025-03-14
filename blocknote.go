package malak

import (
	"fmt"
	"html"
	"strings"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

type Block struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Props    map[string]interface{} `json:"props"`
	Content  interface{}            `json:"content"`
	Children []Block                `json:"children"`
}

type DefaultProps struct {
	BackgroundColor string `json:"backgroundColor"`
	TextColor       string `json:"textColor"`
	TextAlignment   string `json:"textAlignment"`
}

type ParagraphBlock struct {
	Block
	Content []InlineContent `json:"content"`
}

type HeadingBlock struct {
	Block
	Props struct {
		DefaultProps
		Level int `json:"level"`
	} `json:"props"`
	Content []InlineContent `json:"content"`
}

type BulletListItemBlock struct {
	Block
	Content []InlineContent `json:"content"`
}

type NumberedListItemBlock struct {
	Block
	Content []InlineContent `json:"content"`
}

type ImageBlock struct {
	Block
	Props struct {
		DefaultProps
		URL          string `json:"url"`
		Caption      string `json:"caption"`
		PreviewWidth int    `json:"previewWidth"`
	} `json:"props"`
}

type AlertBlock struct {
	Block
	Props struct {
		DefaultProps
		Type string `json:"type"` // warning, error, info, success
	} `json:"props"`
	Content []InlineContent `json:"content"`
}

type DashboardBlock struct {
	Block
	Props struct {
		DefaultProps
		SelectedItem string `json:"selectedItem"`
	} `json:"props"`
	Content []InlineContent `json:"content"`
}

type ChartBlock struct {
	Block
	Props struct {
		DefaultProps
		ChartType     string      `json:"chartType"`
		Data          interface{} `json:"data"`
		SelectedChart string      `json:"selectedChart"` // Integration chart ID
	} `json:"props"`
	Content []InlineContent `json:"content"`
}

type TableBlock struct {
	Block
	Content TableContent `json:"content"`
}

type TableContent struct {
	Type string `json:"type"`
	Rows []struct {
		Cells [][]InlineContent `json:"cells"`
	} `json:"rows"`
}

type InlineContent interface {
	inlineContent()
}

type StyledText struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
	Styles Styles `json:"styles"`
}

func (StyledText) inlineContent() {}

type Link struct {
	Type    string       `json:"type"`
	Content []StyledText `json:"content"`
	Href    string       `json:"href"`
}

func (Link) inlineContent() {}

type Styles struct {
	Bold            bool   `json:"bold"`
	Italic          bool   `json:"italic"`
	Underline       bool   `json:"underline"`
	Strikethrough   bool   `json:"strikethrough"`
	TextColor       string `json:"textColor"`
	BackgroundColor string `json:"backgroundColor"`
}

type BlockContents []Block

func (bc BlockContents) HTML(workspaceID uuid.UUID, renderer ChartRenderer) string {
	var html strings.Builder
	for _, block := range bc {
		html.WriteString(convertBlockToHTML(block, workspaceID, renderer))
	}
	return html.String()
}

func convertBlockToHTML(block Block, workspaceID uuid.UUID, renderer ChartRenderer) string {
	var s strings.Builder
	style := getStyleString(block.Props)

	switch block.Type {
	case "heading":
		level, _ := block.Props["level"].(float64)
		content := getSimpleContent(block.Content)
		s.WriteString(fmt.Sprintf("<h%d style='%s'>%s</h%d>", int(level), style, content, int(level)))
	case "paragraph":
		extra := []string{
			"font-size:100%;",
			"line-height: 24px;",
			"margin: 16px;",
		}
		content := getSimpleContent(block.Content)
		s.WriteString(fmt.Sprintf("<p style='%s'>%s</p>", strings.Join(append(extra, style), " "), content))
	case "alert":
		alertType, _ := block.Props["type"].(string)
		content := getSimpleContent(block.Content)
		s.WriteString(fmt.Sprintf(`<div class="alert" data-alert-type="%s">%s</div>`, alertType, content))
	case "dashboard":
		selectedItem, _ := block.Props["selectedItem"].(string)
		content := getSimpleContent(block.Content)
		s.WriteString(fmt.Sprintf(`<div class="dashboard" data-selected-item="%s">%s</div>`, selectedItem, content))
	case "chart":
		selectedChart, _ := block.Props["selectedChart"].(string)
		if selectedChart == "" {
			s.WriteString(`<div class="chart-error">No chart selected</div>`)
			return s.String()
		}

		chartKey, err := renderer.RenderChart(workspaceID, selectedChart)
		if err != nil {
			s.WriteString(fmt.Sprintf(`<div class="chart-error">Failed to render chart: %s</div>`, err.Error()))
		} else {
			s.WriteString(fmt.Sprintf(`<a href='%s' target="_blank"><img src='%s' alt='%s' style="display: block; width: %s; max-width: 600px; height: auto;"></a>`, chartKey, chartKey, "Chart image", "100%"))
		}
	case "numberedListItem", "bulletListItem":
		content := getSimpleContent(block.Content)
		s.WriteString(fmt.Sprintf("<li style='%s'>%s</li>", style, content))
	case "checkListItem":
		checked := ""
		if isChecked, ok := block.Props["checked"].(bool); ok && isChecked {
			checked = " checked"
		}
		content := getSimpleContent(block.Content)
		s.WriteString(fmt.Sprintf("<li style='%s'><input type='checkbox'%s>%s</li>", style, checked, content))
	case "image":
		url, _ := block.Props["url"].(string)
		name, _ := block.Props["name"].(string)

		s.WriteString(fmt.Sprintf(`<a href='%s' target="_blank"><img src='%s' alt='%s' style="display: block; width: %s; max-width: 600px; height: auto;"></a>`, url, url, name, "100%"))

		_ = style

		if caption, ok := block.Props["caption"].(string); ok && caption != "" {
			s.WriteString(fmt.Sprintf("<figcaption>%s</figcaption>", html.EscapeString(caption)))
		}
	default:
	}

	if len(block.Children) > 0 {
		for _, child := range block.Children {
			s.WriteString(convertBlockToHTML(child, workspaceID, renderer))
		}
	}

	return s.String()
}

func getSimpleContent(content interface{}) string {
	var result strings.Builder

	switch v := content.(type) {
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				text, _ := m["text"].(string)
				styles, _ := m["styles"].(map[string]interface{})
				result.WriteString(applyInlineStyles(text, styles))
			}
		}
	case string:
		result.WriteString(html.EscapeString(v))
	default:
	}

	return result.String()
}

func applyInlineStyles(text string, styles map[string]interface{}) string {
	if len(styles) == 0 {
		return html.EscapeString(text)
	}

	var styleStrings []string
	if bold, ok := styles["bold"].(bool); ok && bold {
		styleStrings = append(styleStrings, "font-weight: bold;")
	}
	if italic, ok := styles["italic"].(bool); ok && italic {
		styleStrings = append(styleStrings, "font-style: italic;")
	}
	if underline, ok := styles["underline"].(bool); ok && underline {
		styleStrings = append(styleStrings, "text-decoration: underline;")
	}
	if strikethrough, ok := styles["strikethrough"].(bool); ok && strikethrough {
		styleStrings = append(styleStrings, "text-decoration: line-through;")
	}
	if textColor, ok := styles["textColor"].(string); ok && textColor != "" {
		styleStrings = append(styleStrings, fmt.Sprintf("color: %s;", textColor))
	}
	if backgroundColor, ok := styles["backgroundColor"].(string); ok && backgroundColor != "" {
		styleStrings = append(styleStrings, fmt.Sprintf("background-color: %s;", backgroundColor))
	}

	if len(styleStrings) > 0 {
		return fmt.Sprintf("<span style='%s'>%s</span>", strings.Join(styleStrings, " "), html.EscapeString(text))
	}

	return html.EscapeString(text)
}

func getStyleString(props map[string]interface{}) string {

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

func SanitizeBlocks(blocks BlockContents) (BlockContents, error) {
	policy := bluemonday.UGCPolicy()

	for i := range blocks {
		sanitizeBlock(&blocks[i], policy)
	}

	return blocks, nil
}

func sanitizeBlock(block *Block, policy *bluemonday.Policy) {
	// Sanitize Props
	for key, value := range block.Props {
		if strValue, ok := value.(string); ok {
			block.Props[key] = policy.Sanitize(strValue)
		}
	}

	// Sanitize Content based on block type
	switch block.Type {
	case "alert":
		if alertType, ok := block.Props["type"].(string); ok {
			// Only allow valid alert types
			validTypes := map[string]bool{
				"warning": true,
				"error":   true,
				"info":    true,
				"success": true,
			}
			if !validTypes[alertType] {
				block.Props["type"] = "warning" // default to warning if invalid
			}
		}
		fallthrough // continue to sanitize content
	case "chart":
		if chartType, ok := block.Props["chartType"].(string); ok {
			// Sanitize chart type
			validTypes := map[string]bool{
				"bar":  true,
				"line": true,
				"pie":  true,
				"area": true,
			}
			if !validTypes[chartType] {
				block.Props["chartType"] = "bar" // default to bar if invalid
			}
		}
		fallthrough // continue to sanitize content
	case "dashboard":
		if content, ok := block.Content.([]interface{}); ok {
			for i := range content {
				if item, ok := content[i].(map[string]interface{}); ok {
					if text, ok := item["text"].(string); ok {
						item["text"] = policy.Sanitize(text)
					}
					if styles, ok := item["styles"].(map[string]interface{}); ok {
						for key, value := range styles {
							if strValue, ok := value.(string); ok {
								styles[key] = policy.Sanitize(strValue)
							}
						}
					}
				}
			}
		}
	default:
		switch content := block.Content.(type) {
		case []InlineContent:
			for j := range content {
				sanitizeInlineContent(content[j], policy)
			}
		case TableContent:
			sanitizeTableContent(&content, policy)
			block.Content = content
		}
	}

	// Recursively sanitize children
	for i := range block.Children {
		sanitizeBlock(&block.Children[i], policy)
	}
}

func sanitizeInlineContent(content InlineContent, policy *bluemonday.Policy) {
	switch v := content.(type) {
	case *StyledText:
		v.Text = policy.Sanitize(v.Text)
		sanitizeStyles(&v.Styles, policy)
	case *Link:
		v.Href = policy.Sanitize(v.Href)
		for i := range v.Content {
			v.Content[i].Text = policy.Sanitize(v.Content[i].Text)
			sanitizeStyles(&v.Content[i].Styles, policy)
		}
	}
}

func sanitizeStyles(styles *Styles, policy *bluemonday.Policy) {
	styles.TextColor = policy.Sanitize(styles.TextColor)
	styles.BackgroundColor = policy.Sanitize(styles.BackgroundColor)
}

func sanitizeTableContent(table *TableContent, policy *bluemonday.Policy) {
	for i := range table.Rows {
		for j := range table.Rows[i].Cells {
			for k := range table.Rows[i].Cells[j] {
				sanitizeInlineContent(table.Rows[i].Cells[j][k], policy)
			}
		}
	}
}

type ChartRenderer interface {
	RenderChart(workspaceID uuid.UUID, chartID string) (string, error)
}
