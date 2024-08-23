package jwttoken

import (
	"testing"

	"github.com/ayinke-llc/malak/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func getConfig() config.Config {
	return config.Config{
		Auth: struct {
			Google struct {
				ClientID     string   "yaml:\"client_id\" mapstructure:\"client_id\""
				ClientSecret string   "yaml:\"client_secret\" mapstructure:\"client_secret\""
				RedirectURI  string   "yaml:\"redirect_uri\" mapstructure:\"redirect_uri\""
				Scopes       []string "yaml:\"scopes\" mapstructure:\"scopes\""
				IsEnabled    bool     "yaml:\"is_enabled\" mapstructure:\"is_enabled\""
			} "yaml:\"google\" mapstructure:\"google\""
			JWT struct {
				Key string "yaml:\"key\" mapstructure:\"key\""
			} "yaml:\"jwt\" mapstructure:\"jwt\""
		}{
			JWT: struct {
				Key string "yaml:\"key\" mapstructure:\"key\""
			}{
				Key: "a907e75f80910f5dc5b8c677de1de611ffa80be9d7d9f9dd614c8c7846db1062",
			},
		},
	}
}

func TestJWT_Parse(t *testing.T) {

	manager := New(getConfig())

	userID := uuid.New()

	token, err := manager.GenerateJWToken(JWTokenData{
		UserID:  userID,
		Purpose: PurposeAccess,
	})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	parsedToken, err := manager.ParseJWToken(token.Token)
	require.NoError(t, err)

	t.Log(parsedToken.ExpiresAt)
	require.Equal(t, userID, parsedToken.UserID)
}

func TestJWT_Generate(t *testing.T) {

	manager := New(getConfig())

	token, err := manager.GenerateJWToken(JWTokenData{
		UserID:  uuid.New(),
		Purpose: PurposeAccess,
	})
	require.NoError(t, err)

	require.NotEmpty(t, token.Token)
}
