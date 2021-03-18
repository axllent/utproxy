/*
Copyright © 2020-Now() Ralph Slooten
This file is part of a CLI application.
*/
package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Prints detailed information about the build environment
and the version of this software.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("utproxy %s compiled with %v on %v/%v\n",
			Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
