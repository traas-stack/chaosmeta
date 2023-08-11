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

package app

import (
	"chaosmeta-platform/routers"
	"chaosmeta-platform/util/log"
	"context"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/spf13/cobra"
)

func init() {
	//cobra.OnInitialize(initConfig)
	startCmd.Flags().BoolP("runModel", "m", true, "运行模式")
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动 http 服务",
	Long:  `启动服务`,
	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

func server() {
	if err := Run(); err != nil {
		log.Error(err)
		panic(err)
	}
}

func Run() (err error) {
	ctx := context.WithValue(context.Background(), log.TraceIdKey, "erg3g42g432g")
	log.CtxInfof(ctx, "start chaosmeta-platform")
	beego.SetStaticPath("/swagger", "swagger")
	routers.Init()
	beego.Run()
	return nil
}
