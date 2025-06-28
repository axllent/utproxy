package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/axllent/ghru/v2"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Prints detailed information about the build environment
and the version of this software.`,
	Run: func(cmd *cobra.Command, args []string) {
		ghruConf := ghru.Config{
			Repo:           "axllent/utproxy",
			ArchiveName:    "utproxy_{{.OS}}_{{.Arch}}",
			BinaryName:     "utproxy",
			CurrentVersion: Version,
		}

		update, _ := cmd.Flags().GetBool("update")

		if update {
			// Update the app
			rel, err := ghruConf.SelfUpdate()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Printf("Updated %s to version %s\n", os.Args[0], rel.Tag)
			fmt.Println("Remember to restart your application.")
			os.Exit(0)
		}

		fmt.Printf("utproxy %s compiled with %v on %v/%v\n",
			Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)

		release, err := ghruConf.Latest()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// The latest version is the same version
		if release.Tag == Version {
			os.Exit(0)
		}

		// A newer release is available
		fmt.Printf(
			"Update available: %s\nRun `%s version -u` to update (requires read/write access to install directory).\n",
			release.Tag,
			os.Args[0],
		)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().
		BoolP("update", "u", false, "update to latest version")
}
