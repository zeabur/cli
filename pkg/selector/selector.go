package selector

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/prompt"
	"github.com/zeabur/cli/pkg/zcontext"
	"go.uber.org/zap"
)

type (
	Selector interface {
		ProjectSelector
		ServiceSelector
		EnvironmentSelector
	}

	ProjectSelector interface {
		SelectProject(opts ...SelectProjectOptionsApplyer) (zcontext.BasicInfo, *model.Project, error)
	}

	ServiceSelector interface {
		SelectService(opt SelectServiceOptions) (zcontext.BasicInfo, *model.Service, error)
	}

	EnvironmentSelector interface {
		SelectEnvironment(projectID string) (zcontext.BasicInfo, *model.Environment, error)
	}
)

type selector struct {
	client   api.Client
	log      *zap.SugaredLogger
	prompter prompt.Prompter
}

func New(client api.Client, log *zap.SugaredLogger, prompter prompt.Prompter) Selector {
	return &selector{
		client:   client,
		log:      log,
		prompter: prompter,
	}
}

type SelectProjectOptions struct {
	// CreatePreferred selects "Create New Project…" by default,
	// and not auto-select the only project.
	CreatePreferred bool
}

type SelectProjectOptionsApplyer func(*SelectProjectOptions)

func WithCreatePreferred() SelectProjectOptionsApplyer {
	return func(opt *SelectProjectOptions) {
		opt.CreatePreferred = true
	}
}

func (s *selector) SelectProject(opts ...SelectProjectOptionsApplyer) (zcontext.BasicInfo, *model.Project, error) {
	options := SelectProjectOptions{
		CreatePreferred: false,
	}
	for _, applyer := range opts {
		applyer(&options)
	}

	projects, err := s.client.ListAllProjects(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("list projects failed: %w", err)
	}
	if len(projects) == 1 && !options.CreatePreferred {
		s.log.Infof("Only one project found, select <%s> automatically\n", projects[0].Name)
		project := projects[0]
		return zcontext.NewBasicInfo(project.ID, project.Name), project, nil
	}

	projectsName := make([]string, len(projects))
	for i, project := range projects {
		projectsName[i] = project.Name
	}
	projectsName = append(projectsName, "Create a new project")

	defaultChoice := projectsName[0]
	if options.CreatePreferred {
		defaultChoice = projectsName[len(projectsName)-1] // the final item = create project
	}

	index, err := s.prompter.Select("Select project", defaultChoice, projectsName)
	if err != nil {
		return nil, nil, fmt.Errorf("select project failed: %w", err)
	}

	if index == len(projects) {

		regions, err := s.client.GetRegions(context.Background())
		if err != nil {
			return nil, nil, fmt.Errorf("get regions failed: %w", err)
		}
		regions = regions[1:]

		regionOptions := make([]string, 0, len(regions))
		for _, region := range regions {
			regionOptions = append(regionOptions, fmt.Sprintf("%s (%s)", region.Description, region.Name))
		}

		projectRegionIndex, err := s.prompter.Select("Select project region", "", regionOptions)
		if err != nil {
			return nil, nil, fmt.Errorf("select project region failed: %w", err)
		}

		projectRegion := regions[projectRegionIndex].ID

		project, err := s.client.CreateProject(context.Background(), projectRegion, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("create project failed: %w", err)
		}

		return zcontext.NewBasicInfo(project.ID, project.Name), project, nil
	}

	project := projects[index]

	return zcontext.NewBasicInfo(project.ID, project.Name), project, nil
}

type SelectServiceOptions struct {
	// ProjectID is the project id to select service from
	ProjectID string
	// Auto select the only service in the project
	Auto bool
	// CreateNew shows an option to create a new service
	CreateNew bool
	// FilterFunc filters the services
	FilterFunc func(service *model.Service) bool
}

func (s *selector) SelectService(opt SelectServiceOptions) (zcontext.BasicInfo, *model.Service, error) {
	projectID := opt.ProjectID
	auto := opt.Auto

	services, err := s.client.ListAllServices(context.Background(), projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get services: %w", err)
	}

	if opt.FilterFunc != nil {
		filtered := make([]*model.Service, 0, len(services))
		for _, service := range services {
			if opt.FilterFunc(service) {
				filtered = append(filtered, service)
			}
		}
		services = filtered
	}

	if len(services) == 0 {
		return nil, nil, nil
	}

	if len(services) == 1 && auto {
		s.log.Infof("Only one service in current project, select <%s> automatically\n", services[0].Name)
		service := services[0]
		return zcontext.NewBasicInfo(service.ID, service.Name), service, nil
	}

	serviceNames := make([]string, len(services))
	for i, service := range services {
		serviceNames[i] = service.Name
	}

	if opt.CreateNew {
		serviceNames = append(serviceNames, "Create a new service")
	}

	index, err := s.prompter.Select("Select a service", serviceNames[0], serviceNames)
	if err != nil {
		return nil, nil, fmt.Errorf("select service failed: %w", err)
	}

	if index == len(services) {
		from := []rune("abcdefghijklmnopqrstuvwxyz")
		b := make([]rune, 8)
		for i := range b {
			b[i] = from[rand.Intn(len(from))]
		}
		serviceName := string(b)
		service, err := s.client.CreateEmptyService(context.Background(), projectID, serviceName)
		if err != nil {
			return nil, nil, fmt.Errorf("create service failed: %w", err)
		}
		return zcontext.NewBasicInfo(service.ID, service.Name), service, nil
	}

	service := services[index]

	return zcontext.NewBasicInfo(service.ID, service.Name), service, nil
}

func (s *selector) SelectEnvironment(projectID string) (zcontext.BasicInfo, *model.Environment, error) {
	environments, err := s.client.ListEnvironments(context.Background(), projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get environments: %w", err)
	}

	if len(environments) == 0 {
		return nil, nil, errors.New("there are no environments in current project")
	}

	var index int

	if len(environments) == 1 {
		s.log.Infof("Only one environment in current project, select <%s> automatically\n", environments[0].Name)
		index = 0
	} else {
		environmentNames := make([]string, len(environments))
		for i, environment := range environments {
			environmentNames[i] = environment.Name
		}

		index, err = s.prompter.Select("Select an environment", environments[0].Name, environmentNames)
		if err != nil {
			return nil, nil, fmt.Errorf("select environment failed: %w", err)
		}
	}

	environment := environments[index]

	return zcontext.NewBasicInfo(environment.ID, environment.Name), environment, nil
}
