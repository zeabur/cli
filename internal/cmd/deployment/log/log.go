package log

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	serviceID     string
	serviceName   string
	environmentID string

	deploymentID string

	logType string
	watch   bool
}

const (
	logTypeRuntime = "runtime"
	logTypeBuild   = "build"
)

func NewCmdLog(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "log",
		Short:   "Get deployment logs, if deployment-id is not specified, use serviceID/serviceName and environmentID to get the deployment",
		PreRunE: util.NeedProjectContextWhenNonInteractive(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLog(f, opts)
		},
	}

	zctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.deploymentID, "deployment-id", "", "Deployment ID")
	cmd.Flags().StringVar(&opts.serviceID, "service-id", zctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.serviceName, "service-name", zctx.GetService().GetName(), "Service Name")
	cmd.Flags().StringVar(&opts.environmentID, "env-id", zctx.GetEnvironment().GetID(), "Environment ID")
	cmd.Flags().StringVarP(&opts.logType, "type", "t", logTypeRuntime, "Log type, runtime or build")
	cmd.Flags().BoolVarP(&opts.watch, "watch", "w", false, "Watch logs")

	return cmd
}

func runLog(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runLogInteractive(f, opts)
	} else {
		return runLogNonInteractive(f, opts)
	}
}

func runLogInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.deploymentID == "" {
		zctx := f.Config.GetContext()
		_, err := f.ParamFiller.ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions{
			ProjectCtx:    zctx,
			ServiceID:     &opts.serviceID,
			ServiceName:   &opts.serviceName,
			EnvironmentID: &opts.environmentID,
			CreateNew:     false,
		})
		if err != nil {
			return err
		}
	}

	return runLogNonInteractive(f, opts)
}

func runLogNonInteractive(f *cmdutil.Factory, opts *Options) (err error) {
	if err = paramCheck(opts); err != nil {
		return err
	}

	if opts.watch {
		if opts.deploymentID != "" {
			var logChan <-chan model.Log
			var subscriptionErr error

			switch opts.logType {
			case logTypeRuntime:
				logChan, subscriptionErr = f.ApiClient.WatchRuntimeLogs(context.Background(), opts.deploymentID)
			case logTypeBuild:
				logChan, subscriptionErr = f.ApiClient.WatchBuildLogs(context.Background(), opts.deploymentID)
			default:
				logChan, subscriptionErr = f.ApiClient.WatchRuntimeLogs(context.Background(), opts.deploymentID)
			}
			if subscriptionErr != nil {
				return fmt.Errorf("failed to watch logs: %w", err)
			}

			for log := range logChan {
				f.Printer.Table(log.Header(), log.Rows())
			}
		}

		return nil
	} else {
		var logs model.Logs

		// If deployment id is provided, get deployment by deployment id
		if opts.deploymentID != "" {
			logs, err = logDeploymentByID(f, opts.deploymentID, opts.logType)
		} else {
			// or, get deployment by service id and environment id

			// If service id is not provided, get service id by service name
			if opts.serviceID == "" {
				var service *model.Service
				if service, err = util.GetServiceByName(f.Config, f.ApiClient, opts.serviceName); err != nil {
					return fmt.Errorf("failed to get service: %w", err)
				} else {
					opts.serviceID = service.ID
				}
			}
			logs, err = logDeploymentByServiceAndEnvironment(f, opts.serviceID, opts.environmentID, opts.logType)
		}

		if err != nil {
			return err
		}

		f.Printer.Table(logs.Header(), logs.Rows())

		return nil
	}
}

func logDeploymentByID(f *cmdutil.Factory, deploymentID, logType string) (model.Logs, error) {
	switch logType {
	case logTypeRuntime:
		logs, err := f.ApiClient.GetRuntimeLogs(context.Background(), deploymentID, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to get runtime logs: %w", err)
		}
		return logs, nil
	case logTypeBuild:
		logs, err := f.ApiClient.GetBuildLogs(context.Background(), deploymentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get build logs: %w", err)
		}
		return logs, nil
	default:
		return nil, fmt.Errorf("unknown log type: %s", logType)
	}
}

func logDeploymentByServiceAndEnvironment(f *cmdutil.Factory, serviceID, environmentID, logType string) (model.Logs, error) {
	switch logType {
	case logTypeRuntime:
		logs, err := f.ApiClient.GetRuntimeLogs(context.Background(), "", serviceID, environmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get runtime logs: %w", err)
		}
		return logs, nil
	case logTypeBuild:
		deployment, exist, err := f.ApiClient.GetLatestDeployment(context.Background(), serviceID, environmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest deployment: %w", err)
		}
		if !exist {
			return nil, fmt.Errorf("no deployment found for service %s and environment %s", serviceID, environmentID)
		}
		f.Log.Infof("Deployment ID: %s", deployment.ID)
		logs, err := f.ApiClient.GetBuildLogs(context.Background(), deployment.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get build logs: %w", err)
		}
		return logs, nil
	default:
		return nil, fmt.Errorf("unknown log type: %s", logType)
	}
}

func paramCheck(opts *Options) error {
	if opts.logType != logTypeRuntime && opts.logType != logTypeBuild {
		return errors.New("log type must be runtime or build")
	}
	if opts.deploymentID != "" {
		return nil
	}

	if opts.serviceID == "" && opts.serviceName == "" {
		return errors.New("when deployment-id is not specified, service-id or service-name is required")
	}

	if opts.environmentID == "" {
		return errors.New("when deployment-id is not specified, env-id is required")
	}

	return nil
}
