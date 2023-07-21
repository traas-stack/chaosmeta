package app

import (
	"chaosmeta-platform/pkg/common/log"
	"context"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
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
	log.Info(ctx, "start chaosmeta-platform")
	beego.SetStaticPath("/swagger", "swagger")
	beego.Run()
	return nil
}
