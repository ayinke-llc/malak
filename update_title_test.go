package malak

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkGetFirstHeader(b *testing.B) {
	markdown := `# First Title
	
Some content here
	
## Second Title
More content`

	_, _ = getFirstHeader(UpdateContent(markdown))
}

func TestGetFirstHeader(t *testing.T) {

	tt := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "properly formatted markdown",
			content: `
# First Title
	
Some content here
	
## Second Title
More content
			`,
			expected: "Second Title",
		},
		{
			name: "poorly formatted markdown",
			content: `
# First Title
	
Some content here
	
#### weird Title
More content

## real header
			`,
			expected: "real header",
		},
		{
			name: "no 2nd header",
			content: `
# First Title
	
Some content here
			`,
			expected: "",
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			title, err := getFirstHeader(UpdateContent(v.content))
			require.NoError(t, err)

			require.Equal(t, v.expected, title)
		})
	}
}
