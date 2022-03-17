package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/axllent/utproxy/updater"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Prints detailed information about the build environment
and the version of this software.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		update, _ := cmd.Flags().GetBool("update")

		if update {
			return updateApp()
		}

		fmt.Printf("utproxy %s compiled with %v on %v/%v\n",
			Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)

		latest, _, _, err := updater.GithubLatest(Repo, RepoBinaryName)
		if err == nil && updater.GreaterThan(latest, Version) {
			fmt.Printf(
				"Update available: %s\nRun `%s version -u` to update (requires read/write access to install directory).\n",
				latest,
				os.Args[0],
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().
		BoolP("update", "u", false, "update to latest version")
}

func updateApp() error {
	rel, err := updater.GithubUpdate(Repo, RepoBinaryName, Version)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %s to version %s\n", os.Args[0], rel)
	return nil
}
