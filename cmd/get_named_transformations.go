package cmd

import (
	"github.com/jeremybastin1207/mindia-core/internal/task"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewGetNamedTransformationsCommand(namedTransformationOperator task.NamedTransformationOperator) *cobra.Command {
	return &cobra.Command{
		Use:   "named transformation list",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			data := pterm.TableData{{"Name", "Transformation"}}

			namedTransformations, err := namedTransformationOperator.GetAll()
			if err != nil {
				return
			}
			for _, t := range namedTransformations {
				data = append(data, []string{
					t.Name,
					t.Transformations,
				})
			}
			err = pterm.DefaultTable.WithHasHeader().WithData(data).Render()
			if err != nil {
				return
			}
		},
	}
}
