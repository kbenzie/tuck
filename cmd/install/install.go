package install

import (
	"fmt"
	"tuck/internal/log"
	"tuck/internal/path"
	"tuck/internal/state"

	"github.com/spf13/cobra"
)

var params struct {
	Package string
	Prefix  string
	Release string
	Local   bool
	DryRun  bool
}

var InstallCmd = &cobra.Command{
	Use:   "install [flags] package",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Install a local or remote package",
	Long: `Install a package with a local path or from a GitHub release
with a project slug or URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		params.Package = args[0]
		log.Debugf("install: %+v\n", params)

		params.Prefix = path.Abs(path.Expand(params.Prefix))
		files := []string{}

		if params.Local {
			if !path.Exists(params.Package) {
				log.Fatalln("local package does not exist:", params.Package)
			}
			params.Package = path.Abs(params.Package)

			// check if a similar package has already been installed?
			pkg, _ := state.Get(params.Package)
			if pkg != nil {
				log.Fatalf("package already installed: '%s'\n", params.Package)
			}

			files = path.Stow(params.Package, params.Prefix, params.DryRun)
			for _, file := range files {
				log.Infoln("installed:", file)
			}

			fmt.Printf("tuck installed %d files from '%s' into '%s'\n",
				len(files), path.Contract(params.Package),
				path.Contract(params.Prefix))
		} else {
			// TODO: check if a similar package has already been installed?

			// TODO: release = github.GetRelease(params.Package, params.Release)
			// TODO: asset = selectAsset(release, osInfo)
			// TODO: pkg_archive = downloadAsset(asset)
			// TODO: pkg_dir = archive.Extract(pkg_archive, xdg.DATA_HOME / tuck / params.Package)
			// TODO: stow.Stow(pkg_dir, params.Prefix)
		}

		if !params.DryRun {
			// store list of files installed by package
			state.Install(params.Package, state.Package{
				Prefix:  params.Prefix,
				Release: params.Release,
				Local:   params.Local,
				Files:   files,
			})
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
	InstallCmd.Flags().BoolVarP(&params.DryRun, "--dry-run", "d", false,
		"don't actually install anything")
}
