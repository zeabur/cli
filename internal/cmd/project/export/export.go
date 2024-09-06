package export

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct{}

func NewCmdExport(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export projects to Template Resource YAML",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(f, opts)
		},
	}

	return cmd
}

// runExport will list all projects page by page
func runExport(f *cmdutil.Factory, opts Options) error {
	projects, err := f.ApiClient.ListAllProjects(context.Background())
	if err != nil {
		return err
	}

	f.Printer.Table(projects.Header(), projects.Rows())

	return nil
}
