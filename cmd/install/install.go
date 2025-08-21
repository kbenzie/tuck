package install

import (
	"fmt"
	"os"
	"path/filepath"
	"tuck/internal/archive"
	"tuck/internal/config"
	"tuck/internal/github"
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

		// TODO: create a file lock, defer its deletion, do the same for other
		// commands which mutate the filesystem state
		cfg, err := config.Load()
		if err != nil {
			log.Fatalln(err)
		}
		log.Debugln(cfg)

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

			// TODO: link instead of move for local packages
			files = path.Stow(params.Package, params.Prefix, params.DryRun)
		} else {
			// TODO: check if a similar package has already been installed?

			release, err := github.GetRelease(params.Package, params.Release)
			if err != nil {
				log.Fatalln(err)
			}
			asset, err := github.SelectAsset(release, cfg.Filters)
			if err != nil {
				log.Fatalln(err)
			}

			archivePath := filepath.Join(path.CacheDir, asset.Name)
			// TODO: validate checksum
			err = path.DownloadFile(asset.BrowserDownloadUrl, archivePath)
			if err != nil {
				log.Fatalln(err)
			}

			err = archive.Extract(archivePath, path.CacheDir)
			if err != nil {
				log.Fatalln(err)
			}

			if err := os.Remove(archivePath); err != nil {
				log.Fatalln(err)
			}

			// TODO: detect if ~/.cache/tuck contains a 1 directory or multiple
			// entries
			entries, err := os.ReadDir(path.CacheDir)
			if err != nil {
				log.Fatalln(err)
			}
			dir := ""
			if len(entries) == 1 && entries[0].IsDir() {
				dir = filepath.Join(path.CacheDir, entries[0].Name())
			} else {
				// assume the archive didn't contain a root directory, this
				// relies on the cache directory being empty
				dir = path.CacheDir
			}

			files = path.Stow(dir, params.Prefix, params.DryRun)

			// remove the package directory, cache dir may no longer exist
			os.RemoveAll(dir)
		}

		for _, file := range files {
			log.Infoln("installed:", file)
		}
		fmt.Printf("tuck installed %d files from '%s' into '%s'\n",
			len(files), path.Contract(params.Package),
			path.Contract(params.Prefix))

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
