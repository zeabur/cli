package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	id   string
	name string
}

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete project",
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(f, opts)
		},
	}

	zctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.name, "name", zctx.GetProject().GetName(), "Project Name")
	cmd.Flags().StringVar(&opts.id, "id", zctx.GetProject().GetID(), "Project ID")

	return cmd
}

func runDelete(f *cmdutil.Factory, opts *Options) error {
	if err := checkParams(opts); err == nil {
		return runDeleteNonInteractive(f, opts)
	}

	if f.Interactive {
		return runDeleteInteractive(f, opts)
	} else {
		return runDeleteNonInteractive(f, opts)
	}
}

func runDeleteNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := checkParams(opts); err != nil {
		return err
	}

	if opts.id == "" && opts.name != "" {
		project, err := util.GetProjectByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = project.ID
	}

	if err := deleteProject(f, opts.id); err != nil {
		return err
	}

	return nil
}

func runDeleteInteractive(f *cmdutil.Factory, opts *Options) error {
	_, project, err := f.Selector.SelectProject()
	if err != nil {
		return err
	}

	confirm, err := f.Prompter.Confirm(fmt.Sprintf("Are you sure you want to delete project %q?", project.Name), true)
	if err != nil {
		return err
	}

	if !confirm {
		f.Log.Info("Delete project canceled")
		return nil
	}

	if err := deleteProject(f, project.ID); err != nil {
		return err
	}

	return nil
}

func deleteProject(f *cmdutil.Factory, projectID string) error {
	err := f.ApiClient.DeleteProject(context.Background(), projectID)
	if err != nil {
		f.Log.Error(err)
		return err
	}

	f.Log.Info("Delete project successfully")

	return nil
}

func checkParams(opts *Options) error {
	if opts.name == "" && opts.id == "" {
		return fmt.Errorf("please specify project by --name or --id")
	}

	return nil
}
