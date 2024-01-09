package deploy

import (
	"context"
	"fmt"

	"math/rand"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/zcontext"
)

type Options struct {
	name       string
	domainName string
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

	cmd.Flags().StringVar(&opts.name, "name", "", "Service name")
	cmd.Flags().StringVar(&opts.domainName, "domain", "", "Domain name")

	return cmd
}

func runDeploy(f *cmdutil.Factory, opts *Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching projects ..."),
	)
	s.Start()
	projects, err := f.ApiClient.ListAllProjects(context.Background())
	if err != nil {
		return err
	}
	s.Stop()

	if len(projects) == 0 {
		confirm, err := f.Prompter.Confirm("No projects found, would you like to create one now?", true)
		if err != nil {
			return err
		}
		if confirm {
			project, err := f.ApiClient.CreateProject(context.Background(), "default", nil)
			if err != nil {
				f.Log.Error("Failed to create project: ", err)
				return err
			}
			f.Log.Infof("Project %s created", project.Name)
			f.Config.GetContext().SetProject(zcontext.NewBasicInfo(project.ID, project.Name))

			return nil
		}
	}

	f.Log.Info("Select one project to deploy your service.")

	_, project, err := f.Selector.SelectProject()
	if err != nil {
		return err
	}

	f.Config.GetContext().SetProject(zcontext.NewBasicInfo(project.ID, project.Name))

	_, environment, err := f.Selector.SelectEnvironment(project.ID)
	if err != nil {
		return err
	}

	s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Creating new service ..."),
	)
	s.Start()

	bytes, fileName, err := util.PackZip()
	if err != nil {
		return err
	}

	service, err := f.ApiClient.CreateEmptyService(context.Background(), project.ID, fileName)
	if err != nil {
		return err
	}

	s.Stop()

	s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Uploading codes to Zeabur ..."),
	)
	s.Start()

	_, err = f.ApiClient.UploadZipToService(context.Background(), project.ID, service.ID, environment.ID, bytes)
	if err != nil {
		return err
	}
	s.Stop()

	s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Creating domain ..."),
	)
	s.Start()

	// if opts.domain is empty, create a random domain
	domainName := opts.domainName
	if domainName == "" {
		var from = []rune("abcdefghijklmnopqrstuvwxyz")
		b := make([]rune, 8)
		for i := range b {
			b[i] = from[rand.Intn(len(from))]
		}
		domainName = string(b)
	}

	domain, err := f.ApiClient.AddDomain(context.Background(), service.ID, environment.ID, opts.domainName == "", domainName)
	if err != nil {
		return err
	}

	fmt.Println("Domain created: ", "https://"+*domain)

	s.Stop()

	return nil
}
