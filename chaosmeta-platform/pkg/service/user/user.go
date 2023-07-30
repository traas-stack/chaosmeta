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
	namespace2 "chaosmeta-platform/pkg/models/namespace"
	"chaosmeta-platform/pkg/service/namespace"
	"chaosmeta-platform/util/errors"
	"chaosmeta-platform/util/log"
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
	NormalRole = UserRole("normal")
	AdminRole  = UserRole("admin")
)

func Init() {
	us := UserService{}
	ctx := context.Background()

	_, err := us.Get(ctx, "admin")
	if err == nil {
		log.Error(err)
		return
	}
	_, err = us.Create(ctx, "admin", "admin", string(AdminRole))
	if err != nil {
		log.Error(err)
	}
}

type UserService struct{}

func (a *UserService) InitAdmin(ctx context.Context, name, password string) error {
	user, err := a.Get(ctx, Admin)
	if err == nil && user != nil {
		return nil
	}
	_, err = a.Create(ctx, Admin, Admin, string(AdminRole))
	return err
}

func (a *UserService) IsAdmin(ctx context.Context, name string) bool {
	user := models.User{Email: name}
	if err := models.GetUser(ctx, &user); err != nil {
		return false
	}
	if user.Disabled {
		return false
	}
	return user.Role == models.AdminRole
}

func (a *UserService) Login(ctx context.Context, name, password string) (string, string, error) {
	user := models.User{Email: name}
	if err := models.GetUser(ctx, &user); err != nil {
		return "", "", err
	}
	if user.Disabled || user.IsDeleted {
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
	tocken, err := authentication.GenerateToken(name, string(GrantTypeAccess), 5*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := authentication.GenerateToken(name, string(GrantTypeRefresh), time.Hour*24)
	if err != nil {
		return "", "", err
	}
	return tocken, refreshToken, nil
}

func (a *UserService) Create(ctx context.Context, name, password, role string) (int, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return 0, err
	}

	userGet, err := a.Get(ctx, name)
	if err == nil && userGet.IsDeleted {
		userGet.IsDeleted = false
		userGet.Password = hash
		return userGet.ID, models.UpdateUser(ctx, userGet)
	}

	user := models.User{
		Email:     name,
		Password:  hash,
		Role:      role,
		Disabled:  false,
		IsDeleted: false,
	}

	_, err = models.InsertUser(ctx, &user)
	if err != nil {
		return 0, err
	}

	us := namespace.NamespaceService{}
	return user.ID, us.DefaultAddUsers(ctx, namespace2.AddUsersParam{
		Users: []namespace2.UserData{{
			Id:         user.ID,
			Permission: int(namespace2.AdminPermission),
		}},
	})
}

func (a *UserService) Get(ctx context.Context, name string) (*models.User, error) {
	user := models.User{Email: name}
	if err := models.GetUser(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *UserService) GetList(ctx context.Context, name, role, orderBy string, page, pageSize int) (int64, []models.User, error) {
	return models.QueryUser(ctx, name, role, orderBy, page, pageSize)
}

func (a *UserService) DeleteList(ctx context.Context, name string, deleteIds []int) error {
	if !a.IsAdmin(ctx, name) {
		return fmt.Errorf("not admin")
	}

	if err := models.DeleteUsersByIdList(ctx, deleteIds); err != nil {
		return err
	}
	return namespace2.UsersOrNamespacesDelete(deleteIds, nil)
}

func (a *UserService) UpdatePassword(ctx context.Context, name, newPassword string) error {
	user, err := a.Get(ctx, name)
	if err != nil {
		return err
	}
	if user.Disabled {
		return errors.ErrUnauthorized()
	}
	if !a.IsAdmin(ctx, name) && name != user.Email {
		return errors.ErrUnauthorized()
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hash
	return models.UpdateUser(ctx, user)
}

func (a *UserService) UpdateListRole(ctx context.Context, name string, ids []int, role string) error {
	if !a.IsAdmin(ctx, name) {
		return fmt.Errorf("not admin")
	}

	return models.UpdateUsersRole(ctx, ids, role)
}

func (a *UserService) CheckToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", errors.ErrUnauthorized()
	}
	authentication := Authentication{}
	tokenClaims, err := authentication.VerifyToken(token)
	if err != nil {
		return "", errors.ErrUnauthorized()
	}
	if tokenClaims.GrantType != string(GrantTypeAccess) {
		return "", errors.ErrUnauthorized()
	}
	return tokenClaims.Username, nil
}

func (a *UserService) RefreshToken(ctx context.Context, token string) (string, error) {
	authentication := Authentication{}
	return authentication.RefreshToken(token, string(GrantTypeAccess))
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
