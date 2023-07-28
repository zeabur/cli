package list

import (
	"context"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List projects, use --page-size to specify page size",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	return cmd
}

// Note: don't import other packages directly in this function, or it will be hard to mock and test
// If you want to add new dependencies, please add them in the Options struct

// runList will list all projects page by page
func runList(f *cmdutil.Factory, opts Options) error {
	projects, err := f.ApiClient.ListAllProjects(context.Background())
	if err != nil {
		return err
	}

	f.Printer.Table(projects.Header(), projects.Rows())

	return nil
}
