package api

import (
	"context"
	"errors"
	"time"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) ListServices(ctx context.Context, projectID string, skip, limit int) (*model.Connection[model.Service], error) {

	var query struct {
		Services *model.Connection[model.Service] `graphql:"services(projectID: $projectID, skip: $skip, limit: $limit)"`
	}

	err := c.Query(ctx, &query, V{
		"projectID": ObjectID(projectID),
		"skip":      skip,
		"limit":     limit,
	})

	if err != nil {
		return nil, err
	}

	return query.Services, nil
}

// ListAllServices returns all services owned by the current user.
func (c *client) ListAllServices(ctx context.Context, projectID string) (model.Services, error) {
	query := func(skip, limit int) (*model.Connection[model.Service], error) {
		return c.ListServices(ctx, projectID, skip, limit)
	}

	return listAll(query)
}

func (c *client) ListServicesDetailByEnvironment(ctx context.Context, projectID, environmentID string,
	skip, limit int) (*model.Connection[model.ServiceDetail], error) {
	skip, limit = normalizePagination(skip, limit)

	var query struct {
		Services *model.Connection[model.ServiceDetail] `graphql:"services(projectID: $projectID, skip: $skip, limit: $limit)"`
	}

	err := c.Query(ctx, &query, V{
		"projectID":     ObjectID(projectID),
		"environmentID": ObjectID(environmentID),
		"skip":          skip,
		"limit":         limit,
	})

	if err != nil {
		return nil, err
	}

	return query.Services, nil
}

func (c *client) ListAllServicesDetailByEnvironment(ctx context.Context, projectID, environmentID string) (model.ServiceDetails, error) {
	query := func(skip, limit int) (*model.Connection[model.ServiceDetail], error) {
		return c.ListServicesDetailByEnvironment(ctx, projectID, environmentID, skip, limit)
	}

	return listAll(query)
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

func (c *client) GetServiceDetailByEnvironment(ctx context.Context, id, ownerName, projectName, name, environmentID string) (*model.ServiceDetail, error) {
	if id != "" {
		return c.getServiceDetailByEnvironmentByID(ctx, id, environmentID)
	}

	if ownerName != "" && projectName != "" && environmentID != "" {
		return c.getServiceDetailByEnvironmentByOwnerAndProjectAndName(ctx, ownerName, projectName, name, environmentID)
	}

	return nil, errors.New("either id or ownerName, projectName, and environmentID must be specified")
}

func (c *client) getServiceDetailByEnvironmentByID(ctx context.Context, id string, environmentID string) (*model.ServiceDetail, error) {
	var query struct {
		Service *model.ServiceDetail `graphql:"service(_id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id":            ObjectID(id),
		"environmentID": ObjectID(environmentID),
	})

	if err != nil {
		return nil, err
	}

	return query.Service, nil
}

func (c *client) getServiceDetailByEnvironmentByOwnerAndProjectAndName(ctx context.Context, ownerName string, projectName string, name string, environmentID string) (*model.ServiceDetail, error) {
	var query struct {
		Service *model.ServiceDetail `graphql:"service(owner: $owner, projectName: $projectName, name: $name)"`
	}

	err := c.Query(ctx, &query, V{
		"owner":         ownerName,
		"projectName":   projectName,
		"name":          name,
		"environmentID": ObjectID(environmentID),
	})

	if err != nil {
		return nil, err
	}

	return query.Service, nil
}

func (c *client) ServiceMetric(ctx context.Context, id, environmentID, metricType string, startTime, endTime time.Time) (*model.ServiceMetric, error) {
	var query struct {
		ServiceMetric model.ServiceMetric `graphql:"service(_id: $serviceID)"`
	}

	err := c.Query(ctx, &query, V{
		"serviceID":     ObjectID(id),
		"environmentID": ObjectID(environmentID),
		"metricType":    model.MetricType(metricType),
		"startTime":     startTime,
		"endTime":       endTime,
	})

	if err != nil {
		return nil, err
	}

	return &query.ServiceMetric, nil
}

func (c *client) RestartService(ctx context.Context, id string, environmentID string) error {
	var mutation struct {
		RestartService bool `graphql:"restartService(serviceID: $serviceID, environmentID: $environmentID)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"serviceID":     ObjectID(id),
		"environmentID": ObjectID(environmentID),
	})

	return err
}
func (c *client) RedeployService(ctx context.Context, id string, environmentID string) error {
	var mutation struct {
		RedeployService bool `graphql:"redeployService(serviceID: $serviceID, environmentID: $environmentID)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"serviceID":     ObjectID(id),
		"environmentID": ObjectID(environmentID),
	})

	return err

}
func (c *client) SuspendService(ctx context.Context, id string, environmentID string) error {
	var mutation struct {
		SuspendService bool `graphql:"suspendService(serviceID: $serviceID, environmentID: $environmentID)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"serviceID":     ObjectID(id),
		"environmentID": ObjectID(environmentID),
	})

	return err
}

func (c *client) ExposeService(ctx context.Context, id string, environmentID string, projectID string, name string) (*model.TempTCPPort, error) {
	var mutation struct {
		ExposeService model.TempTCPPort `graphql:"exposeTempTcpPort(serviceID: $serviceID, environmentID: $environmentID, projectID: $projectID, serviceName: $name)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"serviceID":     ObjectID(id),
		"environmentID": ObjectID(environmentID),
		"projectID":     ObjectID(projectID),
		"name":          name,
	})

	if err != nil {
		return nil, err
	}

	return &mutation.ExposeService, nil
}

func (c *client) GetMarketplaceItems(ctx context.Context) ([]model.MarketplaceItem, error) {
	var query struct {
		MarketplaceItems []model.MarketplaceItem `graphql:"marketplaceItems"`
	}

	err := c.Query(ctx, &query, nil)

	if err != nil {
		return nil, err
	}

	return query.MarketplaceItems, nil
}

func (c *client) CreateServiceFromMarketplace(ctx context.Context, projectID string, name string, itemCode string) (*model.Service, error) {
	var mutation struct {
		CreateServiceFromMarketplace model.Service `graphql:"createServiceFromMarketplace(projectID: $projectID, name: $name, itemCode: $itemCode)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"projectID": ObjectID(projectID),
		"name":      name,
		"itemCode":  itemCode,
	})

	if err != nil {
		return nil, err
	}

	return &mutation.CreateServiceFromMarketplace, nil
}

func (c *client) SearchGitRepositories(ctx context.Context, keyword *string) ([]model.GitRepo, error) {
	var query struct {
		SearchGitRepositories []model.GitRepo `graphql:"searchGitRepositories(limit: 5, provider: GITHUB, keyword: $keyword)"`
	}

	err := c.Query(ctx, &query, V{
		"keyword": keyword,
	})

	if err != nil {
		return nil, err
	}

	return query.SearchGitRepositories, nil
}

func (c *client) CreateService(ctx context.Context, projectID string, name string, repoID int, branchName string) (*model.Service, error) {
	var mutation struct {
		CreateService model.Service `graphql:"createService(projectID: $projectID, template: $template, name: $name, gitProvider: $gitProvider repoID: $repoID, branchName: $branchName)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"projectID": ObjectID(projectID),
		// specify template as "ServiceTemplate" type
		"template":    ServiceTemplate("GIT"),
		"gitProvider": GitProvider("GITHUB"),
		"name":        name,
		"repoID":      repoID,
		"branchName":  branchName,
	})

	if err != nil {
		return nil, err
	}

	return &mutation.CreateService, nil
}
