package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/util"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id   string
	name string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get project",
		Long:    "Get project, use --id or --name to specify the project",
		PreRunE: util.DefaultIDNameByContext(f.Config.GetContext().GetProject(), &opts.id, &opts.name),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	util.AddProjectParam(cmd, &opts.id, &opts.name)

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
	if _, err := f.ParamFiller.ProjectByName(&opts.id, &opts.name); err != nil {
		return err
	}

	return runGetNonInteractive(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}

	ownerName := f.Config.GetUsername()

	project, err := f.ApiClient.GetProject(context.Background(), opts.id, ownerName, opts.name)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	f.Printer.Table(project.Header(), project.Rows())

	return nil
}

func paramCheck(opts *Options) error {
	if opts.id == "" && opts.name == "" {
		return fmt.Errorf("please specify --id or --name")
	}

	return nil
}
