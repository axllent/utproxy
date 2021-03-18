package cmd

import (
	"fmt"
	"sort"

	"github.com/axllent/utproxy/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [test1] [test2]",
	Short: "Manually run tests in your config",
	Long: `Manually run test in your configuration file, outputs to console.
If no args are supplied then all tests are run.`,
	Aliases: []string{"check"},
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		tests := args

		if len(args) == 0 {
			s := viper.GetStringMap("services")
			for k := range s {
				tests = append(tests, k)
			}
			sort.Strings(tests)
		}

		for _, key := range tests {
			err := app.Check(key)
			if err != nil {
				fmt.Printf("%-15s %v\n", key, err)
			} else {
				fmt.Printf("%-15s up\n", key)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
