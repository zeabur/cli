// Package set is the subcommand to set the context for the CLI.
package set

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/zcontext"
)

type Options struct {
	id   string
	name string

	ct contextType
}

type contextType = string

const (
	project              contextType = "project"
	projectShorthand     contextType = "proj"
	environment          contextType = "environment"
	environmentShorthand contextType = "env"
	service              contextType = "service"
	serviceShorthand     contextType = "svc"
)

func NewCmdSet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "set <context-type>",
		Short: "Set Contexts(project, environment, service), either by ID or by name",
		Long: heredoc.Doc(`Set Contexts either by ID or by name,
			For example:
				zeabur context set project --id=1234567890
				zeabur context set proj --id=1234567890
				zeabur context set environment --name=dev
				zeabur context set env --id=1234567890
				zeabur context set service --name=svc1
				zeabur context set svc --name=svc1`,
		),
		Args:       cobra.ExactArgs(1),
		ValidArgs:  []string{project, environment, service},
		ArgAliases: []string{projectShorthand, environmentShorthand, serviceShorthand},
		RunE: func(cmd *cobra.Command, args []string) error {
			// the first argument is the context type
			opts.ct = args[0]
			return runSet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "ID of the project, environment, or service")
	cmd.Flags().StringVar(&opts.name, "name", "", "Name of the project, environment, or service")

	return cmd
}

func runSet(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runSetInteractive(f, opts)
	}

	return runSetNonInteractive(f, opts)
}

func runSetInteractive(f *cmdutil.Factory, opts *Options) error {
	switch opts.ct {
	case project, projectShorthand:
		return selectProject(f, opts)
	case environment, environmentShorthand:
		return selectEnvironment(f, opts)
	case service, serviceShorthand:
		return selectService(f, opts)
	}

	return fmt.Errorf("invalid context type: %s, the context type should be one of project, environment, or service", opts.ct)
}

func runSetNonInteractive(f *cmdutil.Factory, opts *Options) error {
	switch opts.ct {
	case project, projectShorthand:
		return setProject(f, opts.id, opts.name, true)
	case environment, environmentShorthand:
		return setEnvironment(f, opts.id, opts.name, true)
	case service, serviceShorthand:
		return setService(f, opts.id, opts.name, true)
	}

	return fmt.Errorf("invalid context type: %s", opts.ct)
}

func setProject(f *cmdutil.Factory, id, name string, shouldCheck bool) error {
	// if should check, it means either id or name is empty,
	// so we need to get the project first

	// if should check is false, it means both id and name are not empty,
	// and it has been checked in the previous step

	if id == "" && name == "" {
		return fmt.Errorf("either --id or --name should be specified")
	}

	if !shouldCheck && (id == "" || name == "") {
		return fmt.Errorf("invalid call to setProject, shouldCheck is false but id or name is empty")
	}

	if shouldCheck {
		ctx := context.Background()
		project, err := f.ApiClient.GetProject(ctx, id, f.Config.GetUsername(), name)
		if err != nil {
			return fmt.Errorf("failed to get project: %w", err)
		}
		f.Config.GetContext().SetProject(zcontext.NewBasicInfo(project.ID, project.Name))

	} else {
		f.Config.GetContext().SetProject(zcontext.NewBasicInfo(id, name))
	}

	// clear environment and service context when project context is set
	f.Config.GetContext().ClearService()
	f.Config.GetContext().ClearEnvironment()

	return nil
}

func setEnvironment(f *cmdutil.Factory, id, name string, shouldCheck bool) error {
	// if should check, it means either id or name is empty,
	// so we need to get the environment first

	// if should check is false, it means both id and name are not empty,
	// and it has been checked in the previous step

	if id == "" && name == "" {
		return fmt.Errorf("either --id or --name should be specified")
	}

	// we can only check environment by id, name is not supported
	// so, if shouldCheck is true, id must not be empty
	if id == "" && shouldCheck {
		return fmt.Errorf("invalid call to setEnvironment, shouldCheck is true but id is empty")
	}

	if !shouldCheck && (id == "" || name == "") {
		return fmt.Errorf("invalid call to setEnvironment, shouldCheck is false but id or name is empty")
	}

	if err := checkProjectHasBeenSet(f); err != nil {
		return err
	}

	if shouldCheck {
		ctx := context.Background()
		environment, err := f.ApiClient.GetEnvironment(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get environment: %w", err)
		}
		f.Config.GetContext().SetEnvironment(zcontext.NewBasicInfo(environment.ID, environment.Name))

	} else {
		f.Config.GetContext().SetEnvironment(zcontext.NewBasicInfo(id, name))
	}

	return nil
}

