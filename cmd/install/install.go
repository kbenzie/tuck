package install

import (
	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install [flags] package",
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Install a local or remote package",
	Long: `Install a package with a local path or from a GitHub release
with a project slug or URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		println("package:", args[0])
		println("prefix:", prefix)
		println("version:", version)
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
