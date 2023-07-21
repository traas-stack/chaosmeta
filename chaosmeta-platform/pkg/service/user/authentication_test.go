package user

import (
	"chaosmeta-platform/config"
	"fmt"
	"testing"
	"time"
)

func init() {
	config.DefaultRunOptIns = &config.Config{
		SecretKey: "samson",
	}
}

func TestAuthentication_VerifyToken(t *testing.T) {
	authentication := Authentication{}
	tocken, err := authentication.GenerateToken("samson", "admin", 1*time.Minute)
	fmt.Println(tocken)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := authentication.VerifyToken(tocken)
	fmt.Sprintln(*claims, err)
}
