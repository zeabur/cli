package api

import (
	"context"
	"github.com/zeabur/cli/pkg/model"
)

func (c *client) ListVariables(ctx context.Context, serviceID string, environmentID string) (model.Variables, error) {
	var query struct {
		Service struct {
			Variables model.Variables `graphql:"variables(environmentID: $environmentID, exposed: true)"`
		} `graphql:"service(_id: $serviceID)"`
	}

	err := c.Query(ctx, &query, V{
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
	})

	if err != nil {
		return nil, err
	}

	return query.Service.Variables, nil
}

func (c *client) UpdateVariables(ctx context.Context, serviceID string, environmentID string, data map[string]string) (bool, error) {
	var mutation struct {
		UpdateVariables struct {
			UpdateEnvironmentVariable bool `graphql:"updateEnvironmentVariable"`
		} `graphql:"updateEnvironmentVariable(environmentID: $environmentID, serviceID: $serviceID, data: $data)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"environmentID": ObjectID(environmentID),
		"serviceID":     ObjectID(serviceID),
		"data":          data,
	})

	if err != nil {
		return false, err
	}

	return mutation.UpdateVariables.UpdateEnvironmentVariable, nil
}
