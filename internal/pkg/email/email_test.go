package email

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getOptions() SendOptions {
	return SendOptions{
		Plain:     "This is my email in plain text ",
		HTML:      "This is my email in html format",
		Sender:    "yo@lanre.wtf",
		Recipient: "lanre@ayinke.ventures",
		DKIM: struct {
			Sign       bool
			PrivateKey []byte
		}{
			Sign: false,
		},
	}
}

func TestSendOptions_Validate(t *testing.T) {

	t.Run("plain email is empty", func(t *testing.T) {
		opts := getOptions()

		opts.Plain = ""

		err := opts.Validate()
		require.Error(t, err)

		require.Contains(t, err.Error(), "plain copy of email")
	})

	t.Run("html email is empty", func(t *testing.T) {
		opts := getOptions()

		opts.HTML = ""

		err := opts.Validate()
		require.Error(t, err)

		require.Contains(t, err.Error(), "html copy of email")
	})

	t.Run("no sender", func(t *testing.T) {
		opts := getOptions()

		opts.Sender = ""

		err := opts.Validate()
		require.Error(t, err)

		require.Contains(t, err.Error(), "please provide sender")
	})

	t.Run("no Recipient", func(t *testing.T) {
		opts := getOptions()

		opts.Recipient = ""

		err := opts.Validate()
		require.Error(t, err)

		require.Contains(t, err.Error(), "please provide recipient")
	})
}
