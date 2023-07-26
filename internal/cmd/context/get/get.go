package get

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	return cmd
}

func runGet(f *cmdutil.Factory, opts *Options) error {
	project := f.Config.GetContext().GetProject()
	environment := f.Config.GetContext().GetEnvironment()
	service := f.Config.GetContext().GetService()

	header := []string{"Context", "Name", "ID"}
	data := [][]string{
		{"Project", project.GetName(), project.GetID()},
		{"Environment", environment.GetName(), environment.GetID()},
		{"Service", service.GetName(), service.GetID()},
	}

	for _, line := range data {
		if line[1] == "" {
			line[1] = "<not set>"
		}
		if line[2] == "" {
			line[2] = "<not set>"
		}
	}

	f.Printer.Table(header, data)

	return nil
}
