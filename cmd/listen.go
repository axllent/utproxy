package cmd

import (
	"github.com/axllent/utproxy/app"
	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Start daemon",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		app.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)
}
