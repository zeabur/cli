package list

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	//todo: support service name
	serviceID     string
	serviceName   string
	environmentID string
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List deployments",
		Aliases: []string{"ls"},
		PreRunE: util.NeedProjectContext(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	zctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.serviceID, "service-id", zctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.serviceName, "service-name", zctx.GetService().GetName(), "Service Name")
	cmd.Flags().StringVar(&opts.environmentID, "env-id", zctx.GetEnvironment().GetID(), "Environment ID")

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
	zctx := f.Config.GetContext()
	_, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.serviceID, &opts.serviceName, &opts.environmentID)
	if err != nil {
		return err
	}

	return runListNonInteractive(f, opts)
}

func runListNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}

	// If service id is not provided, get service id by service name
	if opts.serviceID == "" {
		if service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.serviceName); err != nil {
			return fmt.Errorf("failed to get service: %w", err)
		} else {
			opts.serviceID = service.ID
		}
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
	if opts.serviceID == "" && opts.serviceName == "" {
		return errors.New("service-id or service-name is required")
	}

	if opts.environmentID == "" {
		return errors.New("environment is required")
	}

	return nil
}
