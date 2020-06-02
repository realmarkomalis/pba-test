package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
	secretKey []byte
	Issuer    string
	ValidFor  time.Duration
}

type tokenServiceClaims struct {
	jwt.StandardClaims
	Claims interface{} `json:"claims"`
}

func NewTokenService(secretKey, issuer string, minutesValid int) *TokenService {
	validFor := time.Duration(minutesValid) * time.Minute

	return &TokenService{
		secretKey: []byte(secretKey),
		Issuer:    issuer,
		ValidFor:  validFor,
	}
}

func (t *TokenService) keyFunc(token *jwt.Token) (interface{}, error) {
	return t.secretKey, nil
}

func (t *TokenService) Decode(tokenString string, out interface{}) error {
	claims := tokenServiceClaims{
		Claims: out,
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims, t.keyFunc)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid token")
	}

	if _, ok := token.Claims.(*tokenServiceClaims); !ok {
		return errors.New("Failed to convert claims to the provided output interface")
	}

	return nil
}

func (t *TokenService) Encode(baseClaims interface{}) (string, error) {
	claims := tokenServiceClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(t.ValidFor).Unix(),
			Issuer:    t.Issuer,
		},
		baseClaims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *TokenService) NewToken(baseClaims interface{}) (string, error) {
	claims := tokenServiceClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(t.ValidFor).Unix(),
			Issuer:    t.Issuer,
		},
		baseClaims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *TokenService) IsValidToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, t.keyFunc)

	if err != nil {
		return false, fmt.Errorf("Error validating token: %s", err)
	}

	return token.Valid, nil
}
