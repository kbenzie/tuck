package cmd

import (
	"fmt"
	"tuck/internal/log"
	"tuck/internal/state"

	"github.com/spf13/cobra"
)

var listParams struct {
	Quiet bool
}

var listCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.MatchAll(cobra.OnlyValidArgs),
	Short: "List installed packages",
	Long:  `List installed packages.`,
	Run: func(cmd *cobra.Command, args []string) {
		pkgs, err := state.GetAll()
		if err != nil {
			log.Fatalln(err)
		}
		if !listParams.Quiet {
			fmt.Println(len(*pkgs), "packages are installed")
		}
		for name, _ := range *pkgs {
			fmt.Printf("%s\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listParams.Quiet, "--quiet", "q", false,
		"list only package names, nothing else")
}
