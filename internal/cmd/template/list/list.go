package list

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct{}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List templates",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	return cmd
}

func runList(f *cmdutil.Factory, opts Options) error {
	templates, err := f.ApiClient.ListAllTemplates(context.Background())
	if err != nil {
		return err
	}

	f.Printer.Table(templates.Header(), templates.Rows())

	return nil
}