func setService(f *cmdutil.Factory, id, name string, shouldCheck bool) error {
	// if should check, it means either id or name is empty,
	// so we need to get the project first

	// if should check is false, it means both id and name are not empty,
	// and it has been checked in the previous step

	if id == "" && name == "" {
		return fmt.Errorf("either --id or --name should be specified")
	}

	if !shouldCheck && (id == "" || name == "") {
		return fmt.Errorf("invalid call to setService, shouldCheck is false but id or name is empty")
	}

	if err := checkProjectHasBeenSet(f); err != nil {
		return err
	}

	if shouldCheck {
		ctx := context.Background()
		service, err := f.ApiClient.GetService(ctx, id, f.Config.GetUsername(), f.Config.GetContext().GetProject().GetName(), name)
		if err != nil {
			return fmt.Errorf("failed to get service: %w", err)
		}
		f.Config.GetContext().SetService(zcontext.NewBasicInfo(service.ID, service.Name))

	} else {
		f.Config.GetContext().SetService(zcontext.NewBasicInfo(id, name))
	}

	return nil
}

func selectProject(f *cmdutil.Factory, opts *Options) error {
	// if flag is set, use it directly, it turns into non-interactive mode automatically
	if opts.id != "" || opts.name != "" {
		return setProject(f, opts.id, opts.name, true)
	}

	// else, show a list of projects to select
	projects, err := f.ApiClient.ListAllProjects(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	if len(projects) == 0 {
		return fmt.Errorf("no project found in your account")
	}

	var index int

	// if there is only one project, no need to ask
	if len(projects) == 1 {
		f.Log.Info("There is only one project in your account, select it automatically")
		index = 0
	} else {
		projectNames := make([]string, len(projects))
		for i, project := range projects {
			projectNames[i] = project.Name
		}

		index, err = f.Prompter.Select("Select a project", projects[0].Name, projectNames)
		if err != nil {
			return err
		}

		f.Log.Info("Project selected: ")
	}

	logProject(f, projects[index])

	confirm := true

	// if project is already set, we need to clear the environment and service, and set the project.
	// So we need to ask user to confirm.
	if !f.Config.GetContext().GetProject().Empty() {
		oldProject := f.Config.GetContext().GetProject().GetName()
		prompt := fmt.Sprintf("Project is already set(%s), do you want to change it?"+
			"(Once changed, the environment and service will be cleared.)", oldProject)
		confirm, err = f.Prompter.Confirm(prompt, true)
		if err != nil {
			return fmt.Errorf("failed to confirm: %w", err)
		}
	}

	if confirm {
		err = setProject(f, projects[index].ID, projects[index].Name, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func selectEnvironment(f *cmdutil.Factory, opts *Options) error {
	if err := checkProjectHasBeenSet(f); err != nil {
		return err
	}

	// if flag is set, use it directly, it turns into non-interactive mode automatically
	if opts.id != "" || opts.name != "" {
		return setEnvironment(f, opts.id, opts.name, true)
	}

	var (
		ctx       = context.Background()
		projectID = f.Config.GetContext().GetProject().GetID()
	)

	environments, err := f.ApiClient.ListEnvironments(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get environments of project %s: %w", projectID, err)
	}

	if len(environments) == 0 {
		return fmt.Errorf("there are no environments in project %s", projectID)
	}

	var index int

	if len(environments) == 1 {
		f.Log.Info("There is only one environment in your project, select it automatically")
		index = 0
	} else {
		environmentNames := make([]string, len(environments))
		for i, environment := range environments {
			environmentNames[i] = environment.Name
		}

		index, err = f.Prompter.Select("Select an environment", environments[0].Name, environmentNames)
		if err != nil {
			return err
		}
	}

	logEnvironment(f, environments[index])

	err = setEnvironment(f, environments[index].ID, environments[index].Name, false)
	if err != nil {
		return err
	}

	return nil
}

func selectService(f *cmdutil.Factory, opts *Options) error {
	if err := checkProjectHasBeenSet(f); err != nil {
		return err
	}

	// if flag is set, use it directly, it turns into non-interactive mode automatically
	if opts.id != "" || opts.name != "" {
		return setService(f, opts.id, opts.name, true)
	}

	projectID := f.Config.GetContext().GetProject().GetID()

	services, err := f.ApiClient.ListAllServices(context.Background(), projectID)
	if err != nil {
		return fmt.Errorf("failed to get services of project %s, %w", projectID, err)
	}

	if len(services) == 0 {
		return fmt.Errorf("there are no services in project %s", projectID)
	}

	var index int

	if len(services) == 1 {
		f.Log.Info("There is only one service in your project, select it automatically")
		index = 0
	} else {
		serviceNames := make([]string, len(services))
		for i, service := range services {
			serviceNames[i] = service.Name
		}

		index, err = f.Prompter.Select("Select a service", services[0].Name, serviceNames)
		if err != nil {
			return err
		}
	}

	logService(f, services[index])

	err = setService(f, services[index].ID, services[index].Name, false)
	if err != nil {
		return err
	}

	return nil
}

func checkProjectHasBeenSet(f *cmdutil.Factory) error {
	if f.Config.GetContext().GetProject().Empty() {
		return fmt.Errorf("you must set project context first")
	}
	return nil
}

// todo: pretty print
func logProject(f *cmdutil.Factory, project *model.Project) {
	f.Log.Info(project)
}

// todo: pretty print
func logEnvironment(f *cmdutil.Factory, environment *model.Environment) {
	f.Log.Info(environment)
}

func logService(f *cmdutil.Factory, service *model.Service) {
	f.Log.Info(service)
}
