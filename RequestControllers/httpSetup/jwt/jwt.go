package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	DevMode       = true
	DevClientName = "Client Three Inc"
	DevClientID   = "1a1b24b8-f439-4334-a91c-ba30a814614c"
)

type ClientClaims struct {
	ClientName string `json:"client_name"`
	ClientID   string `json:"client_id"`
	jwt.StandardClaims
}

type JWTManager struct {
	secretKey   string
	expiryHours int
}

func NewJWTManager(secretKey string, expiryHours int) *JWTManager {
	return &JWTManager{
		secretKey:   secretKey,
		expiryHours: expiryHours,
	}
}

func (m *JWTManager) GenerateToken(clientName, clientID string) (string, error) {
	claims := &ClientClaims{
		ClientName: clientName,
		ClientID:   clientID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(m.expiryHours)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) ValidateToken(tokenStr string) (*ClientClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&ClientClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*ClientClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
