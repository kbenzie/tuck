package remove

import (
	"fmt"
	"os"
	"tuck/internal/log"
	"tuck/internal/path"
	"tuck/internal/state"

	"github.com/spf13/cobra"
)

var params struct {
	Package string
}

var RemoveCmd = &cobra.Command{
	Use:   "remove [flags] package",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Remove an installed package",
	Long: `Remove a package with a local path or from a GitHub release
with a project slug or URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		params.Package = args[0]
		log.Infof("remove: %+v\n", params)
		pkg, err := state.Get(params.Package)
		if err != nil {
			log.Fatalln(err)
		}
		// remove files
		for _, file := range pkg.Files {
			os.Remove(file)
			log.Infoln("removed:", file)
		}
		state.Remove(params.Package)
		fmt.Printf("tuck removed %d files from '%s' out of '%s'\n",
			len(pkg.Files), path.Contract(params.Package),
			path.Contract(pkg.Prefix))
		// TODO: remove empty directories
	},
}

func init() {
	RemoveCmd.Aliases = append(RemoveCmd.Aliases, "rm")
}
