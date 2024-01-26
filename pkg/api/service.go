package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/spf13/viper"
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

func (c *client) ServiceMetric(ctx context.Context, id, projectID, environmentID, metricType string, startTime, endTime time.Time) (*model.ServiceMetric, error) {
	var query struct {
		ServiceMetric model.ServiceMetric `graphql:"service(_id: $serviceID)"`
	}

	err := c.Query(ctx, &query, V{
		"serviceID":     ObjectID(id),
		"projectID":     ObjectID(projectID),
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

func (c *client) ServiceInstructions(ctx context.Context, id, environmentID string) ([]model.ServiceInstruction, error) {
	var query struct {
		ServiceInstructions []model.ServiceInstruction `graphql:"instructions(serviceID: $serviceID, environmentID: $environmentID)"`
	}

	err := c.Query(ctx, &query, V{
		"serviceID":     ObjectID(id),
		"environmentID": ObjectID(environmentID),
	})

	if err != nil {
		return nil, err
	}

	return query.ServiceInstructions, nil
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

func (c *client) GetPrebuiltItems(ctx context.Context) ([]model.PrebuiltItem, error) {
	var query struct {
		PrebuiltItems []model.PrebuiltItem `graphql:"prebuiltMarketplaceItems"`
	}

	err := c.Query(ctx, &query, nil)

	if err != nil {
		return nil, err
	}

	return query.PrebuiltItems, nil
}

func (c *client) CreatePrebuiltService(ctx context.Context, projectID string, marketplaceCode string) (*model.Service, error) {
	var mutation struct {
		CreatePrebuiltService model.Service `graphql:"createGenericService(projectID: $projectID, marketplaceCode: $marketplaceCode)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"projectID":       ObjectID(projectID),
		"marketplaceCode": marketplaceCode,
	})

	if err != nil {
		return nil, err
	}

	return &mutation.CreatePrebuiltService, nil
}

func (c *client) SearchGitRepositories(ctx context.Context, keyword *string) ([]model.GitRepo, error) {
	var query struct {
		SearchGitRepositories []model.GitRepo `graphql:"searchGitRepositories(Limit: 5, provider: GITHUB, keyword: $keyword)"`
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
		"projectID":   ObjectID(projectID),
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

func (c *client) CreateEmptyService(ctx context.Context, projectID string, name string) (*model.Service, error) {
	var mutation struct {
		CreateService model.Service `graphql:"createService(projectID: $projectID, template: $template, name: $name)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"projectID": ObjectID(projectID),
		"template":  ServiceTemplate("GIT"),
		"name":      name,
	})

	if err != nil {
		return nil, err
	}

	return &mutation.CreateService, nil
}

func (c *client) UploadZipToService(ctx context.Context, projectID string, serviceID string, environmentID string, zipBytes []byte) (*model.Service, error) {
	url := "https://gateway.zeabur.com/projects/" + projectID + "/services/" + serviceID + "/deploy"

	method := "POST"

	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	err := multipartWriter.WriteField("environment", environmentID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fileWriter, err := multipartWriter.CreateFormFile("code", "zeabur.zip")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	_, err = io.Copy(fileWriter, bytes.NewReader(zipBytes))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = multipartWriter.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest(method, url, &requestBody)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	token := viper.GetString("token")

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	req.Header.Set("Cookie", "token="+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	return nil, nil
}
