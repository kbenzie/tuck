package path

import (
	"os"
	"path/filepath"
	"strings"
	"tuck/internal/log"

	"github.com/adrg/xdg"
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
		path = filepath.Join(xdg.Home, path[1:])
	}
	return path
}

func Contract(path string) string {
	if strings.HasPrefix(path, xdg.Home) {
		path = "~" + path[len(xdg.Home):]
	}
	return path
}
