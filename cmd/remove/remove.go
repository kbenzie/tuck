package remove

import (
	"tuck/internal/log"

	"github.com/spf13/cobra"
)

var params struct {
	Prefix  string
	Local   bool
	Package string
}

var RemoveCmd = &cobra.Command{
	Use:   "remove [flags] package",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Remove an installed package",
	Long: `Remove a package with a local path or from a GitHub release
with a project slug or URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		params.Package = args[0]
		log.Infof("remove: %+v\n", params)
	},
}

func init() {
	RemoveCmd.Flags().StringVarP(&params.Prefix, "prefix", "p",
		"~/.local", "install prefix path")
	RemoveCmd.Flags().BoolVarP(&params.Local, "local", "l", false,
		"treat package as local path")
}
