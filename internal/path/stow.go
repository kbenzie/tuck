package path

import (
	"tuck/internal/log"

	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

func isStdDirLayout(dirs []os.DirEntry) bool {
	for _, entry := range dirs {
		if !entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), "bin") {
			return true
		}
		// TODO: do more in depth checks?
		// strings.HasSuffix(entry.Name(), "lib")
		// strings.HasSuffix(entry.Name(), "share")
	}
	return false
}

func isExecutable(path string) bool {
	// Assumes the file exists and is regular
	if strings.HasPrefix(runtime.GOOS, "windows") {
		ext := strings.ToLower(filepath.Ext(path))
		if slices.Contains([]string{
			".exe",
			".bat",
			".cmd",
			".com",
			".ps1",
		}, ext) {
			return true
		}
	} else {
		stat, _ := os.Stat(path)
		mode := stat.Mode().Perm()
		return mode&0111 != 0
	}
	return false
}

func isManPage(path string) bool {
	// TODO: Also detect ".1.gz"?
	return strings.HasSuffix(path, ".1")
}

func Stow(src string, dst string, dryRun bool) []string {
	stows := []string{}

	entries, err := os.ReadDir(src)
	if err != nil {
		log.Fatalln(err)
	}

	if isStdDirLayout(entries) {
		log.Debugln("detected package content has standard directory layout")

		// FIXME: don't stow everything but files in the root instead only stow
		// files in directories we know about, currently this will install any
		// directory in the archive even when that doesn't make sense, e.g.
		// completions should likely live elsewhere.
		dirs := []string{}
		files := []string{}

		for _, entry := range entries {
			if entry.IsDir() {
				filepath.WalkDir(filepath.Join(src, entry.Name()),
					func(path string, d os.DirEntry, err error) error {
						if d.IsDir() {
							dirs = append(dirs, path)
						} else {
							files = append(files, path)
						}
						return err
					})
			}
		}

		// create directories in dst
		for _, indir := range dirs {
			reldir, err := filepath.Rel(src, indir)
			if err != nil {
				log.Fatalln(err)
			}
			outdir := filepath.Join(dst, reldir)
			if !Exists(outdir) {
				if !dryRun {
					err := os.MkdirAll(outdir, os.ModePerm)
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}

		// move files to dst
		for _, infile := range files {
			relfile, err := filepath.Rel(src, infile)
			if err != nil {
				log.Fatalln(err)
			}
			outfile := filepath.Join(dst, relfile)
			if !dryRun {
				err = os.Rename(infile, outfile)
				if err != nil {
					log.Fatalln(err)
				}
			}
			stows = append(stows, outfile)
		}

	} else {
		log.Debugln("detected package content has non-standard directory layout")

		// recursively enumerate entries in src path
		dirs := []string{}
		files := []string{}
		filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
			if d.IsDir() {
				dirs = append(dirs, path)
			} else {
				files = append(files, path)
			}
			return err
		})

		// copy files of interest but not docs/license/etc
		bins := []string{}
		manpages := []string{}
		for _, file := range files {
			if isExecutable(file) {
				bins = append(bins, file)
			} else if isManPage(file) {
				manpages = append(manpages, file)
			}
		}

		// TODO: completions := []string{}
		// for _, file := range files {
		//	if isCompletionScriptFor(bins, file) {
		//		completions = append(completions, file)
		//	}
		// }

		if len(bins) > 0 {
			binDir := filepath.Join(dst, "bin")
			if !dryRun {
				err := os.MkdirAll(binDir, os.ModePerm)
				if err != nil {
					log.Fatalln(err)
				}
			}
			for _, inbin := range bins {
				outbin := filepath.Join(binDir, filepath.Base(inbin))
				if !dryRun {
					err := os.Rename(inbin, outbin)
					if err != nil {
						log.Fatalln(err)
					}
				}
				stows = append(stows, outbin)
			}
		}

		if len(manpages) > 0 {
			// TODO: Don't assume man1 directory
			manDir := filepath.Join(dst, "share/man/man1")
			if !dryRun {
				err = os.MkdirAll(manDir, os.ModePerm)
				if err != nil {
					log.Fatalln(err)
				}
			}
			for _, src := range manpages {
				outbin := filepath.Join(manDir, filepath.Base(src))
				if !dryRun {
					err := os.Rename(src, outbin)
					if err != nil {
						log.Fatalln(err)
					}
				}
				stows = append(stows, outbin)
			}
		}
	}

	return stows
}
