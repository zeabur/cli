package api

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zeabur/cli/pkg/model"
)

// V means graphql variables, it's a alias of map[string]interface{}
type V map[string]interface{}

// Client is the interface of the Zeabur API client.
type Client interface {
	GetUserInfo(ctx context.Context) (*model.User, error)

	ListProjects(ctx context.Context, skip, limit int) (*model.ProjectConnection, error)
	GetProject(ctx context.Context, id primitive.ObjectID, ownerName string, name string) (*model.Project, error)
}
