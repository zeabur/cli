package expose

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string // use id to specify service

	projectID string // use projectID and serviceName to specify service
	name      string

	environmentID string // environmentID is required
}

func NewCmdExpose(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "expose",
		Short: "Expose a service temporarily",
		Long: `Expose a service temporarily, default 3600 seconds.
example:
      zeabur service expose # cli will try to get service from context or prompt to select one
	  zeabur service expose --id xxxxx --environment-id xxxx # use id and environment-id to expose service
      zeabur service expose --name xxxxx --project-id xxxx --environment-id xxxx # use name, project-id and environment-id to expose service
      zeabur service expose --name xxxxx --environment-id xxxx # if project context is set, use name, environment-id to expose service
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExpose(f, opts)
		},
	}

	ctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.id, "id", ctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.name, "name", ctx.GetService().GetName(), "Service name")
	cmd.Flags().StringVar(&opts.projectID, "project-id", ctx.GetProject().GetID(), "Service project name")
	cmd.Flags().StringVar(&opts.environmentID, "environment-id", ctx.GetEnvironment().GetID(), "Service environment ID")

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
	err := paramCheck(opts)
	if err != nil {
		return err
	}

	ctx := context.Background()

	tempTCPPort, err := f.ApiClient.ExposeService(ctx, opts.id, opts.environmentID, opts.projectID, opts.name)
	if err != nil {
		return fmt.Errorf("failed to expose service: %w", err)
	}

	f.Log.Infof("Service is exposed on port %d, and will be closed after %d seconds\n",
		tempTCPPort.NodePort, tempTCPPort.RemainSeconds)

	return nil
}

func runExposeInteractive(f *cmdutil.Factory, opts *Options) error {
	// if serviceID is not set, we need to find services by projectID, so projectID is required
	if opts.id == "" {
		if opts.projectID == "" {
			projects, err := f.ApiClient.ListAllProjects(context.Background())
			if err != nil {
				return fmt.Errorf("failed to list projects: %w", err)
			}
			if len(projects) == 0 {
				return fmt.Errorf("no project found")
			}
			if len(projects) == 1 {
				opts.projectID = projects[0].ID
				f.Log.Infof("Only one project found, select <%s> automatically\n", projects[0].Name)
			} else {
				projectNames := make([]string, len(projects))
				for i, project := range projects {
					projectNames[i] = project.Name
				}
				index, err := f.Prompter.Select("Select project", projects[0].Name, projectNames)
				if err != nil {
					return fmt.Errorf("failed to select project: %w", err)
				}
				opts.projectID = projects[index].ID
			}
		}
	}
	// if environmentID is not set, list environments and select one
	if opts.environmentID == "" {
		envs, err := f.ApiClient.ListEnvironments(context.Background(), opts.projectID)
		if err != nil {
			return fmt.Errorf("failed to list environments: %w", err)
		}
		if len(envs) == 0 {
			return fmt.Errorf("no environment found")
		}
		if len(envs) == 1 {
			opts.environmentID = envs[0].ID
			f.Log.Infof("Only one environment found, select <%s> automatically\n", envs[0].Name)
		} else {
			envNames := make([]string, len(envs))
			for i, env := range envs {
				envNames[i] = env.Name
			}
			index, err := f.Prompter.Select("Select environment", envs[0].Name, envNames)
			if err != nil {
				return fmt.Errorf("failed to select environment: %w", err)
			}
			opts.environmentID = envs[index].ID
		}
	}

	// either serviceID or (projectID and serviceName) is required
	if opts.id == "" && opts.name == "" {
		services, err := f.ApiClient.ListAllServices(context.Background(), opts.projectID)
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}
		if len(services) == 0 {
			return fmt.Errorf("no service found")
		}
		if len(services) == 1 {
			opts.id = services[0].ID
			f.Log.Infof("Only one service found, select <%s> automatically\n", services[0].Name)
		} else {
			serviceNames := make([]string, len(services))
			for i, service := range services {
				serviceNames[i] = service.Name
			}
			index, err := f.Prompter.Select("Select service", services[0].Name, serviceNames)
			if err != nil {
				return fmt.Errorf("failed to select service: %w", err)
			}
			opts.id = services[index].ID
		}
	}
	return runExposeNonInteractive(f, opts)
}

func paramCheck(opts *Options) error {
	if !(opts.id != "" || (opts.projectID != "" && opts.name != "")) {
		return fmt.Errorf("please specify --id or (--project-id and --name)")
	}

	if opts.environmentID == "" {
		return fmt.Errorf("please specify --environment-id")
	}

	return nil
}
