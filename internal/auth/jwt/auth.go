package authJWT

import (
	"crypto/sha256"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenTTL  = 1 * time.Hour
	refreshTokenTTL = 7 * 24 * time.Hour
)

type JWT struct {
	accessTokenSecret  string
	refreshTokenSecret string
}

func deriveSecret(base, purpose string) string {
	hasher := sha256.New()
	hasher.Write([]byte(base + ":" + purpose))
	return string(hasher.Sum(nil))
}

func create(secret string) *JWT {
	return &JWT{
		accessTokenSecret:  deriveSecret(secret, "access"),
		refreshTokenSecret: deriveSecret(secret, "refresh"),
	}
}

func (j *JWT) GenerateAccessToken(claims jwt.MapClaims) (string, error) {
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(accessTokenTTL).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.accessTokenSecret)
}

func (j *JWT) GenerateRefreshToken(claims jwt.MapClaims) (string, error) {
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(refreshTokenTTL).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.refreshTokenSecret)
}

func (j *JWT) ValidateAccessToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.accessTokenSecret, nil
	})
	return claims, err
}

func (j *JWT) ValidateRefreshToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return j.refreshTokenSecret, nil
	})
	return claims, err
}
