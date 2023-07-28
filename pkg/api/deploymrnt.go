package api

import (
	"context"
	"github.com/zeabur/cli/pkg/model"
)

func (c *client) ListDeployments(ctx context.Context, serviceID string, environmentID string, skip, limit int) (*model.Connection[model.Deployment], error) {
	skip, limit = normalizePagination(skip, limit)

	var query struct {
		Deployments model.Connection[model.Deployment] `graphql:"deployments(serviceID: $serviceID, environmentID: $environmentID, skip: $skip, limit: $limit)"`
	}

	err := c.Query(ctx, &query, V{
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
		"skip":          skip,
		"limit":         limit,
	})

	if err != nil {
		return nil, err
	}

	return &query.Deployments, nil
}

func (c *client) ListAllDeployments(ctx context.Context, serviceID string, environmentID string) (model.Deployments, error) {
	query := func(skip, limit int) (*model.Connection[model.Deployment], error) {
		return c.ListDeployments(ctx, serviceID, environmentID, skip, limit)
	}

	return listAll(query)
}

func (c *client) GetDeployment(ctx context.Context, id string) (*model.Deployment, error) {
	var query struct {
		Deployment *model.Deployment `graphql:"deployment(id: $id)"`
	}

	err := c.Query(ctx, &query, V{"id": ObjectID(id)})
	if err != nil {
		return nil, err
	}

	return query.Deployment, nil
}

func (c *client) GetLatestDeployment(ctx context.Context, serviceID string, environmentID string) (*model.Deployment, bool, error) {
	deployments, err := c.ListDeployments(ctx, serviceID, environmentID, 0, 1)
	if err != nil {
		return nil, false, err
	}

	if len(deployments.Edges) == 0 {
		return nil, false, nil
	}

	return deployments.Edges[0].Node, true, nil
}
