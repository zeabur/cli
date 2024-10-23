package api

import (
	"context"
	"fmt"

	"github.com/zeabur/cli/pkg/model"
)

// ListProjects returns projects owned by the current user.
// Note: the backend hasn't implemented pagination yet, currently we return all projects at once.
func (c *client) ListProjects(ctx context.Context, skip, limit int) (*model.ProjectConnection, error) {
	skip, limit = normalizePagination(skip, limit)

	var query struct {
		Projects model.ProjectConnection `graphql:"projects(skip: $skip, limit: $limit)"`
	}

	err := c.Query(ctx, &query, V{
		"skip":  skip,
		"limit": limit,
	})
	if err != nil {
		return nil, err
	}

	return &query.Projects, nil
}

// ListAllProjects returns all projects owned by the current user.
func (c *client) ListAllProjects(ctx context.Context) (model.Projects, error) {
	skip := 0
	next := true

	var projects []*model.Project

	for next {
		projectCon, err := c.ListProjects(context.Background(), skip, 100)
		if err != nil {
			return nil, err
		}
		for _, project := range projectCon.Edges {
			projects = append(projects, project.Node)
		}

		skip += 5
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
	ownerUsername string, projectName string) (*model.Project, error) {
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

// Create a project with the region and optional name.
func (c *client) CreateProject(ctx context.Context, region string, name *string) (*model.Project, error) {
	var mutation struct {
		CreateProject model.Project `graphql:"createProject(region: $region, name: $name)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"region": region,
		"name":   name,
	})
	if err != nil {
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
