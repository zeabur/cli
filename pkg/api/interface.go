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
	GitAPI
	TemplateAPI
	DomainAPI
}

type (
	UserAPI interface {
		GetUserInfo(ctx context.Context) (*model.User, error)
	}

	ProjectAPI interface {
		ListProjects(ctx context.Context, skip, limit int) (*model.ProjectConnection, error)
		ListAllProjects(ctx context.Context) (model.Projects, error)
		GetProject(ctx context.Context, id string, ownerName string, name string) (*model.Project, error)
		CreateProject(ctx context.Context, region string, name *string) (*model.Project, error)
		DeleteProject(ctx context.Context, id string) error

		GetRegions(ctx context.Context) ([]model.Region, error)
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
		ServiceMetric(ctx context.Context, id, projectID, environmentID, metricType string, startTime, endTime time.Time) (*model.ServiceMetric, error)
		GetPrebuiltItems(ctx context.Context) ([]model.PrebuiltItem, error)
		SearchGitRepositories(ctx context.Context, keyword *string) ([]model.GitRepo, error)

		RestartService(ctx context.Context, id string, environmentID string) error
		RedeployService(ctx context.Context, id string, environmentID string) error
		SuspendService(ctx context.Context, id string, environmentID string) error
		ExposeService(ctx context.Context, id string, environmentID string, projectID string, name string) (*model.TempTCPPort, error)
		CreatePrebuiltService(ctx context.Context, projectID string, marketplaceCode string) (*model.Service, error)
		CreateService(ctx context.Context, projectID string, name string, repoID int, branchName string) (*model.Service, error)
		CreateEmptyService(ctx context.Context, projectID string, name string) (*model.Service, error)
		UploadZipToService(ctx context.Context, projectID string, serviceID string, environmentID string, zipBytes []byte) (*model.Service, error)
	}

	DomainAPI interface {
		AddDomain(ctx context.Context, serviceID string, environmentID string, isGenerated bool, domain string, options ...string) (*string, error)
		ListDomains(ctx context.Context, serviceID string, environmentID string) (model.Domains, error)
		RemoveDomain(ctx context.Context, domain string) (bool, error)
		CheckDomainAvailable(ctx context.Context, domain string, isGenerated bool) (bool, string, error)
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

		WatchRuntimeLogs(ctx context.Context, deploymentID string) (<-chan model.Log, error)
		WatchBuildLogs(ctx context.Context, deploymentID string) (<-chan model.Log, error)
	}

	GitAPI interface {
		GetRepoBranches(ctx context.Context, repoOwner string, repoName string) ([]string, error)
		GetRepoID(repoOwner string, repoName string) (int, error)
		GetRepoInfo() (string, string, error)
		GetRepoBranchesByRepoID(repoID int) ([]string, error)
	}

	TemplateAPI interface {
		ListTemplates(ctx context.Context, skip, limit int) (*model.TemplateConnection, error)
		ListAllTemplates(ctx context.Context) (model.Templates, error)
		GetTemplate(ctx context.Context, code string) (*model.Template, error)

		DeployTemplate(ctx context.Context, code string, variables map[string]interface{}, repoConfigs model.RepoConfigs) (*model.Project, error)
		DeleteTemplate(ctx context.Context, code string) error
	}
)
