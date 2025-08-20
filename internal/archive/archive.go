package archive

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Extract(archive string, outdir string) error {
	switch {
	case strings.HasSuffix(archive, ".tar.gz"):
		return tar("xzf", archive, outdir)
	case strings.HasSuffix(archive, ".tar.xz"):
		return tar("xJf", archive, outdir)
	case strings.HasSuffix(archive, ".tar.bz2"):
		return tar("xjf", archive, outdir)
	default:
		return fmt.Errorf("unsupported archive type: %s", archive)
	}
}

func tar(flags string, archive string, outdir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(outdir)
	defer os.Chdir(cwd)
	if err != nil {
		return err
	}

	cmd := exec.Command("tar", flags, archive)
	err = cmd.Run()
	if err != nil {
		exitErr := err.(*exec.ExitError)
		return fmt.Errorf("extracting '%s' failed:\n%s", archive, exitErr.Stderr)
	}

	return nil
}
