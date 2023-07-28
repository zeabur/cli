package list

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	//todo: support service name
	serviceID     string
	environmentID string

	projectID string
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		projectID: f.Config.GetContext().GetProject().GetID(),
	}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List deployments",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	zctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.serviceID, "service-id", zctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.environmentID, "environment", zctx.GetEnvironment().GetID(), "Environment ID")

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
	if _, err := f.ParamFiller.ServiceWithEnvironment(&opts.projectID, &opts.serviceID, &opts.environmentID); err != nil {
		return err
	}

	return runListNonInteractive(f, opts)
}

func runListNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}

	deployments, err := f.ApiClient.ListAllDeployments(context.Background(), opts.serviceID, opts.environmentID)
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	if len(deployments) == 0 {
		f.Log.Info("No deployments found")
		return nil
	}

	f.Printer.Table(deployments.Header(), deployments.Rows())

	return nil
}

func paramCheck(opts *Options) error {
	if opts.serviceID == "" {
		return errors.New("service-id is required")
	}

	if opts.environmentID == "" {
		return errors.New("environment is required")
	}

	return nil
}
