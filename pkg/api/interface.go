package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

// Client is the interface of the Zeabur API client.
type Client interface {
	GetUserInfo(ctx context.Context) (*model.User, error)

	ListProjects(ctx context.Context, skip, limit int) (*model.ProjectConnection, error)
	ListAllProjects(ctx context.Context) ([]*model.Project, error)
	GetProject(ctx context.Context, id string, ownerName string, name string) (*model.Project, error)

	ListEnvironments(ctx context.Context, projectID string) ([]*model.Environment, error)
	GetEnvironment(ctx context.Context, id string) (*model.Environment, error)

	ListServices(ctx context.Context, projectID string, skip, limit int) (*model.ServiceConnection, error)
	ListAllServices(ctx context.Context, projectID string) ([]*model.Service, error)
	GetService(ctx context.Context, id string, ownerName string, projectName string, name string) (*model.Service, error)
	ExposeService(ctx context.Context, id string, environmentID string, projectID string, name string) (*model.TempTCPPort, error)
}
