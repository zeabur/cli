package get

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	ID string

	OwnerName   string
	ProjectName string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		OwnerName: f.Config.GetUsername(),
	}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get project",
		Long:  "Get project, use --id or --name to specify the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	projectCtx := f.Config.GetContext().GetProject()

	cmd.Flags().StringVar(&opts.ID, "id", projectCtx.GetID(), "Project ID")
	cmd.Flags().StringVarP(&opts.ProjectName, "name", "n", projectCtx.GetName(), "Project name")

	return cmd
}

// Note: don't import other packages directly in this function, or it will be hard to mock and test
// If you want to add new dependencies, please add them in the Options struct

func runGet(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	} else {
		return runGetNonInteractive(f, opts)
	}
}

func runGetInteractive(f *cmdutil.Factory, opts *Options) error {
	if _, err := f.ParamFiller.Project(&opts.ID); err != nil {
		return err
	}

	return runGetNonInteractive(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}

	project, err := f.ApiClient.GetProject(context.Background(), opts.ID, opts.OwnerName, opts.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	logProject(f, project)

	return nil
}

func paramCheck(opts *Options) error {
	if opts.ID == "" && opts.ProjectName == "" {
		return fmt.Errorf("please specify --id or --name")
	}

	return nil
}

func logProject(f *cmdutil.Factory, p *model.Project) {
	projects := model.Projects{p}

	f.Printer.Table(projects.Header(), projects.Rows())
}
