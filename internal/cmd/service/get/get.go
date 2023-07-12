package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
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

	logService(f, service)

	return nil
}

func runGetInteractive(f *cmdutil.Factory, opts *Options) error {
	projectCtx := f.Config.GetContext().GetProject()
	if projectCtx.Empty() {
		return fmt.Errorf("please use `zc project set` to set the project context first")
	}

	projectID, projectName := projectCtx.GetID(), projectCtx.GetName()

	// if id or (projectName and name) is specified, we have used non-interactive mode
	// therefore, now the id and name must be empty

	services, err := f.ApiClient.ListAllServices(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("list services failed: %w", err)
	}

	if len(services) == 0 {
		return fmt.Errorf("no service found in project %s", projectName)
	}

	var index int

	if len(services) == 1 {
		index = 0
		f.Log.Info("Only one service found, select it by default")
	} else {
		serviceNames := make([]string, 0, len(services))
		for _, service := range services {
			serviceNames = append(serviceNames, service.Name)
		}

		index, err = f.Prompter.Select("Select a service", serviceNames[0], serviceNames)
		if err != nil {
			return fmt.Errorf("select service failed: %w", err)
		}
	}

	logService(f, services[index])

	return nil
}

func paramCheck(opts *Options) error {
	if opts.id != "" || (opts.projectName != "" && opts.name != "") {
		return nil
	}

	return fmt.Errorf("please specify --id or (--project-name and --name)")
}

// todo: pretty print service
func logService(f *cmdutil.Factory, service *model.Service) {
	f.Log.Info(service)
}
