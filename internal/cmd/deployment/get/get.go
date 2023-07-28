package get

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	//todo: support service name
	serviceID     string
	environmentID string

	projectID string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get deployments",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	zctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.serviceID, "service-id", zctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.environmentID, "environment", zctx.GetEnvironment().GetID(), "Environment ID")

	return cmd
}

func runGet(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	} else {
		return runGetNonInteractive(f, opts)
	}
}

func runGetInteractive(f *cmdutil.Factory, opts *Options) error {
	// if missing service-id or environment-id, we need to use project id to get them
	if err := paramCheck(opts); err != nil {
		if opts.projectID == "" {
			basicInfo, _, err := f.Selector.SelectProject()
			if err != nil {
				return err
			}
			opts.projectID = basicInfo.GetID()
		}
	}

	if opts.serviceID == "" {
		basicInfo, _, err := f.Selector.SelectService(opts.projectID)
		if err != nil {
			return err
		}
		opts.serviceID = basicInfo.GetID()
	}

	if opts.environmentID == "" {
		basicInfo, _, err := f.Selector.SelectEnvironment(opts.projectID)
		if err != nil {
			return err
		}
		opts.environmentID = basicInfo.GetID()
	}
	return runGetNonInteractive(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}

	deployment, exist, err := f.ApiClient.GetLatestDeployment(context.Background(), opts.serviceID, opts.environmentID)
	if err != nil {
		return fmt.Errorf("failed to get deployments: %w", err)
	}

	if !exist {
		f.Log.Info("No deployments found")
		return nil
	}

	deployments := model.Deployments{deployment}

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
