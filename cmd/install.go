package cmd

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

var installParams struct {
	Package string
	Prefix  string
	Release string
	Local   bool
	DryRun  bool
}

var installCmd = &cobra.Command{
	Use:   "install [flags] package",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Install a local or remote package",
	Long: `Install a package with a local path or from a GitHub release
with a project slug or URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		installParams.Package = args[0]
		log.Debugf("install: %+v\n", installParams)

		installParams.Prefix = path.Abs(path.Expand(installParams.Prefix))
		files := []string{}

		// TODO: create a file lock, defer its deletion, do the same for other
		// commands which mutate the filesystem state
		cfg, err := config.Load()
		if err != nil {
			log.Fatalln(err)
		}
		log.Debugln(cfg)

		if installParams.Local {
			if !path.Exists(installParams.Package) {
				log.Fatalln("local package does not exist:", installParams.Package)
			}
			installParams.Package = path.Abs(installParams.Package)

			// check if a similar package has already been installed?
			pkg, _ := state.Get(installParams.Package)
			if pkg != nil {
				log.Fatalf("package already installed: '%s'\n", installParams.Package)
			}

			// TODO: link instead of move for local packages
			files = path.Stow(installParams.Package, installParams.Prefix, installParams.DryRun)
		} else {
			// TODO: check if a similar package has already been installed?

			release, err := github.GetRelease(installParams.Package, installParams.Release)
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

			files = path.Stow(dir, installParams.Prefix, installParams.DryRun)

			// remove the package directory, cache dir may no longer exist
			os.RemoveAll(dir)
		}

		for _, file := range files {
			log.Infoln("installed:", file)
		}
		fmt.Printf("tuck installed %d files from '%s' into '%s'\n",
			len(files), path.Contract(installParams.Package),
			path.Contract(installParams.Prefix))

		if !installParams.DryRun {
			// store list of files installed by package
			state.Install(installParams.Package, state.Package{
				Prefix:  installParams.Prefix,
				Release: installParams.Release,
				Local:   installParams.Local,
				Files:   files,
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Aliases = append(installCmd.Aliases, "in")
	installCmd.Flags().StringVarP(&installParams.Prefix, "prefix", "p",
		"~/.local", "install prefix path")
	installCmd.Flags().StringVarP(&installParams.Release, "release", "r",
		"latest", "github release to install")
	installCmd.Flags().BoolVarP(&installParams.Local, "local", "l", false,
		"treat package as local path")
	installCmd.Flags().BoolVarP(&installParams.DryRun, "dry-run", "d", false,
		"don't actually install anything")
}
