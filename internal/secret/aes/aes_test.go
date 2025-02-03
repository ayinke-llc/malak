package aes

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "12345678901234567890123456789012" // 32-byte key
	plaintext := "hello world"

	encrypted, err := encrypt([]byte(plaintext), key)
	require.NoError(t, err)
	require.NotEmpty(t, encrypted)

	decrypted, err := decrypt(encrypted, key)
	require.NoError(t, err)
	require.Equal(t, plaintext, string(decrypted))
}

func TestDecryptInvalidData(t *testing.T) {
	key := "12345678901234567890123456789012"
	invalidData := "invalid-base64"

	_, err := decrypt(invalidData, key)
	require.Error(t, err)
}

func TestDecryptTamperedData(t *testing.T) {
	key := "12345678901234567890123456789012"
	plaintext := "sensitive data"

	encrypted, err := encrypt([]byte(plaintext), key)
	require.NoError(t, err)

	// Tamper with encrypted data
	decoded, _ := base64.StdEncoding.DecodeString(encrypted)
	decoded[0] ^= 0xFF // Modify first byte
	tamperedEncrypted := base64.StdEncoding.EncodeToString(decoded)

	_, err = decrypt(tamperedEncrypted, key)
	require.Error(t, err)
}

func TestAESClient(t *testing.T) {
	cfg := config.Config{}

	cfg.Integration.AES.Key = "12345678901234567890123456789012"

	client, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, client)

	secretClient := client.(*aesClient)
	plaintext := "test secret"

	encrypted, err := secretClient.Create(context.TODO(), &secret.CreateSecretOptions{Value: plaintext})
	require.NoError(t, err)

	decrypted, err := secretClient.Get(context.TODO(), encrypted)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}

func TestAESClientInvalidKey(t *testing.T) {
	cfg := config.Config{}
	cfg.Integration.AES.Key = "1234567890123"

	_, err := New(cfg)
	require.Error(t, err)
}
