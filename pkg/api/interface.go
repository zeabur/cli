package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

// Client is the interface of the Zeabur API client.
type Client interface {
	GetUserInfo(ctx context.Context) (*model.User, error)

	ListProjects(ctx context.Context, skip, limit int) (*model.ProjectConnection, error)
	GetProject(ctx context.Context, id string, ownerName string, name string) (*model.Project, error)
}
