package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	ID string

	OwnerName   string
	ProjectName string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		OwnerName: f.Config.GetUsername(),
	}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get project",
		Long:  "Get project, use --id or --name to specify the project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.ID, "id", "", "Project ID")

	cmd.Flags().StringVarP(&opts.ProjectName, "name", "n", "", "Project name")

	return cmd
}

// Note: don't import other packages directly in this function, or it will be hard to mock and test
// If you want to add new dependencies, please add them in the Options struct

func runGet(f *cmdutil.Factory, opts *Options) error {
	// if param check passed, run non-interactive mode first
	if err := paramCheck(opts); err == nil {
		return runGetNonInteractive(f, opts)
	}

	if f.Interactive {
		return runGetInteractive(f, opts)
	} else {
		return runGetNonInteractive(f, opts)
	}
}

func runGetInteractive(f *cmdutil.Factory, opts *Options) error {
	projects, err := getAllProjects(f)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		f.Log.Info("There are no projects in your account")
		return nil
	}

	// if there is only one project, no need to ask
	if len(projects) == 1 {
		logProject(f, projects[0])
		return nil
	}

	projectNames := make([]string, len(projects))
	for i, project := range projects {
		projectNames[i] = project.Name
	}

	index, err := f.Prompter.Select("Select a project", projects[0].Name, projectNames)
	if err != nil {
		return err
	}

	logProject(f, projects[index])

	return nil
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}

	project, err := f.ApiClient.GetProject(context.Background(), opts.ID, opts.OwnerName, opts.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	logProject(f, project)

	return nil
}

func paramCheck(opts *Options) error {
	if opts.ID == "" && opts.ProjectName == "" {
		return fmt.Errorf("please specify --id or --name")
	}

	return nil
}

func getAllProjects(f *cmdutil.Factory) ([]*model.Project, error) {
	skip := 0
	next := true

	var projects []*model.Project

	for next {
		projectCon, err := f.ApiClient.ListProjects(context.Background(), skip, 5)
		if err != nil {
			return nil, err
		}
		for _, project := range projectCon.Edges {
			projects = append(projects, project.Node)
		}

		skip += 5
		next = projectCon.PageInfo.HasNextPage
	}

	return projects, nil
}

// todo: pretty print
func logProject(f *cmdutil.Factory, project *model.Project) {
	f.Log.Info(project)
}
