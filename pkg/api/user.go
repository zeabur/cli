package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) GetUserInfo(ctx context.Context) (*model.User, error) {
	var query struct {
		User model.User `graphql:"me"`
	}

	err := c.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return &query.User, nil
}
