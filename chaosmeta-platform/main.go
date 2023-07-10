package main

import (
	"chaosmeta-platform/logger"
	"chaosmeta-platform/models"
	_ "chaosmeta-platform/routers"
	"context"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	logger.Setup()
	models.Setup()
}

func main() {
	ctx := context.WithValue(context.Background(), logger.TraceIdKey, "erg3g42g432g")
	fmt.Println(models.UpdateUserRole(ctx, 1, models.AdminRole))
	fmt.Println(models.UpdateUserRole(ctx, 3, models.AdminRole))

	//password := "123456"
	//salt := "abcd"
	//hash := sha256.Sum256([]byte(password + salt))
	//hashedPassword := hex.EncodeToString(hash[:])

	beego.Run()
}
