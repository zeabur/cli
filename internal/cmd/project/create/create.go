package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	ProjectName string
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(f, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.ProjectName, "name", "n", "", "Project name")

	return cmd
}

func runCreate(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err == nil {
		return runCreateNonInteractive(f, opts)
	}

	if f.Interactive {
		return runCreateInteractive(f, opts)
	}

	return runCreateNonInteractive(f, opts)

}

func runCreateInteractive(f *cmdutil.Factory, opts *Options) error {
	projectName, err := f.Prompter.Input("Please input project name:", "")
	if err != nil {
		return err
	}

	if err := createProject(f, projectName); err != nil {
		return err
	}

	return nil
}

func runCreateNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}

	err := createProject(f, opts.ProjectName)
	if err != nil {
		return err
	}

	return nil
}

func createProject(f *cmdutil.Factory, projectName string) error {
	project, err := f.ApiClient.CreateProject(context.Background(), projectName)
	if err != nil {
		f.Log.Error(err)
		return err
	}

	f.Log.Infof("Project %s created", project.Name)

	return nil
}

func paramCheck(opts *Options) error {
	if opts.ProjectName == "" {
		return fmt.Errorf("please specify project name with --name")
	}

	return nil
}
