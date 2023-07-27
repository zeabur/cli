package get

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeabur/cli/pkg/model"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id   string
	name string

	projectID   string
	projectName string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		projectID:   f.Config.GetContext().GetProject().GetID(),
		projectName: f.Config.GetContext().GetProject().GetName(),
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

func runGetInteractive(f *cmdutil.Factory, opts *Options) error {
	// if param missing, we should use project id to select service,
	// so, the project context must be set
	if err := paramCheck(opts); err != nil {
		if opts.projectID == "" || opts.projectName == "" {
			project, _, err := f.Selector.SelectProject()
			if err != nil {
				return err
			}
			opts.projectID = project.GetID()
			opts.projectName = project.GetName()
		}
	}

	// if id or (projectName and name) is specified, we have used non-interactive mode
	// therefore, now the id and name must be empty

	_, service, err := f.Selector.SelectService(opts.projectID)
	if err != nil {
		return fmt.Errorf("failed to select service: %w", err)
	}

	logService(f, service)

	return nil
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) error {
	// if only name is specified, user must have set the project context
	if opts.id == "" && opts.name != "" {
		if opts.projectID == "" || opts.projectName == "" {
			return errors.New("since only name is specified, please set project context first")
		}
	}

	ctx := context.Background()
	service, err := f.ApiClient.GetService(ctx, opts.id, f.Config.GetUsername(), opts.projectName, opts.name)
	if err != nil {
		return fmt.Errorf("get service failed: %w", err)
	}

	f.Log.Infof("Selected service: %s(%s)", service.Name, service.ID)

	return nil
}

func paramCheck(opts *Options) error {
	if opts.id != "" || opts.name != "" {
		return nil
	}

	return fmt.Errorf("please specify --id or --name")
}

func logService(f *cmdutil.Factory, service *model.Service) {
	services := model.Services{service}
	f.Printer.Table(services.Header(), services.Rows())
}
