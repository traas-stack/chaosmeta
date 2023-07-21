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
	"chaosmeta-platform/pkg/models"
	"chaosmeta-platform/util/errors"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type GrantType string

var (
	GrantTypeAccess  GrantType = "access"
	GrantTypeRefresh GrantType = "refresh"
	Admin                      = "admin"
)

type UserRole string

const (
	NormalRole  = UserRole("normal")
	AdminRole   = UserRole("admin")
	visitorRole = UserRole("visitor")
)

type User struct{}

func (a *User) InitAdmin(ctx context.Context, name, password string) error {
	user, err := a.Get(ctx, Admin)
	if err == nil && user != nil {
		return nil
	}
	return a.Create(ctx, Admin, Admin, string(AdminRole))
}

func (a *User) IsAdmin(ctx context.Context, name, password string) bool {
	user := models.User{Email: name}
	if err := models.GetUser(ctx, &user); err != nil {
		return false
	}
	if user.Disabled {
		return false
	}
	if !VerifyPassword(password, user.Password) {
		return false
	}
	return user.Role == models.AdminRole
}

func (a *User) Login(ctx context.Context, name, password string) (string, string, error) {
	user := models.User{Email: name}
	if err := models.GetUser(ctx, &user); err != nil {
		return "", "", err
	}
	if user.Disabled {
		return "", "", errors.ErrUnauthorized()
	}
	if !VerifyPassword(password, user.Password) {
		return "", "", errors.ErrUnauthorized()
	}

	user.LastLoginTime = time.Now()
	if err := models.UpdateUser(ctx, &user); err != nil {
		return "", "", err
	}

	authentication := Authentication{}
	tocken, err := authentication.GenerateToken(name, string(GrantTypeAccess), 1*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := authentication.GenerateToken(name, string(GrantTypeRefresh), time.Hour*24)
	if err != nil {
		return "", "", err
	}
	return tocken, refreshToken, nil
}

func (a *User) Create(ctx context.Context, name, password, role string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	user := models.User{
		Email:    name,
		Password: hash,
		Role:     role,
		Disabled: false,
	}

	_, err = models.InsertUser(ctx, &user)
	return err
}

func (a *User) Get(ctx context.Context, name string) (*models.User, error) {
	user := models.User{Email: name}
	if err := models.GetUser(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *User) GetList(ctx context.Context, name, role, orderBy string, offset, limit int) (int64, []models.User, error) {
	return models.QueryUser(ctx, name, role, orderBy, offset, limit)
}

func (a *User) DeleteList(ctx context.Context, name, password string, deleteIds []int) error {
	if !a.IsAdmin(ctx, name, password) {
		return fmt.Errorf("not admin")
	}

	if err := models.DeleteUsersByIdList(ctx, deleteIds); err != nil {
		return err
	}
	return models.UsersOrNamespacesDelete(deleteIds, nil)
}

func (a *User) UpdatePassword(ctx context.Context, name, password, newPassword string) error {
	user, err := a.Get(ctx, name)
	if err != nil {
		return err
	}
	if user.Disabled {
		return errors.ErrUnauthorized()
	}
	if !VerifyPassword(password, user.Password) {
		return errors.ErrUnauthorized()
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hash
	return models.UpdateUser(ctx, user)
}

func (a *User) UpdateRole(ctx context.Context, name, password, changeUserName, role string) error {
	if !a.IsAdmin(ctx, name, changeUserName) {
		return fmt.Errorf("not admin")
	}
	user, err := a.Get(ctx, name)
	if err != nil {
		return err
	}
	user.Role = role
	return models.UpdateUser(ctx, user)
}

func (a *User) CheckToken(ctx context.Context, token string) error {
	if token == "" {
		return errors.ErrUnauthorized()
	}
	authentication := Authentication{}
	tokenClaims, err := authentication.VerifyToken(token)
	if err != nil {
		return errors.ErrUnauthorized()
	}
	if tokenClaims.GrantType != string(GrantTypeAccess) {
		return errors.ErrUnauthorized()
	}
	return nil
}

// Generate a user's hashed password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Verify that the user's password is correct
func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println(err)
	return err == nil
}
