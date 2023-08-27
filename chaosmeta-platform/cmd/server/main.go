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
