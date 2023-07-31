package api

import (
	"context"
	"time"

	"github.com/zeabur/cli/pkg/model"
)

// Client is the interface of the Zeabur API client.
type Client interface {
	UserAPI
	ProjectAPI
	ServiceAPI
	EnvironmentAPI
}

type (
	UserAPI interface {
		GetUserInfo(ctx context.Context) (*model.User, error)
	}

	ProjectAPI interface {
		ListProjects(ctx context.Context, skip, limit int) (*model.ProjectConnection, error)
		ListAllProjects(ctx context.Context) (model.Projects, error)
		GetProject(ctx context.Context, id string, ownerName string, name string) (*model.Project, error)
		CreateProject(ctx context.Context, name string) (*model.Project, error)
		DeleteProject(ctx context.Context, id string) error
	}

	EnvironmentAPI interface {
		ListEnvironments(ctx context.Context, projectID string) (model.Environments, error)
		GetEnvironment(ctx context.Context, id string) (*model.Environment, error)
	}

	ServiceAPI interface {
		ListServices(ctx context.Context, projectID string, skip, limit int) (*model.ServiceConnection, error)
		ListAllServices(ctx context.Context, projectID string) (model.Services, error)
		ListServicesDetailByEnvironment(ctx context.Context, projectID, environmentID string, skip, limit int) (*model.ServiceDetailConnection, error)
		ListAllServicesDetailByEnvironment(ctx context.Context, projectID, environmentID string) (model.ServiceDetails, error)
		GetService(ctx context.Context, id string, ownerName string, projectName string, name string) (*model.Service, error)
		ServiceMetric(ctx context.Context, id, environmentID, metricType string, startTime, endTime time.Time) (*model.ServiceMetric, error)

		RestartService(ctx context.Context, id string, environmentID string) error
		RedeployService(ctx context.Context, id string, environmentID string) error
		SuspendService(ctx context.Context, id string, environmentID string) error
		ExposeService(ctx context.Context, id string, environmentID string, projectID string, name string) (*model.TempTCPPort, error)
	}
)
