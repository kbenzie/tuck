package cmd

import (
	"fmt"
	"os"
	"tuck/internal/log"

	"github.com/spf13/cobra"
)

var params struct {
	Verbose int
}

var rootCmd = &cobra.Command{
	Use:   "tuck",
	Short: "To fit (packages) securely or snugly (into ~/.local/bin)",
	Long: `To fit (packages) securely or snugly (into ~/.local/bin)

Tuck has the following goals:

* Manage installation of local packages, similar to GNU Stow
* Easily download and install packages from GitHub releases
* Be a single statically compiled binary that's easy to install `,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		switch params.Verbose {
		case 0:
			log.SetLevel(log.LevelWarn)
		case 1:
			log.SetLevel(log.LevelInfo)
		default:
			log.SetLevel(log.LevelDebug)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().CountVarP(&params.Verbose, "verbose", "v",
		"enable verbose output")
}

func SetVersion(version string, commit string, date string) {
	rootCmd.Version = fmt.Sprintf("%s (%s - %s)", version, commit, date)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
