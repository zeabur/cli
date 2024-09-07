package export

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	ProjectID   string
	ProjectName string

	EnvironmentID string
}

func NewCmdExport(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export projects to Template Resource YAML",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(f, opts)
		},
	}

	util.AddProjectParam(cmd, &opts.ProjectID, &opts.ProjectName)
	cmd.Flags().StringVar(&opts.EnvironmentID, "environment", "", "Environment ID to export. If not specified, we export the default environment.")

	return cmd
}

// runExport will export the project to a template resource YAML.
func runExport(f *cmdutil.Factory, opts Options) error {
	if opts.ProjectID == "" && opts.ProjectName == "" {
		return fmt.Errorf("please specify project by --name or --id")
	}

	if opts.ProjectID == "" && opts.ProjectName != "" {
		project, err := util.GetProjectByName(f.Config, f.ApiClient, opts.ProjectName)
		if err != nil {
			return fmt.Errorf("get project %s failed: %w", opts.ProjectName, err)
		}
		opts.ProjectID = project.ID
	}

	if opts.EnvironmentID == "" {
		environments, err := f.ApiClient.ListEnvironments(context.Background(), opts.ProjectID)
		if err != nil {
			return fmt.Errorf("list environments for project<%s> failed: %w", opts.ProjectID, err)
		}

		if len(environments) == 0 {
			return fmt.Errorf("no environment found in project %s", opts.ProjectID)
		}

		opts.EnvironmentID = environments[0].ID
	}

	exportedTemplate, err := f.ApiClient.ExportProject(context.Background(), opts.ProjectID, opts.EnvironmentID)
	if err != nil {
		return fmt.Errorf("export environment<%s> of project<%s> failed: %w", opts.EnvironmentID, opts.ProjectID, err)
	}

	for _, warning := range exportedTemplate.Warnings {
		f.Log.Warn(warning)
	}

	fmt.Println(exportedTemplate.ResourceYAML)

	return nil
}
