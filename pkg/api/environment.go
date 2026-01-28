package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) ListEnvironments(ctx context.Context, projectID string) (model.Environments, error) {
	var query struct {
		Environments []*model.Environment `graphql:"environments(projectID: $projectID)"`
	}

	err := c.Query(ctx, &query, V{
		"projectID": ObjectID(projectID),
	})
	if err != nil {
		return nil, err
	}

	return query.Environments, nil
}

func (c *client) GetEnvironment(ctx context.Context, id string) (*model.Environment, error) {
	var query struct {
		Environment model.Environment `graphql:"environment(_id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id": ObjectID(id),
	})
	if err != nil {
		return nil, err
	}

	return &query.Environment, nil
}
