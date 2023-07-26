package list

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"strings"
)

type Options struct {
	projectID     string
	environmentID string
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List environments",
		Long:    `List environments, if environment-id is provided, list services in the environment in detail`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	ctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.projectID, "project-id", ctx.GetProject().GetID(), "Project ID")
	cmd.Flags().StringVar(&opts.environmentID, "environment-id", ctx.GetEnvironment().GetID(), "Environment ID")

	return cmd
}

func runList(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runListInteractive(f, opts)
	} else {
		return runListNonInteractive(f, opts)
	}
}

func runListInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.projectID == "" {
		projects, err := f.ApiClient.ListAllProjects(context.Background())
		if err != nil {
			return fmt.Errorf("list projects failed: %w", err)
		}
		if len(projects) == 0 {
			return fmt.Errorf("no projects found")
		}
		if len(projects) == 1 {
			opts.projectID = projects[0].ID
			f.Log.Infof("Only one project found, select %s automatically\n", projects[0].Name)
		} else {
			projectsName := make([]string, len(projects))
			for i, project := range projects {
				projectsName[i] = project.Name
			}
			index, err := f.Prompter.Select("Select project", projectsName[0], projectsName)
			if err != nil {
				return fmt.Errorf("select project failed: %w", err)
			}
			opts.projectID = projects[index].ID
		}
	}

	return runListNonInteractive(f, opts)
}

func runListNonInteractive(f *cmdutil.Factory, opts *Options) error {
	err := paramCheck(opts)
	if err != nil {
		return err
	}

	if opts.environmentID == "" {
		return listServicesBrief(f, opts.projectID)
	} else {
		return listServicesDetailByEnvironment(f, opts.projectID, opts.environmentID)
	}
}

func listServicesBrief(f *cmdutil.Factory, projectID string) error {
	services, err := f.ApiClient.ListAllServices(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("list services failed: %w", err)
	}

	if len(services) == 0 {
		f.Log.Infof("No services found")
		return nil
	}

	header := []string{"ID", "Name", "Type", "CreatedAt"}
	rows := make([][]string, 0, len(services))

	for _, service := range services {
		row := make([]string, len(header))
		row[0] = service.ID
		row[1] = service.Name
		row[2] = service.Template
		row[3] = service.CreatedAt.Format("2006-01-02 15:04:05")
		rows = append(rows, row)
	}

	cmdutil.PrintTable(header, rows)

	return nil
}

func listServicesDetailByEnvironment(f *cmdutil.Factory, projectID, environmentID string) error {
	services, err := f.ApiClient.ListAllServicesDetailByEnvironment(context.Background(), projectID, environmentID)
	if err != nil {
		return fmt.Errorf("list services failed: %w", err)
	}

	if len(services) == 0 {
		f.Log.Infof("No services found")
		return nil
	}

	header := []string{"ID", "Name", "Status", "Domains", "Type", "GitTrigger", "CreatedAt"}
	rows := make([][]string, 0, len(services))

	for _, service := range services {
		row := make([]string, len(header))
		row[0] = service.ID
		row[1] = service.Name
		row[2] = service.Status
		domains := make([]string, len(service.Domains))
		for i, domain := range service.Domains {
			domains[i] = domain.Domain
		}
		row[3] = strings.Join(domains, ",")
		row[4] = service.Template
		gitTrigger := ""
		if service.GitTrigger != nil {
			gitTrigger = fmt.Sprintf("%s(%s)", service.GitTrigger.BranchName, service.GitTrigger.Provider)
		} else {
			gitTrigger = "None"
		}
		row[5] = gitTrigger
		row[6] = service.CreatedAt.Format("2006-01-02 15:04:05")
		rows = append(rows, row)
	}

	cmdutil.PrintTable(header, rows)

	return nil
}

func paramCheck(opts *Options) error {
	if opts.projectID == "" {
		return fmt.Errorf("project-id is required")
	}

	return nil
}
