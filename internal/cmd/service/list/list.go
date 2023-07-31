package list

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	projectID     string
	environmentID string
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		projectID: f.Config.GetContext().GetProject().GetID(),
	}
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List services",
		Long:    `List services, if env-id is provided, list services in the environment in detail`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	ctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.environmentID, "env-id", ctx.GetEnvironment().GetID(), "Environment ID")

	return cmd
}

func runList(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runListInteractive(f, opts)
	} else {
		return runListNonInteractive(f, opts)
	}
}

func runListInteractive(f *cmdutil.Factory, opts *Options) error {
	// fetch project id from context
	opts.projectID = f.Config.GetContext().GetProject().GetID()

	// if project id is not set, prompt to select one
	if _, err := f.ParamFiller.Project(&opts.projectID); err != nil {
		return err
	}

	return runListNonInteractive(f, opts)
}

func runListNonInteractive(f *cmdutil.Factory, opts *Options) error {
	err := paramCheck(opts)
	if err != nil {
		return err
	}

	if opts.environmentID == "" {
		return listServicesBrief(f, opts.projectID)
	} else {
		return listServicesDetailByEnvironment(f, opts.projectID, opts.environmentID)
	}
}

func listServicesBrief(f *cmdutil.Factory, projectID string) error {
	services, err := f.ApiClient.ListAllServices(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("list services failed: %w", err)
	}

	if len(services) == 0 {
		f.Log.Infof("No services found")
		return nil
	}

	f.Printer.Table(services.Header(), services.Rows())

	return nil
}

func listServicesDetailByEnvironment(f *cmdutil.Factory, projectID, environmentID string) error {
	services, err := f.ApiClient.ListAllServicesDetailByEnvironment(context.Background(), projectID, environmentID)
	if err != nil {
		return fmt.Errorf("list services failed: %w", err)
	}

	if len(services) == 0 {
		f.Log.Infof("No services found")
		return nil
	}

	f.Printer.Table(services.Header(), services.Rows())

	return nil
}

func paramCheck(opts *Options) error {
	if opts.projectID == "" {
		return fmt.Errorf("project-id is required")
	}

	return nil
}
