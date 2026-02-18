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
	projectID     string
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

	cmd.Flags().StringVar(&opts.projectID, "project-id", zctx.GetProject().GetID(), "Project ID")
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
	zctx := f.Config.GetContext()

	if opts.projectID == "" {
		opts.projectID = zctx.GetProject().GetID()
	}

	if opts.deploymentID == "" {
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
	// Resolve serviceID from serviceName first
	if opts.serviceID == "" && opts.serviceName != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.serviceName)
		if err != nil {
			return fmt.Errorf("failed to get service: %w", err)
		}
		opts.serviceID = service.ID
	}

	// When serviceID is available, always resolve projectID and environmentID from the service
	// instead of relying on context (which may point to a different project).
	if opts.serviceID != "" {
		service, err := f.ApiClient.GetService(context.Background(), opts.serviceID, "", "", "")
		if err != nil {
			return fmt.Errorf("failed to get service: %w", err)
		}
		if service.Project != nil {
			opts.projectID = service.Project.ID
		}
		envID, resolveErr := util.ResolveEnvironmentID(f.ApiClient, opts.projectID)
		if resolveErr != nil {
			return resolveErr
		}
		opts.environmentID = envID
	}

	// Fallback: resolve environmentID from context project if still empty
	if opts.deploymentID == "" && opts.environmentID == "" {
		projectID := opts.projectID
		if projectID == "" {
			projectID = f.Config.GetContext().GetProject().GetID()
		}
		envID, resolveErr := util.ResolveEnvironmentID(f.ApiClient, projectID)
		if resolveErr != nil {
			return resolveErr
		}
		opts.environmentID = envID
	}

	if err = paramCheck(opts); err != nil {
		return err
	}

	if opts.watch {
		return watchLogs(f, opts)
	}
	return queryLogs(f, opts)
}

func queryLogs(f *cmdutil.Factory, opts *Options) error {
	var logs model.Logs
	var err error

	switch opts.logType {
	case logTypeRuntime:
		logs, err = f.ApiClient.GetRuntimeLogs(context.Background(), opts.serviceID, opts.environmentID, opts.deploymentID)
		if err != nil {
			return fmt.Errorf("failed to get runtime logs: %w", err)
		}
	case logTypeBuild:
		deploymentID := opts.deploymentID
		if deploymentID == "" {
			deployment, exist, e := f.ApiClient.GetLatestDeployment(context.Background(), opts.serviceID, opts.environmentID)
			if e != nil {
				return fmt.Errorf("failed to get latest deployment: %w", e)
			}
			if !exist {
				return fmt.Errorf("no deployment found for service %s and environment %s", opts.serviceID, opts.environmentID)
			}
			deploymentID = deployment.ID
			f.Log.Infof("Deployment ID: %s", deploymentID)
		}
		logs, err = f.ApiClient.GetBuildLogs(context.Background(), deploymentID)
		if err != nil {
			return fmt.Errorf("failed to get build logs: %w", err)
		}
	default:
		return fmt.Errorf("unknown log type: %s", opts.logType)
	}

	f.Printer.Table(logs.Header(), logs.Rows())
	return nil
}

func watchLogs(f *cmdutil.Factory, opts *Options) error {
	var logChan <-chan model.Log
	var err error

	switch opts.logType {
	case logTypeRuntime:
		if opts.serviceID == "" || opts.environmentID == "" {
			return errors.New("service-id and env-id are required for watching runtime logs")
		}
		logChan, err = f.ApiClient.WatchRuntimeLogs(context.Background(), opts.projectID, opts.serviceID, opts.environmentID, opts.deploymentID)
	case logTypeBuild:
		deploymentID := opts.deploymentID
		if deploymentID == "" {
			deployment, exist, e := f.ApiClient.GetLatestDeployment(context.Background(), opts.serviceID, opts.environmentID)
			if e != nil {
				return fmt.Errorf("failed to get latest deployment: %w", e)
			}
			if !exist {
				return fmt.Errorf("no deployment found for service %s and environment %s", opts.serviceID, opts.environmentID)
			}
			deploymentID = deployment.ID
			f.Log.Infof("Deployment ID: %s", deploymentID)
		}
		logChan, err = f.ApiClient.WatchBuildLogs(context.Background(), opts.projectID, deploymentID)
	default:
		return fmt.Errorf("unknown log type: %s", opts.logType)
	}

	if err != nil {
		return fmt.Errorf("failed to watch logs: %w", err)
	}

	for log := range logChan {
		f.Printer.Table(log.Header(), log.Rows())
	}

	return nil
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
