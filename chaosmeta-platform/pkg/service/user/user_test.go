package user

import (
	"chaosmeta-platform/config"
	"context"
	"fmt"
	"testing"
)

func init() {
	setUp()
}

func setUp() {
	if err := config.InitConfigWithFilePath("/Users/samson/GolandProjects/chaosmeta/chaosmeta-platform/config"); err != nil {
		panic("config init failed")
	}
	config.Setup()
}

func TestUser_Login(t *testing.T) {
	setUp()
	a := &UserService{}
	got, got1, err := a.Login(context.Background(), "liusongshan.lss@alibaba-inc.com", "samson")

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("token:", got)
	fmt.Println("fresh token:", got1)
}

func TestUser_CheckToken(t *testing.T) {
	a := &UserService{}
	_, err := a.CheckToken(context.Background(), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImxpdXNvbmdzaGFuLmxzc0BhbGliYWJhLWluYy5jb20iLCJncmFudFR5cGUiOiJhY2Nlc3MiLCJleHAiOjE2ODk3NjY5NTcsImlzcyI6ImNoYW9zbWV0YV9pc3N1ZXIiLCJuYmYiOjE2ODk3NjY4OTd9.DniQEMuR5MyaR3beJvM17dm4qdl_wI3Pc93GV1OBKeg")

	if err != nil {
		t.Fatal(err)
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	setUp()
	a := &UserService{}
	err := a.UpdatePassword(context.Background(), "liusongshan.lss@alibaba-inc.com", "samson")

	if err != nil {
		t.Fatal(err)
	}
}

func TestUser_Create(t *testing.T) {
	setUp()
	a := &UserService{}
	_, err := a.Create(context.Background(), "liusongshan.lss@alibaba-inc.com", "samson", string(AdminRole))

	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < 10; i++ {
		a := &UserService{}
		_, err := a.Create(context.Background(), fmt.Sprintf("liusongshan.lss%d@alibaba-inc.com", i), "samson", string(NormalRole))

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestUser_Delete(t *testing.T) {
	a := &UserService{}
	err := a.DeleteList(context.Background(), "liusongshan.lss@alibaba-inc.com", []int{2})

	if err != nil {
		t.Fatal(err)
	}
}

func TestUser_GetList(t *testing.T) {
	a := &UserService{}

	count, usrList, err := a.GetList(context.Background(), "liusongshan", Admin, "create_time", 0, 5)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(count, usrList)
}

func TestUser_Get(t *testing.T) {
	a := &UserService{}
	usr, err := a.Get(context.Background(), "liusongshan.lss@alibaba-inc.com")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(usr)

}
