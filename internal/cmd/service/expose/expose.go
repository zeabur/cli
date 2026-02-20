package expose

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id   string // use id or name to specify service
	name string

	environmentID string // environmentID is required
}

func NewCmdExpose(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "expose",
		Short: "Expose a service temporarily",
		Long: `Expose a service temporarily, default 3600 seconds.
example:
      zeabur service expose # cli will prompt to select a service
	  zeabur service expose --id xxxxx --env-id xxxx # use id and env-id to expose service
      zeabur service expose --name xxxxx --env-id xxxx # use name, env-id to expose service
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExpose(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)

	return cmd
}

func runExpose(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runExposeInteractive(f, opts)
	} else {
		return runExposeNonInteractive(f, opts)
	}
}

func runExposeNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" && opts.name != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = service.ID
	}

	if opts.id == "" {
		return fmt.Errorf("please specify --id or --name")
	}

	// Resolve environment and project from the service
	service, err := f.ApiClient.GetService(context.Background(), opts.id, "", "", "")
	if err != nil {
		return fmt.Errorf("get service failed: %w", err)
	}

	if opts.environmentID == "" {
		if service.Project == nil || service.Project.ID == "" {
			return fmt.Errorf("service has no associated project")
		}
		envID, err := util.ResolveEnvironmentID(f.ApiClient, service.Project.ID)
		if err != nil {
			return err
		}
		opts.environmentID = envID
	}

	projectID := ""
	if service.Project != nil {
		projectID = service.Project.ID
	}

	ctx := context.Background()
	tempTCPPort, err := f.ApiClient.ExposeService(ctx, opts.id, opts.environmentID, projectID, opts.name)
	if err != nil {
		return fmt.Errorf("failed to expose service: %w", err)
	}

	f.Log.Infof("Service is exposed on port %d, and will be closed after %d seconds\n",
		tempTCPPort.NodePort, tempTCPPort.RemainSeconds)

	return nil
}

func runExposeInteractive(f *cmdutil.Factory, opts *Options) error {
	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions{
		ProjectCtx:    f.Config.GetContext(),
		ServiceID:     &opts.id,
		ServiceName:   &opts.name,
		EnvironmentID: &opts.environmentID,
		CreateNew:     false,
	}); err != nil {
		return err
	}

	return runExposeNonInteractive(f, opts)
}
