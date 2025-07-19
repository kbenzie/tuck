package path

import (
	"os"
	"path/filepath"
	"strings"
	"tuck/internal/log"
)

func Abs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Errorln(err)
	}
	return abs
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Expand(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Errorln(err)
		}
		path = filepath.Join(home, path[1:])
	}
	return path
}

func Contract(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Errorln(err)
	}
	if strings.HasPrefix(path, home) {
		path = "~" + path[len(home):]
	}
	return path
}
