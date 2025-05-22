package install

import (
	"github.com/spf13/cobra"
	"tuck/internal/log"
)

var InstallCmd = &cobra.Command{
	Use:   "install [flags] package",
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Install a local or remote package",
	Long: `Install a package with a local path or from a GitHub release
with a project slug or URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		var pkg string = args[0]
		log.Info("package:", pkg)
		log.Info("prefix:", prefix)
		log.Info("version:", version)

		// TODO: get latest release from github api
	},
}

var prefix string
var version string

func init() {
	InstallCmd.Aliases = append(InstallCmd.Aliases, "in")
	InstallCmd.Flags().StringVarP(&prefix, "prefix", "p",
		"~/.local", "install prefix path")
	InstallCmd.Flags().StringVarP(&version, "tag", "t",
		"latest", "release tag to install")
}
