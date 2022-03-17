package cmd

import (
	"fmt"
	"os"

	"github.com/axllent/utproxy/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Config file
	cfgFile string

	// Version of the app
	Version = "dev"

	// Repo on Github for updater
	Repo = "axllent/utproxy"

	// RepoBinaryName on Github for updater
	RepoBinaryName = "utproxy"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "utproxy",
	Short: `UTProxy - service proxy for uptime monitors

Documentation: https://github.com/axllent/utproxy`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		"",
		"config file (default is /etc/utproxy.yaml)",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".up" (without extension).
		viper.AddConfigPath("/etc")
		viper.SetConfigName("utproxy")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if viper.ConfigFileUsed() == "" {
			fmt.Println("Config file not found.")
			fmt.Println("Save your config file in `/etc/utproxy.yaml` or use the `-c` option.")
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	if len(viper.GetStringMap("Services")) == 0 {
		fmt.Println("No services defined in", viper.ConfigFileUsed())
		os.Exit(1)
	}

	app.LogRequestsTo = viper.GetString("Log")

	viper.WatchConfig() // reload config on change
}
