package main

import (
	"chaosmeta-platform/pkg/common/log"
	//_ "chaosmeta-platform/routers"
	"context"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	log.Init()
	//models.Setup(cfg)

}

func main() {
	// TODO: for test
	ctx := context.WithValue(context.Background(), log.TraceIdKey, "erg3g42g432g")

	log.Info(ctx, "start chaosmeta-platform")
	beego.Run()
}
