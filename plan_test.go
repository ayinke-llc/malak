package malak

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {

	tt := []struct {
		name     string
		amount   int64
		takeN    int64
		hasError bool
	}{
		{
			name:     "cannot take from an empty counter",
			amount:   0,
			takeN:    10,
			hasError: true,
		},
		{
			name:     "cannot take more than avaialble",
			amount:   10,
			takeN:    11,
			hasError: true,
		},
		{
			name:     "can take exact value",
			amount:   10,
			takeN:    10,
			hasError: false,
		},
		{
			name:     "can take less value",
			amount:   10,
			takeN:    9,
			hasError: false,
		},
		{
			name:     "zero take 1",
			amount:   0,
			takeN:    1,
			hasError: true,
		},
		{
			name:     "zero take 0",
			amount:   0,
			takeN:    0,
			hasError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := Counter(tc.amount)
			err := c.TakeN(tc.takeN)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
