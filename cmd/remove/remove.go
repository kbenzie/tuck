package remove

import (
	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use: "remove [flags] package",
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Remove an installed package",
	Long: `Remove a package with a local path or from a GitHub release
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
	RemoveCmd.Flags().StringVarP(&prefix, "prefix", "p",
		"~/.local", "install prefix path")
	RemoveCmd.Flags().StringVarP(&version, "tag", "t",
		"latest", "release tag to install")
}
