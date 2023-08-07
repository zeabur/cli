package api

import (
	"context"
	"fmt"

	"github.com/hasura/go-graphql-client/pkg/jsonutil"
	"github.com/zeabur/cli/pkg/model"
)

func (c *client) GetRuntimeLogs(ctx context.Context, deploymentID, serviceID, environmentID string) (model.Logs, error) {
	if deploymentID != "" {
		return c.getRuntimeLogsByDeploymentID(ctx, deploymentID)
	}

	if serviceID != "" && environmentID != "" {
		return c.getRuntimeLogsByServiceIDAndEnvironmentID(ctx, serviceID, environmentID)
	}

	return nil, fmt.Errorf("invalid arguments")
}

func (c *client) getRuntimeLogsByDeploymentID(ctx context.Context, deploymentID string) (model.Logs, error) {
	var query struct {
		Logs model.Logs `graphql:"runtimeLogs(deploymentID: $deploymentID)"`
	}

	err := c.Query(ctx, &query, V{
		"deploymentID": ObjectID(deploymentID),
	})

	fmt.Println("query", query)

	if err != nil {
		return nil, err
	}

	return query.Logs, nil
}

func (c *client) getRuntimeLogsByServiceIDAndEnvironmentID(ctx context.Context, serviceID, environmentID string) (model.Logs, error) {
	var query struct {
		Logs model.Logs `graphql:"runtimeLogs(serviceID: $serviceID, environmentID: $environmentID)"`
	}

	err := c.Query(ctx, &query, V{
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
	})

	if err != nil {
		return nil, err
	}

	return query.Logs, nil
}

func (c *client) GetBuildLogs(ctx context.Context, deploymentID string) (model.Logs, error) {
	var query struct {
		Logs model.Logs `graphql:"buildLogs(deploymentID: $deploymentID)"`
	}

	err := c.Query(ctx, &query, V{
		"deploymentID": ObjectID(deploymentID),
	})

	if err != nil {
		return nil, err
	}

	return query.Logs, nil
}

func (c *client) WatchRuntimeLogs(ctx context.Context, deploymentID, serviceID, environmentID string) (<-chan model.Log, error) {
	logs := make(chan model.Log)
	//todo: implement
	return logs, nil
}

// todo: to be tested and implemented
func (c *client) WatchBuildLogs(ctx context.Context, deploymentID string) (<-chan model.Log, error) {
	logs := make(chan model.Log, 100)
	type query struct {
		Log model.Log `graphql:"buildLogReceived(deploymentID: $deploymentID)"`
	}

	q := query{}

	handler := func(dataValue []byte, errValue error) error {
		if errValue != nil {
			// handle error
			// if returns error, it will failback to `onError` event
			fmt.Println(errValue)
			return nil
		}
		data := query{}
		// use the github.com/hasura/go-graphql-client/pkg/jsonutil package
		err := jsonutil.UnmarshalGraphQL(dataValue, &data)
		if err != nil {
			// handle error
			// if returns error, it will failback to `onError` event
			fmt.Println(err)
			return nil
		}

		logs <- data.Log

		return nil
	}

	subID, err := c.sub.Subscribe(&q, V{
		"deploymentID": ObjectID(deploymentID),
	}, handler)

	if err != nil {
		return nil, err
	}

	// todo: handle unsubscribe
	fmt.Println(subID)

	return logs, nil
}
