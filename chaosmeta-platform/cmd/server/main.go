package main

import (
	"chaosmeta-platform/cmd/server/app"
	_ "chaosmeta-platform/cmd/server/docs"
)

// @title Chaosmeta
// @version 1.0 （必填）
// @description This is chaosmeta-platform api docs.
// @license.name Apache 2.0
// @contact.name go-swagger帮助文档
// @contact.url https://github.com/traas-stack/chaosmeta
// @host 127.0.0.1:8080

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						CHAOSMETA_PLATFORM_TOKEN
//	@description				用户令牌

// @BasePath /chaosmeta/api/v1
func main() {
	app.Execute()
}
