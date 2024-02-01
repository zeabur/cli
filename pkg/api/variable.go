package api

import (
	"context"
	"github.com/zeabur/cli/pkg/model"
)

func (c *client) ListVariables(ctx context.Context, serviceID string, environmentID string) (model.Variables, model.Variables, error) {
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
		return nil, nil, err
	}

	variableList := make(model.Variables, 0, len(query.Service.Variables))
	readonlyVariableList := make(model.Variables, 0, len(query.Service.Variables))
	for _, variable := range query.Service.Variables {
		if variable.ServiceID == serviceID {
			variableList = append(variableList, variable)
		} else {
			readonlyVariableList = append(readonlyVariableList, variable)
		}
	}

	return variableList, readonlyVariableList, nil
}

func (c *client) UpdateVariables(ctx context.Context, serviceID string, environmentID string, data map[string]string) (bool, error) {
	var mutation struct {
		UpdateEnvironmentVariable bool `graphql:"updateEnvironmentVariable(environmentID: $environmentID, serviceID: $serviceID, data: $data)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"environmentID": ObjectID(environmentID),
		"serviceID":     ObjectID(serviceID),
		"data":          MapString(data),
	})

	if err != nil {
		return false, err
	}

	return mutation.UpdateEnvironmentVariable, nil
}
