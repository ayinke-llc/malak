package smtp

import (
	"context"
	"testing"

	"github.com/ayinke-llc/malak"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/pkg/email"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSMTP_Send(t *testing.T) {

	containerReq := testcontainers.ContainerRequest{
		Image:        "axllent/mailpit",
		ExposedPorts: []string{"8025/tcp", "1025/tcp"},
		WaitingFor:   wait.ForHTTP("/").WithPort("8025/tcp"),
		Env: map[string]string{
			"MP_SMTP_AUTH_ACCEPT_ANY":     "1",
			"MP_SMTP_AUTH_ALLOW_INSECURE": "1",
		},
	}

	mailpitContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	require.NoError(t, err)

	port, err := mailpitContainer.MappedPort(context.Background(), "1025")
	require.NoError(t, err)

	client, err := New(getConfig(port.Int()))

	require.NoError(t, err)

	err = client.Send(context.Background(), email.SendOptions{
		HTML:      "This is my email in html format",
		Sender:    "yo@lanre.wtf",
		Recipient: "lanre@ayinke.ventures",
		DKIM: struct {
			Sign       bool
			PrivateKey []byte
		}{
			Sign: false,
		},
	})

	require.NoError(t, err)
}

func TestNew_Errors(t *testing.T) {

	t.Run("smtp host is empty", func(t *testing.T) {

		cfg := getConfig(1025)
		cfg.Email.SMTP.Host = ""

		_, err := New(cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "smtp host")
	})

	t.Run("smtp usernasme is empty", func(t *testing.T) {

		cfg := getConfig(1025)
		cfg.Email.SMTP.Username = ""

		_, err := New(cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "smtp username")
	})

	t.Run("smtp password is empty", func(t *testing.T) {

		cfg := getConfig(1025)
		cfg.Email.SMTP.Password = ""

		_, err := New(cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "smtp password")
	})

}

func getConfig(port int) config.Config {
	return config.Config{
		Email: struct {
			Provider   config.EmailProvider "mapstructure:\"provider\" yaml:\"provider\""
			Sender     malak.Email          "mapstructure:\"sender\" yaml:\"sender\""
			SenderName string               "mapstructure:\"sender_name\" yaml:\"sender_name\""
			SMTP       struct {
				Host     string "mapstructure:\"host\" yaml:\"host\""
				Port     int    "mapstructure:\"port\" yaml:\"port\""
				Username string "mapstructure:\"username\" yaml:\"username\""
				Password string "mapstructure:\"password\" yaml:\"password\""
				UseTLS   bool   "yaml:\"use_tls\" mapstructure:\"use_tls\""
			} "mapstructure:\"smtp\" yaml:\"smtp\""
		}{
			Provider:   config.EmailProviderSmtp,
			Sender:     malak.Email("yo@oops.com"),
			SenderName: "Malak Updates",
			SMTP: struct {
				Host     string "mapstructure:\"host\" yaml:\"host\""
				Port     int    "mapstructure:\"port\" yaml:\"port\""
				Username string "mapstructure:\"username\" yaml:\"username\""
				Password string "mapstructure:\"password\" yaml:\"password\""
				UseTLS   bool   "yaml:\"use_tls\" mapstructure:\"use_tls\""
			}{
				Username: "random",
				Password: "random",
				Port:     port,
				Host:     "localhost",
				UseTLS:   false,
			},
		},
	}
}
