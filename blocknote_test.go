package malak

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSimpleContent(t *testing.T) {
	tests := []struct {
		name     string
		content  interface{}
		expected string
	}{
		{
			name: "Single StyledText",
			content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Hello, World!",
					"styles": map[string]interface{}{
						"bold": true,
					},
				},
			},
			expected: "<span style='font-weight: bold;'>Hello, World!</span>",
		},
		{
			name: "Multiple StyledText",
			content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "Hello, ",
					"styles": map[string]interface{}{
						"bold": true,
					},
				},
				map[string]interface{}{
					"type": "text",
					"text": "World!",
					"styles": map[string]interface{}{
						"italic": true,
					},
				},
			},
			expected: "<span style='font-weight: bold;'>Hello, </span><span style='font-style: italic;'>World!</span>",
		},
		{
			name:     "String Content",
			content:  "Plain Text",
			expected: "Plain Text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSimpleContent(tt.content)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestApplyInlineStyles(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		styles   map[string]interface{}
		expected string
	}{
		{
			name:     "No Styles",
			text:     "Plain Text",
			styles:   map[string]interface{}{},
			expected: "Plain Text",
		},
		{
			name: "Bold",
			text: "Bold Text",
			styles: map[string]interface{}{
				"bold": true,
			},
			expected: "<span style='font-weight: bold;'>Bold Text</span>",
		},
		{
			name: "Multiple Styles",
			text: "Styled Text",
			styles: map[string]interface{}{
				"bold":            true,
				"italic":          true,
				"textColor":       "blue",
				"backgroundColor": "lightgray",
			},
			expected: "<span style='font-weight: bold; font-style: italic; color: blue; background-color: lightgray;'>Styled Text</span>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyInlineStyles(tt.text, tt.styles)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetStyleString(t *testing.T) {
	tests := []struct {
		name     string
		props    map[string]interface{}
		expected string
	}{
		{
			name:     "Empty Props",
			props:    map[string]interface{}{},
			expected: "",
		},
		{
			name: "Text Color",
			props: map[string]interface{}{
				"textColor": "red",
			},
			expected: "color:red;",
		},
		{
			name: "Multiple Properties",
			props: map[string]interface{}{
				"textColor":       "blue",
				"backgroundColor": "lightgray",
				"textAlignment":   "right",
			},
			expected: "color:blue; background-color:lightgray; text-align:right;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStyleString(tt.props)
			require.Equal(t, tt.expected, result)
		})
	}
}
