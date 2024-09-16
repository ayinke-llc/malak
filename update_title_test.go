package malak

import (
	"testing"
)

func BenchmarkGetFirstHeader(b *testing.B) {
	markdown := `# First Title
	
Some content here
	
## Second Title
More content`

	getFirstHeader(markdown)
}
