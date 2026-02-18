package api

import (
	"context"
	"fmt"

	"github.com/hasura/go-graphql-client/pkg/jsonutil"
	"github.com/zeabur/cli/pkg/model"
)

func (c *client) GetRuntimeLogs(ctx context.Context, serviceID, environmentID, deploymentID string) (model.Logs, error) {
	if serviceID == "" {
		return nil, fmt.Errorf("serviceID is required for runtime logs")
	}

	if deploymentID != "" {
		var query struct {
			Logs model.Logs `graphql:"runtimeLogs(serviceID: $serviceID, environmentID: $environmentID, deploymentID: $deploymentID)"`
		}

		err := c.Query(ctx, &query, V{
			"serviceID":     ObjectID(serviceID),
			"environmentID": ObjectID(environmentID),
			"deploymentID":  ObjectID(deploymentID),
		})
		if err != nil {
			return nil, err
		}

		return query.Logs, nil
	}

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

func (c *client) WatchRuntimeLogs(ctx context.Context, projectID, serviceID, environmentID, deploymentID string) (<-chan model.Log, error) {
	if deploymentID != "" {
		return c.watchRuntimeLogsWithDeployment(projectID, serviceID, environmentID, deploymentID)
	}
	return c.watchRuntimeLogs(projectID, serviceID, environmentID)
}

func (c *client) watchRuntimeLogs(projectID, serviceID, environmentID string) (<-chan model.Log, error) {
	logs := make(chan model.Log, 100)
	subClient := c.sub

	type subscription struct {
		Log model.Log `graphql:"runtimeLogReceived(projectID: $projectID, serviceID: $serviceID, environmentID: $environmentID)"`
	}

	sub := subscription{}

	_, err := subClient.Subscribe(&sub, V{
		"projectID":     ObjectID(projectID),
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
	}, func(dataValue []byte, errValue error) error {
		if errValue != nil {
			fmt.Println(errValue)
			return nil
		}

		if dataValue == nil {
			return nil
		}

		data := subscription{}
		err := jsonutil.UnmarshalGraphQL(dataValue, &data)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		logs <- data.Log
		return nil
	})
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(logs)
		_ = subClient.Run()
	}()

	return logs, nil
}

func (c *client) watchRuntimeLogsWithDeployment(projectID, serviceID, environmentID, deploymentID string) (<-chan model.Log, error) {
	logs := make(chan model.Log, 100)
	subClient := c.sub

	type subscription struct {
		Log model.Log `graphql:"runtimeLogReceived(projectID: $projectID, serviceID: $serviceID, environmentID: $environmentID, deploymentID: $deploymentID)"`
	}

	sub := subscription{}

	_, err := subClient.Subscribe(&sub, V{
		"projectID":     ObjectID(projectID),
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
		"deploymentID":  ObjectID(deploymentID),
	}, func(dataValue []byte, errValue error) error {
		if errValue != nil {
			fmt.Println(errValue)
			return nil
		}

		if dataValue == nil {
			return nil
		}

		data := subscription{}
		err := jsonutil.UnmarshalGraphQL(dataValue, &data)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		logs <- data.Log
		return nil
	})
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(logs)
		_ = subClient.Run()
	}()

	return logs, nil
}

func (c *client) WatchBuildLogs(ctx context.Context, projectID, deploymentID string) (<-chan model.Log, error) {
	logs := make(chan model.Log, 100)
	subClient := c.sub

	type subscription struct {
		Log model.Log `graphql:"buildLogReceived(projectID: $projectID, deploymentID: $deploymentID)"`
	}

	sub := subscription{}

	_, err := subClient.Subscribe(&sub, V{
		"projectID":    ObjectID(projectID),
		"deploymentID": ObjectID(deploymentID),
	}, func(dataValue []byte, errValue error) error {
		if errValue != nil {
			fmt.Println(errValue)
			return nil
		}

		if dataValue == nil {
			return nil
		}

		data := subscription{}
		err := jsonutil.UnmarshalGraphQL(dataValue, &data)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		logs <- data.Log
		return nil
	})
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(logs)
		_ = subClient.Run()
	}()

	return logs, nil
}
