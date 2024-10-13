package malak

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStyleString(t *testing.T) {

	props := map[string]any{
		"textColor":       "red",
		"backgroundColor": "blue",
		"textAlignment":   "align-center",
	}

	s := getStyleString(props)

	require.NotEmpty(t, s)

	require.Contains(t, s, "color:red")
	require.Contains(t, s, "background-color:blue")
}
