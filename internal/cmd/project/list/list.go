package list

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	PageSize int
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

	// todo: short flag?
	cmd.Flags().IntVar(&opts.PageSize, "page-size", 5, "Page size")

	return cmd
}

// Note: don't import other packages directly in this function, or it will be hard to mock and test
// If you want to add new dependencies, please add them in the Options struct

// runList will list all projects page by page
// if interactive, it will ask if you want to continue to the next page
// if non-interactive, it will list all projects in the one time
func runList(f *cmdutil.Factory, opts Options) error {
	skip := 0
	next := true

	var projects model.Projects

	firstPage := true

	for next {
		projectCon, err := f.ApiClient.ListProjects(context.Background(), skip, opts.PageSize)
		if err != nil {
			return err
		}
		for _, project := range projectCon.Edges {
			projects = append(projects, project.Node)
		}

		skip += opts.PageSize
		next = projectCon.PageInfo.HasNextPage

		if f.Interactive {
			if firstPage {
				firstPage = false
			}
			f.Printer.Table(projects.Header(), projects.Rows())
			projects = nil // reset projects
			if next {
				var err error
				next, err = f.Prompter.Confirm("next page?", true)
				if err != nil {
					return fmt.Errorf("failed to confirm: %w", err)
				}
			}
		}
	}

	if !f.Interactive {
		f.Printer.Table(projects.Header(), projects.Rows())
	}

	return nil
}
