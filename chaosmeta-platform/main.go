package main

import (
	"chaosmeta-platform/models"
	"chaosmeta-platform/pkg/logger"
	_ "chaosmeta-platform/routers"
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	cfg, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		panic(any(fmt.Sprintf("init config error: %s", err.Error())))
	}

	logger.Setup(cfg)
	models.Setup(cfg)
}

func main() {
	// TODO: for test
	ctx := context.WithValue(context.Background(), logger.TraceIdKey, "erg3g42g432g")
	//password := "123456"
	//salt := "abcd"
	//hash := sha256.Sum256([]byte(password + salt))
	//hashedPassword := hex.EncodeToString(hash[:])

	id, err := models.InsertUser(ctx, &models.User{
		Name:     "test",
		Password: "test",
		Role:     models.AdminRole,
	})
	if err != nil {
		panic(any(fmt.Sprintf("insert user error: %s", err.Error())))
	}
	fmt.Println(models.UpdateUserRole(ctx, id, models.NormalRole))
	logger.Info(ctx, "start chaosmeta-platform")

	beego.Run()
}
