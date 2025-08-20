package path

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tuck/internal/log"

	"github.com/adrg/xdg"
)

var (
	CacheDir string = filepath.Join(xdg.CacheHome, "tuck")
	StateDir string = filepath.Join(xdg.StateHome, "tuck")
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

func IsDir(path string) bool {
	info, err := os.Stat(StateDir)
	if !os.IsNotExist(err) {
		return false
	}
	if info == nil || !info.IsDir() {
		return false
	}
	return true
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

func DownloadFile(url string, outpath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error downloading '%s': %d", url, response.StatusCode)
	}
	outfile, err := os.Create(outpath)
	log.Debugf("created file '%s'", outpath)
	if err != nil {
		return err
	}
	defer outfile.Close()
	bytes, err := io.Copy(outfile, response.Body)
	log.Debugf("%d bytes written to '%s'", bytes, outpath)
	return err
}

func init() {
	if !Exists(CacheDir) {
		if err := os.MkdirAll(CacheDir, os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}
	if !Exists(StateDir) {
		if err := os.MkdirAll(StateDir, os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}
}
