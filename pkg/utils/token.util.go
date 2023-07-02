package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// private dan public tokennya adalah utf-8 yang di encode ke base64 saat akan membuat atau memvalidasi token maka tokennya di kembalikan ke utf-8 untuk masuk di func jwt.ParseRSAPrivateKeyFromPEM

func CreateToken(ttl time.Duration, payload any, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}

	// parse the key from base64 back to utf8
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"sub": payload,
		"exp": now.Add(ttl).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	fmt.Println(claims)
	return token, nil
}

func ValidateToken(token string, publicKey string) (any, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		fmt.Println(t)
		// write any validation u want
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		// return the key to parse function
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	// if cannot convert to jwt mapclaims or the parsed token is invalid
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}
	fmt.Println(claims)
	return claims["sub"], nil
}
