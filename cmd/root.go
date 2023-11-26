package cmd

import (
	"fmt"
	"os"

	"github.com/jeremybastin1207/mindia-core/internal/task"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mindia-cli",
	Short: "Mindia's cli",
}

func Execute(getMediaTask task.GetMediaTask,
	apiKeyOperator task.ApiKeyOperator,
	namedTransformationOperator task.NamedTransformationOperator) {
	/* 	rootCmd.AddCommand(NewGetMediaCommand(getMediaTask))
	   	rootCmd.AddCommand(NewGetApikeyCommand(namedTransformationOperator))
	   	rootCmd.AddCommand(NewGetNamedTransformationsCommand(apiKeyOperator))
	*/
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
