package user

import (
	"chaosmeta-platform/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AuthenticationController interface {
	GenerateToken(username, grantType string, expireDuration time.Duration) (string, error)
	VerifyToken(tokenString string) (*Claims, error)
	RefreshToken(refreshToken, grantType string) (string, error)
}

type Authentication struct{}

type Claims struct {
	Username  string `json:"username"`
	GrantType string `json:"grantType"`
	jwt.StandardClaims
}

func (a *Authentication) GenerateToken(username, grantType string, expireDuration time.Duration) (string, error) {
	expire := time.Now().Add(expireDuration)
	claims := &Claims{
		Username:  username,
		GrantType: grantType,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: expire.Unix(),
			Issuer:    "chaosmeta_issuer",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.DefaultRunOptIns.SecretKey))
}

func (a *Authentication) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.DefaultRunOptIns.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}

func (a *Authentication) RefreshToken(refreshToken, grantType string) (string, error) {
	claim, err := a.VerifyToken(refreshToken)
	if err != nil {
		return "", err
	}

	newAccessToken, err := a.GenerateToken(claim.Username, grantType, time.Hour)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}
