package app

import (
	"chaosmeta-platform/config"
	"chaosmeta-platform/pkg/service/inject"
	"chaosmeta-platform/pkg/service/namespace"
	"chaosmeta-platform/pkg/service/user"
	"chaosmeta-platform/util/log"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default in $HOME/.chaosmeta/config.yaml)")

	rootCmd.AddCommand(serverCmd)

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(err)
	}
	config.InitConfig()

	log.SetDefaultLogOption(log.LogOption{
		LogPath:    config.DefaultRunOptIns.Log.Path,
		MaxAge:     7,
		MaxSize:    64,
		MaxBackups: 3,
		OutPutType: "BothFileAndStdErrPut",
		Level:      config.DefaultRunOptIns.Log.Level,
	})

	config.Setup()
	user.Init()
	namespace.Init()
	if err := inject.Init(); err != nil {
		log.Panic(err)
	}
	//if err := clientset.Init(); err != nil {
	//	log.Panic(err)
	//}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chaosmeta",
	Short: "c",
	Long:  `混沌工程`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
