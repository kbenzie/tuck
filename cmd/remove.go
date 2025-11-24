package cmd

import (
	"fmt"
	"os"
	"tuck/internal/log"
	"tuck/internal/path"
	"tuck/internal/state"

	"github.com/spf13/cobra"
)

var removeParams struct {
	Package string
}

var removeCmd = &cobra.Command{
	Use:   "remove [flags] package",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Remove an installed package",
	Long: `Remove a package with a local path or from a GitHub release
with a project slug or URL.`,
	ValidArgsFunction: removeValidArgsFunc,
	Run: func(cmd *cobra.Command, args []string) {
		removeParams.Package = args[0]
		log.Infof("remove: %+v\n", removeParams)
		pkg, err := state.Get(removeParams.Package)
		if err != nil {
			log.Fatalln(err)
		}
		if pkg == nil {
			log.Errorln("package not installed:", removeParams.Package)
			return
		}
		// remove files
		for _, file := range pkg.Files {
			os.Remove(file)
			log.Infoln("removed:", file)
		}
		state.Remove(removeParams.Package)
		fmt.Printf("tuck removed %d files from '%s' out of '%s'\n",
			len(pkg.Files), path.Contract(removeParams.Package),
			path.Contract(pkg.Prefix))
		// TODO: remove empty directories
	},
}

func removeValidArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var completions []string
	pkgs, err := state.GetAll()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	for name := range *pkgs {
		completions = append(completions, name)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Aliases = append(removeCmd.Aliases, "rm")
}
