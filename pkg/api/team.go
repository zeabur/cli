package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

// ListTeams returns every team the current caller belongs to, with the
// caller's role in each team. Drives `zeabur workspace list` and the lazy
// verify run at CLI startup.
func (c *client) ListTeams(ctx context.Context) ([]model.Team, error) {
	var query struct {
		Teams []model.Team `graphql:"teams"`
	}

	if err := c.Query(ctx, &query, nil); err != nil {
		return nil, err
	}
	return query.Teams, nil
}
