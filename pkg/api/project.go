package api

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zeabur/cli/pkg/model"
)

// ListProjects returns projects owned by the current user.
// Note: the backend hasn't implemented pagination yet, currently we return all projects at once.
func (c *client) ListProjects(ctx context.Context, skip, limit int) (*model.ProjectConnection, error) {
	if skip < 0 {
		skip = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 5
	}

	var query struct {
		Projects model.ProjectConnection `graphql:"projects(skip: $skip, limit: $limit)"`
	}

	err := c.Query(ctx, &query, V{
		"skip":  skip,
		"limit": limit,
	})
	if err != nil {
		return nil, err
	}

	return &query.Projects, nil
}

// GetProject returns a project by (its ID), or (owner name and project name).
func (c *client) GetProject(ctx context.Context, id primitive.ObjectID, ownerName string, name string) (*model.Project, error) {
	var query struct {
		Project model.Project `graphql:"project(id: $id, owner: $owner, name: $name)"`
	}

	err := c.Query(ctx, &query, V{
		"_id":   id,
		"owner": ownerName,
		"name":  name,
	})
	if err != nil {
		return nil, err
	}

	return &query.Project, nil
}
