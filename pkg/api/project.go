package api

import (
	"context"
	"fmt"

	"github.com/zeabur/cli/pkg/model"
)

// ListProjects returns projects owned by the given owner. An empty ownerID
// means the caller's personal projects (the pre-workspace default — preserved
// for callers that don't yet pass an owner).
//
// Note: the backend hasn't implemented pagination yet, currently we return all
// projects at once.
func (c *client) ListProjects(ctx context.Context, ownerID string, skip, limit int) (*model.ProjectConnection, error) {
	skip, limit = normalizePagination(skip, limit)

	if ownerID == "" {
		var query struct {
			Projects model.ProjectConnection `graphql:"projects(skip: $skip, limit: $limit)"`
		}
		if err := c.Query(ctx, &query, V{
			"skip":  skip,
			"limit": limit,
		}); err != nil {
			return nil, err
		}
		return &query.Projects, nil
	}

	var query struct {
		Projects model.ProjectConnection `graphql:"projects(ownerID: $ownerID, skip: $skip, limit: $limit)"`
	}
	if err := c.Query(ctx, &query, V{
		"ownerID": ObjectID(ownerID),
		"skip":    skip,
		"limit":   limit,
	}); err != nil {
		return nil, err
	}
	return &query.Projects, nil
}

// ListAllProjects walks every page of ListProjects for the given owner.
func (c *client) ListAllProjects(ctx context.Context, ownerID string) (model.Projects, error) {
	skip := 0
	next := true

	var projects []*model.Project

	for next {
		// Propagate the caller's context so cancellation / deadlines
		// reach each page request (CodeRabbit PLA-1590 review).
		projectCon, err := c.ListProjects(ctx, ownerID, skip, 100)
		if err != nil {
			return nil, err
		}
		for _, project := range projectCon.Edges {
			projects = append(projects, project.Node)
		}

		skip += 100
		next = projectCon.PageInfo.HasNextPage
	}

	return projects, nil
}

// GetProject returns a project by (its ID), or (owner name and project name).
func (c *client) GetProject(ctx context.Context, id string, ownerUsername string, projectName string) (*model.Project, error) {
	if id == "" {
		return c.getProjectByOwnerUsernameAndProject(ctx, ownerUsername, projectName)
	}

	return c.getProjectByID(ctx, id)
}

func (c *client) ExportProject(ctx context.Context, id string, environmentID string) (*model.ExportedTemplate, error) {
	var query struct {
		Project struct {
			ExportedTemplate model.ExportedTemplate `graphql:"exportedTemplate(environmentID: $environmentID)"`
		} `graphql:"project(_id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id":            ObjectID(id),
		"environmentID": environmentID,
	})
	if err != nil {
		return nil, err
	}

	return &query.Project.ExportedTemplate, nil
}

func (c *client) getProjectByID(ctx context.Context, id string) (*model.Project, error) {
	var query struct {
		Project model.Project `graphql:"project(_id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id": ObjectID(id),
	})
	if err != nil {
		return nil, err
	}

	return &query.Project, nil
}

func (c *client) getProjectByOwnerUsernameAndProject(ctx context.Context,
	ownerUsername string, projectName string,
) (*model.Project, error) {
	var query struct {
		Project model.Project `graphql:"project(owner: $owner, name: $name)"`
	}

	err := c.Query(ctx, &query, V{
		"owner": ownerUsername,
		"name":  projectName,
	})
	if err != nil {
		return nil, err
	}

	return &query.Project, nil
}

// CreateProject creates a project under the given owner. An empty ownerID
// creates the project under the caller's personal account (the pre-workspace
// default).
func (c *client) CreateProject(ctx context.Context, ownerID, region string, name *string) (*model.Project, error) {
	if ownerID == "" {
		var mutation struct {
			CreateProject model.Project `graphql:"createProject(region: $region, name: $name)"`
		}
		if err := c.Mutate(ctx, &mutation, V{
			"region": region,
			"name":   name,
		}); err != nil {
			return nil, err
		}
		return &mutation.CreateProject, nil
	}

	var mutation struct {
		CreateProject model.Project `graphql:"createProject(ownerID: $ownerID, region: $region, name: $name)"`
	}
	if err := c.Mutate(ctx, &mutation, V{
		"ownerID": ObjectID(ownerID),
		"region":  region,
		"name":    name,
	}); err != nil {
		return nil, err
	}
	return &mutation.CreateProject, nil
}

// Delete a project with the given id
func (c *client) DeleteProject(ctx context.Context, id string) error {
	var mutation struct {
		DeleteProject bool `graphql:"deleteProject(_id: $id)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id": ObjectID(id),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) GetRegions(ctx context.Context) ([]model.Region, error) {
	var query struct {
		Regions []model.Region `graphql:"regions"`
	}

	err := c.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return query.Regions, nil
}

func (c *client) GetServers(ctx context.Context) ([]model.Server, error) {
	var query struct {
		Servers []model.Server `graphql:"servers"`
	}

	err := c.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return query.Servers, nil
}

func (c *client) GetGenericRegions(ctx context.Context) ([]model.GenericRegion, error) {
	regions, err := c.GetRegions(ctx)
	if err != nil {
		return nil, fmt.Errorf("get regions: %w", err)
	}
	servers, err := c.GetServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get servers: %w", err)
	}

	genericRegions := make([]model.GenericRegion, 0, len(regions)+len(servers))
	for _, region := range regions {
		genericRegions = append(genericRegions, region)
	}
	for _, server := range servers {
		genericRegions = append(genericRegions, server)
	}

	return genericRegions, nil
}

// CloneProject clones a project to a target region.
func (c *client) CloneProject(ctx context.Context, projectID, environmentID, targetRegion string, suspendOldProject bool) (*model.CloneProjectResult, error) {
	var mutation struct {
		CloneProject model.CloneProjectResult `graphql:"cloneProject(projectId: $projectId, environmentId: $environmentId, targetRegion: $targetRegion, suspendOldProject: $suspendOldProject, preserveGroupsOrder: true)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"projectId":         ObjectID(projectID),
		"environmentId":     ObjectID(environmentID),
		"targetRegion":      targetRegion,
		"suspendOldProject": suspendOldProject,
	})
	if err != nil {
		return nil, err
	}

	return &mutation.CloneProject, nil
}

// CloneProjectStatus queries the status of a project clone operation.
func (c *client) CloneProjectStatus(ctx context.Context, newProjectID string) (*model.CloneProjectStatusResult, error) {
	var query struct {
		CloneProjectStatus model.CloneProjectStatusResult `graphql:"cloneProjectStatus(newProjectId: $newProjectId)"`
	}

	err := c.Query(ctx, &query, V{
		"newProjectId": ObjectID(newProjectID),
	})
	if err != nil {
		return nil, err
	}

	return &query.CloneProjectStatus, nil
}
