package install

import (
	"tuck/internal/log"
	"tuck/internal/path"

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
			if !path.Exists(params.Package) {
				log.Fatalln("local package does not exist:", params.Package)
			}
			pkg := path.Abs(params.Package)
			entries := path.Stow(pkg, path.Expand(params.Prefix))
			log.Debugf("stowed %d entries\n", len(entries))
			for _, entry := range entries {
				log.Debugf("  %s\n", entry)
			}
			// TODO: store list of files installed by package
		} else {
			// TODO: release = github.GetRelease(params.Package, params.Release)
			// TODO: asset = selectAsset(release, osInfo)
			// TODO: pkg_archive = downloadAsset(asset)
			// TODO: pkg_dir = archive.Extract(pkg_archive, xdg.DATA_HOME / tuck / params.Package)
			// TODO: stow.Stow(pkg_dir, params.Prefix)
		}
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
