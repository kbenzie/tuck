package install

import (
	"os"
	"tuck/internal/log"

	"github.com/spf13/cobra"
)

var params struct {
	Prefix  string
	Release string
	Local   bool
	Package string
}

var InstallCmd = &cobra.Command{
	Use:   "install [flags] package",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Install a local or remote package",
	Long: `Install a package with a local path or from a GitHub release
with a project slug or URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		params.Package = args[0]
		log.Infof("install: %+v\n", params)
		if params.Local {
			if _, err := os.Stat(params.Package); os.IsNotExist(err) {
				log.Fatalln("local path does not exist:", params.Package)
			}
		} else {
			// TODO: download = github.DownloadRelease(params.Package, params.Release)
			// TODO: extract
		}
		// TODO: stow
	},
}

func init() {
	InstallCmd.Aliases = append(InstallCmd.Aliases, "in")
	InstallCmd.Flags().StringVarP(&params.Prefix, "prefix", "p",
		"~/.local", "install prefix path")
	InstallCmd.Flags().StringVarP(&params.Release, "release", "r",
		"latest", "github release to install")
	InstallCmd.Flags().BoolVarP(&params.Local, "local", "l", false,
		"treat package as local path")
}
