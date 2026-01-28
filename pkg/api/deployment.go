package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) ListDeployments(ctx context.Context, serviceID string, environmentID string, skip, limit int) (*model.DeploymentConnection, error) {
	_, limit = normalizePagination(skip, limit)

	var query struct {
		Deployments model.DeploymentConnection `graphql:"deployments(serviceID: $serviceID, environmentID: $environmentID, perPage: $perPage)"`
	}

	err := c.Query(ctx, &query, V{
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
		"perPage":       limit,
	})
	if err != nil {
		return nil, err
	}

	return &query.Deployments, nil
}

func (c *client) ListAllDeployments(ctx context.Context, serviceID string, environmentID string) (model.Deployments, error) {
	conn, err := c.ListDeployments(ctx, serviceID, environmentID, 0, 5)
	if err != nil {
		return nil, err
	}

	deployments := make(model.Deployments, 0, len(conn.Edges))
	for _, edge := range conn.Edges {
		deployments = append(deployments, edge.Node)
	}

	return deployments, nil
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
