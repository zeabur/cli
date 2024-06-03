package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	id   string
	name string
	yes  bool
}

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete project",
		Aliases: []string{"del"},
		PreRunE: util.DefaultIDNameByContext(f.Config.GetContext().GetProject(), &opts.id, &opts.name),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(f, opts)
		},
	}

	util.AddProjectParam(cmd, &opts.id, &opts.name)
	cmd.Flags().BoolVarP(&opts.yes, "yes", "y", false, "Delete project without confirmation")

	return cmd
}

func runDelete(f *cmdutil.Factory, opts *Options) error {
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

	if !opts.yes {
		f.Log.Info("Please use --yes to confirm deletion without interactive prompt")
		return nil
	}

	project, err := f.ApiClient.GetProject(context.Background(), opts.id, "", "")
	if err != nil {
		f.Log.Error(err)
		return err
	}

	if err := deleteProject(f, project); err != nil {
		return err
	}

	return nil
}

func runDeleteInteractive(f *cmdutil.Factory, opts *Options) error {
	_, project, err := f.Selector.SelectProject()
	if err != nil {
		return err
	}

	if !opts.yes {
		confirm, err := f.Prompter.Confirm(fmt.Sprintf("Are you sure you want to delete project %q (%s) ?", project.Name, project.ID), false)
		if err != nil {
			return err
		}

		if !confirm {
			f.Log.Info("Delete project canceled")
			return nil
		}
	}

	if err := deleteProject(f, project); err != nil {
		return err
	}

	return nil
}

func deleteProject(f *cmdutil.Factory, project *model.Project) error {
	err := f.ApiClient.DeleteProject(context.Background(), project.ID)
	if err != nil {
		f.Log.Error(err)
		return err
	}

	f.Log.Infof("Delete project %s (%s) successfully", project.Name, project.ID)
	return nil
}

func checkParams(opts *Options) error {
	if opts.name == "" && opts.id == "" {
		return fmt.Errorf("please specify project by --name or --id")
	}

	return nil
}
