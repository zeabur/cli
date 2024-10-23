package deploy

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/selector"
	"github.com/zeabur/cli/pkg/zcontext"
)

type Options struct {
	name       string
	domainName string

	// create will create a new service instead of selecting one if true
	create bool

	// specify a service ID and environment ID to deploy on

	serviceID     string
	environmentID string
}

func NewCmdDeploy(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   "Deploy local project to Zeabur with one command",
		PreRunE: util.NeedProjectContextWhenNonInteractive(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeploy(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.serviceID, "service-id", "", "Service ID to redeploy on")
	cmd.Flags().StringVar(&opts.environmentID, "environment-id", "", "Environment ID to redeploy on")
	cmd.Flags().StringVar(&opts.name, "name", "", "Service name")
	cmd.Flags().StringVar(&opts.domainName, "domain", "", "Domain name")
	cmd.Flags().BoolVar(&opts.create, "create", false, "Create a new service")

	return cmd
}

func runDeploy(f *cmdutil.Factory, opts *Options) error {
	var environment *model.Environment
	var service *model.Service
	var err error

	bytes, fileName, err := util.PackZip()
	if err != nil {
		return fmt.Errorf("packing zip: %w", err)
	}
	opts.name = fileName

	if opts.serviceID == "" {
		service, environment, err = selectInteractively(f, opts)
		if err != nil {
			return err
		}
	} else {
		service, err = f.ApiClient.GetService(context.Background(), opts.serviceID, "", "", "")
		if err != nil {
			return err
		}

		environment, err = f.ApiClient.GetEnvironment(context.Background(), opts.environmentID)
		if err != nil {
			return err
		}
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Uploading codes to Zeabur ..."),
	)
	s.Start()

	_, err = f.ApiClient.UploadZipToService(context.Background(), service.Project.ID, service.ID, environment.ID, bytes)
	if err != nil {
		return err
	}
	s.Stop()

	domainName := opts.domainName

	if domainName == "" {
		fmt.Println("Service deployed successfully, you can access it via:")
		fmt.Println("https://dash.zeabur.com/projects/" + service.Project.ID + "/services/" + service.ID + "?envID=" + environment.ID)
		return nil
	}

	s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Creating domain ..."),
	)
	s.Start()

	domain, err := f.ApiClient.AddDomain(context.Background(), service.ID, environment.ID, false, domainName)
	if err != nil {
		return err
	}

	fmt.Println("Domain created: ", "https://"+*domain)

	s.Stop()

	return nil
}

func selectInteractively(f *cmdutil.Factory, opts *Options) (*model.Service, *model.Environment, error) {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching projects ..."),
	)
	s.Start()
	projects, err := f.ApiClient.ListAllProjects(context.Background())
	if err != nil {
		return nil, nil, err
	}
	s.Stop()

	if len(projects) == 0 {
		confirm, err := f.Prompter.Confirm("No projects found. Would you like to create one now?", true)
		if err != nil {
			return nil, nil, err
		}
		if confirm {
			project, err := f.ApiClient.CreateProject(context.Background(), "default", nil)
			if err != nil {
				f.Log.Error("Failed to create project: ", err)
				return nil, nil, err
			}
			f.Log.Infof("Project %s created. Run this command again to deploy on it.", project.Name)
			f.Config.GetContext().SetProject(zcontext.NewBasicInfo(project.ID, project.Name))

			return nil, nil, nil
		}
	}

	f.Log.Info("Select one project to deploy your service.")

	_, project, err := f.Selector.SelectProject()
	if err != nil {
		return nil, nil, err
	}

	f.Config.GetContext().SetProject(zcontext.NewBasicInfo(project.ID, project.Name))

	_, environment, err := f.Selector.SelectEnvironment(project.ID)
	if err != nil {
		return nil, nil, err
	}

	f.Log.Info("Select one service to deploy or create a new one.")

	var service *model.Service
	if !opts.create {
		_, service, err = f.Selector.SelectService(selector.SelectServiceOptions{
			ProjectID: project.ID,
			Auto:      false,
			CreateNew: true,
			FilterFunc: func(service *model.Service) bool {
				return service.Template == "GIT"
			},
		})
		if err != nil {
			return nil, nil, err
		}
	}

	if service == nil {
		f.Log.Info("No service found, create a new one.")

		s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Creating new service ..."),
		)
		s.Start()

		name := opts.name

		service, err = f.ApiClient.CreateEmptyService(context.Background(), project.ID, name)
		if err != nil {
			return nil, nil, err
		}

		s.Stop()
	}

	return service, environment, nil
}
