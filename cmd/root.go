package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"tuck/cmd/install"
	"tuck/cmd/remove"
)

var rootCmd = &cobra.Command{
	Use:   "tuck",
	Short: "To fit (packages) securely or snugly (into ~/.local/bin)",
	Long: `To fit (packages) securely or snugly (into ~/.local/bin)

Tuck has the following goals:

* Manage installation of local packages, similar to GNU Stow
* Easily download and install packages from GitHub releases
* Be a single statically compiled binary that's easy to install `,
}

func init() {
	rootCmd.AddCommand(install.InstallCmd)
	rootCmd.AddCommand(remove.RemoveCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
