package malak

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsImageFromURL(t *testing.T) {

	tt := []struct {
		endpoint string
		hasError bool
		name     string
	}{
		{
			name:     "no url provided",
			hasError: true,
			endpoint: "",
		},
		{
			name:     "bad url",
			hasError: true,
			endpoint: "http://localhost:44000",
		},
		{
			name:     "google.com",
			hasError: true,
			endpoint: "https://google.com",
		},
		{
			name:     "unsplash",
			hasError: false,
			endpoint: "https://images.unsplash.com/photo-1737467023078-a694673d7cb3",
		},
	}

	for _, tc := range tt {

		t.Run(tc.name, func(t *testing.T) {
			r, err := IsImageFromURL(tc.endpoint)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.True(t, r)
			require.NoError(t, err)
		})
	}
}
