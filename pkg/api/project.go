package api

import (
	"context"

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
func (c *client) GetProject(ctx context.Context, id string, ownerUsername string, projectName string) (*model.Project, error) {
	if id == "" {
		return c.getProjectByOwnerUsernameAndProject(ctx, ownerUsername, projectName)
	}

	return c.getProjectByID(ctx, id)
}

func (c *client) getProjectByID(ctx context.Context, id string) (*model.Project, error) {
	var query struct {
		Project model.Project `graphql:"project(_id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id": ObjectID(id),
	})

	if err != nil {
		return nil, err
	}

	return &query.Project, nil
}

func (c *client) getProjectByOwnerUsernameAndProject(ctx context.Context,
	ownerUsername string, projectName string) (*model.Project, error) {
	var query struct {
		Project model.Project `graphql:"project(owner: $owner, name: $name)"`
	}

	err := c.Query(ctx, &query, V{
		"owner": ownerUsername,
		"name":  projectName,
	})

	if err != nil {
		return nil, err
	}

	return &query.Project, nil
}
