package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/zeabur/cli/pkg/constant"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/util"
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
	skip, limit int,
) (*model.Connection[model.ServiceDetail], error) {
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
		CreatePrebuiltService model.Service `graphql:"createPrebuiltService(projectID: $projectID, marketplaceCode: $marketplaceCode)"`
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

func (c *client) CreatePrebuiltServiceRaw(ctx context.Context, projectID string, rawSchema string) (*model.Service, error) {
	var mutation struct {
		CreateCustomService model.Service `graphql:"createPrebuiltService(projectID: $projectID, rawSchema: $rawSchema)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"projectID": ObjectID(projectID),
		"schema":    rawSchema,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.CreateCustomService, nil
}

func (c *client) CreatePrebuiltServiceCustom(ctx context.Context, projectID string, schema model.ServiceSpecSchemaInput) (*model.Service, error) {
	var mutation struct {
		CreateCustomService model.Service `graphql:"createPrebuiltService(projectID: $projectID, schema: $schema)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"projectID": ObjectID(projectID),
		"schema":    schema,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.CreateCustomService, nil
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
		CreateService model.Service `graphql:"createService(projectID: $projectID, template: $template, name: $name, gitProvider: $gitProvider, repoID: $repoID, branchName: $branchName)"`
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
	// Step 1: Calculate SHA256 hash of content
	h := sha256.New()
	if _, err := h.Write(zipBytes); err != nil {
		return nil, fmt.Errorf("failed to calculate content hash: %w", err)
	}
	contentHash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Step 2: Create upload session
	createUploadReq := struct {
		ContentHash          string `json:"content_hash"`
		ContentHashAlgorithm string `json:"content_hash_algorithm"`
		ContentLength        int64  `json:"content_length"`
	}{
		ContentHash:          contentHash,
		ContentHashAlgorithm: "sha256",
		ContentLength:        int64(len(zipBytes)),
	}

	createUploadBody, err := json.Marshal(createUploadReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create upload request: %w", err)
	}

	createUploadResp, err := http.NewRequestWithContext(ctx, "POST", constant.ZeaburServerURL+"/v2/upload", bytes.NewReader(createUploadBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create upload request: %w", err)
	}

	token := viper.GetString("token")
	createUploadResp.Header.Set("Content-Type", "application/json")
	createUploadResp.Header.Set("Cookie", "token="+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(createUploadResp)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, util.FormatHTTPError("failed to create upload session", resp)
	}

	var uploadSession struct {
		PresignHeader struct {
			ContentType string `json:"Content-Type"`
		} `json:"presign_header"`
		PresignMethod string `json:"presign_method"`
		PresignURL    string `json:"presign_url"`
		UploadID      string `json:"upload_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&uploadSession); err != nil {
		return nil, fmt.Errorf("failed to decode upload session response: %w", err)
	}

	// Step 3: Upload file to S3
	uploadReq, err := http.NewRequestWithContext(ctx, uploadSession.PresignMethod, uploadSession.PresignURL, bytes.NewReader(zipBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 upload request: %w", err)
	}

	uploadReq.Header.Set("Content-Type", uploadSession.PresignHeader.ContentType)
	uploadReq.Header.Set("Content-Length", strconv.FormatInt(int64(len(zipBytes)), 10))

	uploadResp, err := client.Do(uploadReq)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusOK {
		return nil, util.FormatHTTPError("failed to upload to S3", uploadResp)
	}

	// Step 4: Prepare upload for deployment
	prepareReq := struct {
		UploadType    string `json:"upload_type"`
		ServiceID     string `json:"service_id"`
		EnvironmentID string `json:"environment_id"`
	}{
		UploadType:    "existing_service",
		ServiceID:     serviceID,
		EnvironmentID: environmentID,
	}

	prepareBody, err := json.Marshal(prepareReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal prepare request: %w", err)
	}

	prepareResp, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/v2/upload/%s/prepare", constant.ZeaburServerURL, uploadSession.UploadID),
		bytes.NewReader(prepareBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create prepare request: %w", err)
	}

	prepareResp.Header.Set("Content-Type", "application/json")
	prepareResp.Header.Set("Cookie", "token="+token)

	resp, err = client.Do(prepareResp)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, util.FormatHTTPError("failed to prepare upload", resp)
	}

	return nil, nil
}

func (c *client) GetDNSName(ctx context.Context, serviceID string) (string, error) {
	var query struct {
		Service struct {
			DnsName string `graphql:"dnsName"`
		} `graphql:"service(_id: $serviceID)"`
	}

	err := c.Query(ctx, &query, V{
		"serviceID": ObjectID(serviceID),
	})
	if err != nil {
		return "", err
	}

	return query.Service.DnsName, nil
}

func (c *client) GetPortForwardingMode(ctx context.Context, serviceID string, environmentID string) (model.PortForwardingMode, error) {
	var query struct {
		Service struct {
			PortForwardingMode model.PortForwardingMode `graphql:"portForwardingMode(environmentID: $environmentID)"`
		} `graphql:"service(_id: $serviceID)"`
	}
	err := c.Query(ctx, &query, V{
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
	})
	if err != nil {
		return "", err
	}
	return query.Service.PortForwardingMode, nil
}

func (c *client) UpdatePortForwardingMode(ctx context.Context, serviceID string, environmentID string, mode model.PortForwardingMode) error {
	var mutation struct {
		UpdatePortForwardingMode bool `graphql:"updatePortForwardingMode(serviceID: $serviceID, environmentID: $environmentID, mode: $mode)"`
	}
	return c.Mutate(ctx, &mutation, V{
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
		"mode":          mode,
	})
}

func (c *client) GetServicePorts(ctx context.Context, serviceID string, environmentID string) ([]model.ServicePort, error) {
	var query struct {
		Service struct {
			Ports []model.ServicePort `graphql:"ports(environmentID: $environmentID)"`
		} `graphql:"service(_id: $serviceID)"`
	}
	err := c.Query(ctx, &query, V{
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
	})
	if err != nil {
		return nil, err
	}
	return query.Service.Ports, nil
}

func (c *client) GetPortForwardedHost(ctx context.Context, serviceID string) (string, error) {
	var query struct {
		Service struct {
			PortForwardedHost string `graphql:"portForwardedHost"`
		} `graphql:"service(_id: $serviceID)"`
	}
	err := c.Query(ctx, &query, V{
		"serviceID": ObjectID(serviceID),
	})
	if err != nil {
		return "", err
	}
	return query.Service.PortForwardedHost, nil
}

func (c *client) UpdateImageTag(ctx context.Context, serviceID, environmentID, tag string) error {
	var mutation struct {
		UpdateImageTag bool `graphql:"updateServiceImage(serviceID: $serviceID, environmentID: $environmentID, tag: $tag)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
		"tag":           tag,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) DeleteService(ctx context.Context, id string) error {
	var mutation struct {
		DeleteService bool `graphql:"deleteService(_id: $id)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id": ObjectID(id),
	})

	return err
}

func (c *client) ExecuteCommand(ctx context.Context, serviceID string, environmentID string, command []string) (*model.CommandResult, error) {
	var mutation struct {
		ExecuteCommand model.CommandResult `graphql:"executeCommand(serviceID: $serviceID, environmentID: $environmentID, command: $command)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"serviceID":     ObjectID(serviceID),
		"environmentID": ObjectID(environmentID),
		"command":       command,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.ExecuteCommand, nil
}
