package app

import (
	"github.com/spf13/cobra"
)

func init() {
	serverCmd.AddCommand(startCmd)
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "混沌工程统一入口",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
