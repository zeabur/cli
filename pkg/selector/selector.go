package selector

import (
	"context"
	"errors"
	"fmt"

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
		SelectProject() (zcontext.BasicInfo, *model.Project, error)
	}

	ServiceSelector interface {
		SelectService(projectID string) (zcontext.BasicInfo, *model.Service, error)
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

func (s *selector) SelectProject() (zcontext.BasicInfo, *model.Project, error) {
	projects, err := s.client.ListAllProjects(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("list projects failed: %w", err)
	}
	if len(projects) == 0 {
		return nil, nil, fmt.Errorf("no projects found")
	}
	if len(projects) == 1 {
		s.log.Infof("Only one project found, select <%s> automatically\n", projects[0].Name)
		project := projects[0]
		return zcontext.NewBasicInfo(project.ID, project.Name), project, nil
	}

	projectsName := make([]string, len(projects))
	for i, project := range projects {
		projectsName[i] = project.Name
	}
	projectsName = append(projectsName, "Create a new project")
	index, err := s.prompter.Select("Select project", projectsName[0], projectsName)
	if err != nil {
		return nil, nil, fmt.Errorf("select project failed: %w", err)
	}

	if index == len(projects) {

		regions, err := s.client.GetRegions(context.Background())
		if err != nil {
			return nil, nil, fmt.Errorf("get regions failed: %w", err)
		}
		regions = regions[1:]

		regionIDs := make([]string, 0, len(regions))
		for _, region := range regions {
			regionIDs = append(regionIDs, region.ID)
		}

		projectRegionIndex, err := s.prompter.Select("Select project region", "", regionIDs)
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

func (s *selector) SelectService(projectID string) (zcontext.BasicInfo, *model.Service, error) {
	services, err := s.client.ListAllServices(context.Background(), projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get services: %w", err)
	}

	if len(services) == 0 {
		return nil, nil, errors.New("there are no services in current project")
	}

	if len(services) == 1 {
		s.log.Infof("Only one service in current project, select <%s> automatically\n", services[0].Name)
		service := services[0]
		return zcontext.NewBasicInfo(service.ID, service.Name), service, nil
	}

	serviceNames := make([]string, len(services))

	for i, service := range services {
		serviceNames[i] = service.Name
	}

	index, err := s.prompter.Select("Select a service", serviceNames[0], serviceNames)
	if err != nil {
		return nil, nil, fmt.Errorf("select service failed: %w", err)
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
