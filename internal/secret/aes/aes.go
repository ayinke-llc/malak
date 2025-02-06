package aes

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/ayinke-llc/hermes"
	"github.com/ayinke-llc/malak/config"
	"github.com/ayinke-llc/malak/internal/secret"
)

type aesClient struct {
	key string
}

func New(cfg config.Config) (secret.SecretClient, error) {
	if hermes.IsStringEmpty(cfg.Secrets.AES.Key) {
		return nil, errors.New("please provide your AES key")
	}

	keyLen := len(cfg.Secrets.AES.Key)
	if keyLen != 32 {
		return nil, errors.New("AES key must be 32 bytes")
	}

	return &aesClient{
		key: cfg.Secrets.AES.Key,
	}, nil
}

func (i *aesClient) Close() error {
	return nil
}

func (i *aesClient) Get(ctx context.Context,
	key string) (string, error) {

	val, err := decrypt(key, i.key)
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (i *aesClient) Create(ctx context.Context,
	opts *secret.CreateSecretOptions) (string, error) {
	return encrypt([]byte(opts.Value), i.key)
}

func encrypt(plaintext []byte, keyString string) (string, error) {
	key := []byte(keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and include nonce in the final output
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Encode as base64 for safe storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encodedData string, keyString string) ([]byte, error) {
	key := []byte(keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decode the base64 data
	ciphertext, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, encryptedData, nil)
}
