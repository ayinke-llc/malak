package jwttoken

import (
	"errors"
	"fmt"
	"time"

	"github.com/ayinke-llc/malak/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// ENUM(access,refresh)
type Purpose uint8

type JWTokenData struct {
	Token     string
	Purpose   Purpose
	UserID    uuid.UUID
	ExpiresAt time.Time
}

type jwtokenManager struct {
	signingKey string
}

type JWTokenManager interface {
	GenerateJWToken(JWTokenData) (JWTokenData, error)
	ParseJWToken(string) (JWTokenData, error)
}

func New(cfg config.Config) JWTokenManager {
	return &jwtokenManager{
		signingKey: cfg.Auth.JWT.Key,
	}
}

func (t *jwtokenManager) GenerateJWToken(data JWTokenData) (JWTokenData, error) {
	claims := jwt.MapClaims{
		"signer": "malak",
		"id":     data.UserID,
		"exp":    time.Now().Add(time.Hour * 168), // 7 days
	}
	jwtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtoken.SignedString([]byte(t.signingKey))
	if err != nil {
		return JWTokenData{}, fmt.Errorf("GenerateJWToken/SignedString: sign jwtoken failed: %w", err)
	}

	data.Token = token
	return data, nil
}

func (t *jwtokenManager) ParseJWToken(JWToken string) (JWTokenData, error) {
	parsedJWToken, err := jwt.Parse(JWToken, func(JWToken *jwt.Token) (i interface{}, e error) {
		if _, ok := JWToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("ParseJWToken/Parse: error: unexpected signing method: %v", JWToken.Header["alg"])
		}
		return []byte(t.signingKey), nil
	})

	if err != nil && parsedJWToken == nil {
		return JWTokenData{}, fmt.Errorf("ParseJWToken/Parse: parse JWToken failed: %w", err)
	}

	claims, ok := parsedJWToken.Claims.(jwt.MapClaims)
	if !ok {
		return JWTokenData{}, fmt.Errorf("ParseJWToken/parsedJWToken.Claims: error: JWToken wrong claims")
	}

	id, ok := claims["id"].(string)
	if !ok {
		return JWTokenData{}, errors.New("user_id not exists")
	}

	expiresAt, ok := claims["exp"].(int64)
	if !ok {
		return JWTokenData{}, errors.New("ParseJWToken/parseJWTokenClaim/exp: expiration date not found")
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return JWTokenData{}, err
	}

	return JWTokenData{
		UserID:    userID,
		ExpiresAt: time.Unix(expiresAt, 0),
	}, nil
}
