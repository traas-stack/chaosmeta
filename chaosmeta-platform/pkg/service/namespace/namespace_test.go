package namespace

import (
	"chaosmeta-platform/config"
	"chaosmeta-platform/pkg/models/namespace"
	"context"
	"fmt"
	"testing"
)

func init() {
	if err := config.InitConfigWithFilePath("/Users/samson/GolandProjects/chaosmeta/chaosmeta-platform/config"); err != nil {
		panic("config init failed")
	}
	config.Setup()
}

func TestNamespaceService_CreateNamespace(t *testing.T) {
	s := &NamespaceService{}
	if err := s.Create(context.Background(), "冷凇的空间", "一次重要的里程碑", "liusongshan.lss@alibaba-inc.com"); err != nil {
		t.Fatal("CreateNamespace() error", err)
	}
	for i := 0; i < 10; i++ {
		if err := s.Create(context.Background(), fmt.Sprintf("冷凇的空间:%d", i), "一次重要的里程碑", "liusongshan.lss@alibaba-inc.com"); err != nil {
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

func TestNamespaceService_GetUsers(t *testing.T) {
	s := &NamespaceService{}
	users, count, err := s.GetUsers(context.Background(), 0, "liusongshan.lss", -1, "create_time", 0, 5)
	if err != nil {
		t.Fatal("CreateNamespace() error", err)
	}
	fmt.Println(users, count)
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
