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
	DeploymentAPI
	LogAPI
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
	}

	EnvironmentAPI interface {
		ListEnvironments(ctx context.Context, projectID string) (model.Environments, error)
		GetEnvironment(ctx context.Context, id string) (*model.Environment, error)
	}

	ServiceAPI interface {
		ListServices(ctx context.Context, projectID string, skip, limit int) (*model.Connection[model.Service], error)
		ListAllServices(ctx context.Context, projectID string) (model.Services, error)
		ListServicesDetailByEnvironment(ctx context.Context, projectID, environmentID string, skip, limit int) (*model.Connection[model.ServiceDetail], error)
		ListAllServicesDetailByEnvironment(ctx context.Context, projectID, environmentID string) (model.ServiceDetails, error)
		GetService(ctx context.Context, id, ownerName, projectName, name string) (*model.Service, error)
		GetServiceDetailByEnvironment(ctx context.Context, id, ownerName, projectName, name, environmentID string) (*model.ServiceDetail, error)
		ServiceMetric(ctx context.Context, id, environmentID, metricType string, startTime, endTime time.Time) (*model.ServiceMetric, error)
		GetMarketplaceItems(ctx context.Context) ([]model.MarketplaceItem, error)

		RestartService(ctx context.Context, id string, environmentID string) error
		RedeployService(ctx context.Context, id string, environmentID string) error
		SuspendService(ctx context.Context, id string, environmentID string) error
		ExposeService(ctx context.Context, id string, environmentID string, projectID string, name string) (*model.TempTCPPort, error)
		CreateServiceFromMarketplace(ctx context.Context, projectID string, name string, itemCode string) (*model.Service, error)
	}

	DeploymentAPI interface {
		ListDeployments(ctx context.Context, serviceID string, environmentID string, skip, limit int) (*model.Connection[model.Deployment], error)
		ListAllDeployments(ctx context.Context, serviceID string, environmentID string) (model.Deployments, error)
		GetDeployment(ctx context.Context, id string) (*model.Deployment, error)
		GetLatestDeployment(ctx context.Context, serviceID string, environmentID string) (*model.Deployment, bool, error)
	}

	LogAPI interface {
		// GetRuntimeLogs returns the logs of a service, two cases of parameters:
		// 1. only deploymentID
		// 2. deploymentID and serviceID
		GetRuntimeLogs(ctx context.Context, deploymentID, serviceID, environmentID string) (model.Logs, error)
		GetBuildLogs(ctx context.Context, deploymentID string) (model.Logs, error)

		WatchRuntimeLogs(ctx context.Context, deploymentID, serviceID, environmentID string) (<-chan model.Log, error)
		WatchBuildLogs(ctx context.Context, deploymentID string) (<-chan model.Log, error)
	}
)
