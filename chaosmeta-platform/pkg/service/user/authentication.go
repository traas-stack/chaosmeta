/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package user

import (
	"chaosmeta-platform/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

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

func (a *Authentication) RefreshToken(token, grantType string) (string, error) {
	claim, err := a.VerifyToken(token)
	if err != nil {
		return "", err
	}

	newAccessToken, err := a.GenerateToken(claim.Username, grantType, 10*time.Minute)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}
