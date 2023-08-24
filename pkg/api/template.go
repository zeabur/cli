package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) ListTemplates(ctx context.Context, skip, limit int) (*model.TemplateConnection, error) {
	skip, limit = normalizePagination(skip, limit)

	var query struct {
		Templates model.TemplateConnection `graphql:"templates(skip: $skip, limit: $limit)"`
	}

	err := c.Query(ctx, &query, V{
		"skip":  skip,
		"limit": limit,
	})
	if err != nil {
		return nil, err
	}

	return &query.Templates, nil
}

func (c *client) ListAllTemplates(ctx context.Context) (model.Templates, error) {
	skip := 0
	next := true

	var templates []*model.Template

	for next {
		templateCon, err := c.ListTemplates(context.Background(), skip, 100)
		if err != nil {
			return nil, err
		}
		for _, template := range templateCon.Edges {
			templates = append(templates, template.Node)
		}

		skip += 5
		next = templateCon.PageInfo.HasNextPage
	}

	return templates, nil
}
