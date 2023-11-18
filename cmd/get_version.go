package cmd

import (
	"fmt"

	"github.com/jeremybastin1207/mindia-core/internal/settings"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Mindia",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Mindia", settings.Version)
	},
}
