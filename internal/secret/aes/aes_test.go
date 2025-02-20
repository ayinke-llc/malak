package aes

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		key       string
		wantErr   bool
	}{
		{
			name:      "basic encryption and decryption",
			plaintext: "hello world",
			key:       "12345678901234567890123456789012",
			wantErr:   false,
		},
		{
			name:      "empty string",
			plaintext: "",
			key:       "12345678901234567890123456789012",
			wantErr:   false,
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-=[]{}|;:,.<>?`~'\"\\",
			key:       "12345678901234567890123456789012",
			wantErr:   false,
		},
		{
			name:      "unicode characters",
			plaintext: "Hello 世界 🌍 привет мир",
			key:       "12345678901234567890123456789012",
			wantErr:   false,
		},
		{
			name:      "very long string",
			plaintext: strings.Repeat("a", 1000),
			key:       "12345678901234567890123456789012",
			wantErr:   false,
		},
		{
			name:      "invalid key length",
			plaintext: "hello world",
			key:       "short-key",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := encrypt([]byte(tt.plaintext), tt.key)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, encrypted)

			decrypted, err := decrypt(encrypted, tt.key)
			require.NoError(t, err)
			require.Equal(t, tt.plaintext, string(decrypted))
		})
	}
}

func TestDecryptInvalidData(t *testing.T) {
	tests := []struct {
		name       string
		ciphertext string
		key        string
	}{
		{
			name:       "invalid base64",
			ciphertext: "invalid-base64",
			key:        "12345678901234567890123456789012",
		},
		{
			name:       "empty string",
			ciphertext: "",
			key:        "12345678901234567890123456789012",
		},
		{
			name:       "malformed base64",
			ciphertext: "SGVsbG8gV29ybGQ=!", // Valid base64 with invalid character
			key:        "12345678901234567890123456789012",
		},
		{
			name:       "too short after base64 decode",
			ciphertext: base64.StdEncoding.EncodeToString([]byte("short")),
			key:        "12345678901234567890123456789012",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decrypt(tt.ciphertext, tt.key)
			require.Error(t, err)
		})
	}
}

func TestDecryptTamperedData(t *testing.T) {
	key := "12345678901234567890123456789012"
	tests := []struct {
		name      string
		plaintext string
		tamperFn  func([]byte) []byte
	}{
		{
			name:      "modify first byte",
			plaintext: "sensitive data",
			tamperFn: func(b []byte) []byte {
				b[0] ^= 0xFF
				return b
			},
		},
		{
			name:      "modify last byte",
			plaintext: "sensitive data",
			tamperFn: func(b []byte) []byte {
				b[len(b)-1] ^= 0xFF
				return b
			},
		},
		{
			name:      "truncate data",
			plaintext: "sensitive data",
			tamperFn: func(b []byte) []byte {
				return b[:len(b)-1]
			},
		},
		{
			name:      "append data",
			plaintext: "sensitive data",
			tamperFn: func(b []byte) []byte {
				return append(b, 0xFF)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := encrypt([]byte(tt.plaintext), key)
			require.NoError(t, err)

			decoded, err := base64.StdEncoding.DecodeString(encrypted)
			require.NoError(t, err)

			tampered := tt.tamperFn(decoded)
			tamperedEncrypted := base64.StdEncoding.EncodeToString(tampered)

			_, err = decrypt(tamperedEncrypted, key)
			require.Error(t, err)
		})
	}
}

func TestAESClient(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		plaintext string
		wantErr   bool
	}{
		{
			name:      "valid key and plaintext",
			key:       "12345678901234567890123456789012",
			plaintext: "test secret",
			wantErr:   false,
		},
		{
			name:      "empty key",
			key:       "",
			plaintext: "test secret",
			wantErr:   true,
		},
		{
			name:      "short key",
			key:       "short",
			plaintext: "test secret",
			wantErr:   true,
		},
		{
			name:      "exactly 32 bytes key",
			key:       "12345678901234567890123456789012",
			plaintext: "test secret",
			wantErr:   false,
		},
		{
			name:      "key longer than 32 bytes should error",
			key:       "123456789012345678901234567890123", // 33 bytes
			plaintext: "test secret",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{}
			cfg.Secrets.AES.Key = tt.key

			client, err := New(cfg)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, client)

			secretClient := client.(*aesClient)
			encrypted, err := secretClient.Create(t.Context(), &secret.CreateSecretOptions{Value: tt.plaintext})
			require.NoError(t, err)

			decrypted, err := secretClient.Get(t.Context(), encrypted)
			require.NoError(t, err)
			require.Equal(t, tt.plaintext, decrypted)
		})
	}
}

func TestAESClientContextCancellation(t *testing.T) {
	cfg := config.Config{}
	cfg.Secrets.AES.Key = "12345678901234567890123456789012"

	client, err := New(cfg)
	require.NoError(t, err)

	// First create a valid encrypted value
	plaintext := "test value"
	encrypted, err := client.Create(t.Context(), &secret.CreateSecretOptions{Value: plaintext})
	require.NoError(t, err)

	// Now test with cancelled context
	ctx, cancel := context.WithCancel(t.Context())
	cancel() // Cancel context immediately

	// Test Create with cancelled context
	_, err = client.Create(ctx, &secret.CreateSecretOptions{Value: "test"})
	require.NoError(t, err) // Should still work as the operation is synchronous

	// Test Get with cancelled context using valid encrypted value
	decrypted, err := client.Get(ctx, encrypted)
	require.NoError(t, err) // Should still work as the operation is synchronous
	require.Equal(t, plaintext, decrypted)
}
