package get

import (
	"context"
	"fmt"
	"github.com/zeabur/cli/pkg/model"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string

	ownerName   string
	projectName string
	name        string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		ownerName: f.Config.GetUsername(),
	}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	ctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.id, "id", ctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.name, "name", ctx.GetService().GetName(), "Service name")
	cmd.Flags().StringVar(&opts.projectName, "project-name", ctx.GetProject().GetName(), "Service project name")

	return cmd
}

func runGet(f *cmdutil.Factory, opts *Options) error {
	err := paramCheck(opts)

	if f.Interactive {
		// if param check passed, run non-interactive mode
		if err == nil {
			return runGetNonInteractive(f, opts)
		}

		return runGetInteractive(f, opts)
	} else {
		if err != nil {
			return err
		}
		return runGetNonInteractive(f, opts)
	}
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()
	service, err := f.ApiClient.GetService(ctx, opts.id, opts.ownerName, opts.projectName, opts.name)
	if err != nil {
		return fmt.Errorf("get service failed: %w", err)
	}

	f.Log.Infof("Selected service: %s(%s)", service.Name, service.ID)

	return nil
}

func runGetInteractive(f *cmdutil.Factory, opts *Options) error {
	projectCtx := f.Config.GetContext().GetProject()
	if projectCtx.Empty() {
		return fmt.Errorf("please use `zc project set` to set the project context first")
	}

	// if id or (projectName and name) is specified, we have used non-interactive mode
	// therefore, now the id and name must be empty

	_, service, err := f.Selector.SelectService(projectCtx.GetID())
	if err != nil {
		return fmt.Errorf("failed to select service: %w", err)
	}

	logService(f, service)

	return nil
}

func paramCheck(opts *Options) error {
	if opts.id != "" || (opts.projectName != "" && opts.name != "") {
		return nil
	}

	return fmt.Errorf("please specify --id or (--project-name and --name)")
}

func logService(f *cmdutil.Factory, service *model.Service) {
	services := model.Services{service}
	f.Printer.Table(services.Header(), services.Rows())
}
