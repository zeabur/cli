package deploy

import (
	"context"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/zcontext"
)

type Options struct {
	name string
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

	f.Log.Info("You have selected project %s", project.Name)

	f.Config.GetContext().SetProject(zcontext.NewBasicInfo(project.ID, project.Name))

	return nil
}
