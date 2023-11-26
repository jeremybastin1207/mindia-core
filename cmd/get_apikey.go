package cmd

import (
	"github.com/jeremybastin1207/mindia-core/internal/task"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewApiKeyReaderCommand(aikeyOperator task.ApiKeyOperator) *cobra.Command {
	return &cobra.Command{
		Use:   "apikey list",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			data := pterm.TableData{{"Name", "Key"}}

			apikeys, err := aikeyOperator.GetAll()
			if err != nil {

				return
			}
			for _, apikey := range apikeys {
				data = append(data, []string{
					apikey.Name,
					apikey.Key,
				})
			}
			_ = pterm.DefaultTable.WithHasHeader().WithData(data).Render()
		},
	}
}
