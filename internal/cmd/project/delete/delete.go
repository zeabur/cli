package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	ProjectID string
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

	cmd.Flags().StringVar(&opts.ProjectID, "id", "", "Project ID")

	return cmd
}

func runDelete(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err == nil {
		return runDeleteNonInteractive(f, opts)
	}

	if f.Interactive {
		return runDeleteInteractive(f, opts)
	} else {
		return runDeleteNonInteractive(f, opts)
	}
}

func runDeleteNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := deleteProject(f, opts.ProjectID); err != nil {
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

func paramCheck(opts *Options) error {
	if opts.ProjectID == "" {
		return fmt.Errorf("please specify project id by --id")
	}

	return nil
}
