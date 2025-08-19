package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"tuck/cmd/install"
	"tuck/cmd/list"
	"tuck/cmd/remove"
	"tuck/internal/log"
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
			break
		case 1:
			log.SetLevel(log.LevelInfo)
			break
		default:
			log.SetLevel(log.LevelDebug)
			break
		}
	},
}

func init() {
	rootCmd.PersistentFlags().CountVarP(&params.Verbose, "verbose", "v",
		"enable verbose output")
	rootCmd.AddCommand(install.InstallCmd)
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(remove.RemoveCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
