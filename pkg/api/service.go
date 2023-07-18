package api

import (
	"context"
	"errors"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) ListServices(ctx context.Context, projectID string, skip, limit int) (*model.ServiceConnection, error) {
	skip, limit = normalizePagination(skip, limit)

	var query struct {
		Services model.ServiceConnection `graphql:"services(projectID: $projectID, skip: $skip, limit: $limit)"`
	}

	err := c.Query(ctx, &query, V{
		"projectID": ObjectID(projectID),
		"skip":      skip,
		"limit":     limit,
	})

	if err != nil {
		return nil, err
	}

	return &query.Services, nil
}

// ListAllServices returns all services owned by the current user.
func (c *client) ListAllServices(ctx context.Context, projectID string) ([]*model.Service, error) {
	skip := 0
	next := true

	var services []*model.Service

	for next {
		serviceCon, err := c.ListServices(context.Background(), projectID, skip, 100)
		if err != nil {
			return nil, err
		}
		for _, service := range serviceCon.Edges {
			services = append(services, service.Node)
		}

		skip += 5
		next = serviceCon.PageInfo.HasNextPage
	}

	return services, nil
}

func (c *client) GetService(ctx context.Context, id string, ownerName string, projectName string, name string) (*model.Service, error) {
	if id != "" {
		return c.getServiceByID(ctx, id)
	}

	if ownerName != "" && projectName != "" && name != "" {
		return c.getServiceByOwnerAndProjectAndName(ctx, ownerName, projectName, name)
	}

	return nil, errors.New("either id or ownerName, projectName, and name must be specified")
}

func (c *client) getServiceByID(ctx context.Context, id string) (*model.Service, error) {
	var query struct {
		Service model.Service `graphql:"service(_id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id": ObjectID(id),
	})

	if err != nil {
		return nil, err
	}

	return &query.Service, nil
}

func (c *client) getServiceByOwnerAndProjectAndName(ctx context.Context, ownerName string, projectName string, name string) (*model.Service, error) {
	var query struct {
		Service model.Service `graphql:"service(owner: $owner, projectName: $projectName, name: $name)"`
	}

	err := c.Query(ctx, &query, V{
		"owner":       ownerName,
		"projectName": projectName,
		"name":        name,
	})

	if err != nil {
		return nil, err
	}

	return &query.Service, nil
}
