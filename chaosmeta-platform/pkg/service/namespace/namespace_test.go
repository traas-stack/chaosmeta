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

package namespace

import (
	"chaosmeta-platform/config"
	"chaosmeta-platform/pkg/models/namespace"
	"context"
	"fmt"
	"testing"
)

func init() {
	if err := config.InitConfigWithFilePath("/Users/samson/GolandProjects/chaosmeta/chaosmeta-platform/conf"); err != nil {
		panic(err)
	}
	config.Setup()
}

func TestNamespaceService_CreateNamespace(t *testing.T) {
	s := &NamespaceService{}
	if _, err := s.Create(context.Background(), "冷凇的空间", "一次重要的里程碑", "liusongshan.lss@alibaba-inc.com"); err != nil {
		t.Fatal("CreateNamespace() error", err)
	}
	for i := 0; i < 10; i++ {
		if _, err := s.Create(context.Background(), fmt.Sprintf("冷凇的空间:%d", i), "一次重要的里程碑", "liusongshan.lss@alibaba-inc.com"); err != nil {
			t.Fatal("CreateNamespace() error", err)
		}
	}
}

func TestNamespaceService_UpdateNamespace(t *testing.T) {
	s := &NamespaceService{}
	if err := s.Update(context.Background(), "", 1, "冷凇的新空间", "一次新重要的里程碑"); err != nil {
		t.Fatal("CreateNamespace() error", err)
	}
}

func TestNamespaceService_GetNamespace(t *testing.T) {
	s := &NamespaceService{}
	n, err := s.Get(context.Background(), 2)
	if err != nil {
		t.Fatal("CreateNamespace() error", err)
	}
	fmt.Println(n)
}

func TestNamespaceService_DeleteNamespace(t *testing.T) {
	s := &NamespaceService{}
	err := s.Delete(context.Background(), "", 1)
	if err != nil {
		t.Fatal("CreateNamespace() error", err)
	}
}

func TestNamespaceService_AddUsers(t *testing.T) {
	s := &NamespaceService{}
	err := s.AddUsers(context.Background(), "", 1, namespace.AddUsersParam{})
	if err != nil {
		t.Fatal("CreateNamespace() error", err)
	}
}

func TestNamespaceService_GroupedUserNamespaces(t *testing.T) {
	s := &NamespaceService{}
	total, userList, err := s.GroupedUserInNamespaces(context.Background(), 1, "", "", -1, "create_time", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("total", total)
	fmt.Println("userList", userList)
}

func TestNamespaceService_GroupUserNotInNamespaces(t *testing.T) {
	s := &NamespaceService{}
	total, namespaceList, err := s.GroupNamespacesUserNotIn(context.Background(), 3, "", "", "-create_time", 1, 100)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("total", total)
	fmt.Println("namespaceList", namespaceList)
}

func TestNamespaceService_GroupNamespacesByUsername(t *testing.T) {
	s := &NamespaceService{}
	total, namespaceList, err := s.GroupNamespacesByUsername(context.Background(), -1, "liusongshan", "", 0, "-create_time", 1, 100)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("total", total)
	fmt.Println("namespaceList", namespaceList)
}

func TestNamespaceService_GroupAllNamespaces(t *testing.T) {
	s := &NamespaceService{}
	total, namespaceList, err := s.GroupAllNamespaces(context.Background(), "", "liusongshan", "create_time", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("total", total)
	fmt.Println("namespaceList", namespaceList)
}

func TestNamespaceService_RemoveUsers(t *testing.T) {
	s := &NamespaceService{}
	err := s.RemoveUsers(context.Background(), "", []int{2}, 1)
	if err != nil {
		t.Fatal("CreateNamespace() error", err)
	}
}

func TestNamespaceService_ChangeUsersPermission(t *testing.T) {
	s := &NamespaceService{}
	err := s.ChangeUsersPermission(context.Background(), "", []int{2, 3, 4}, 1, namespace.AdminPermission)
	if err != nil {
		t.Fatal("CreateNamespace() error", err)
	}
}
