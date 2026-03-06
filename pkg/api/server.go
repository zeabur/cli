package api

import (
	"context"

	"github.com/zeabur/cli/pkg/model"
)

func (c *client) ListServers(ctx context.Context) (model.ServerListItems, error) {
	var query struct {
		Servers model.ServerListItems `graphql:"servers"`
	}

	err := c.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return query.Servers, nil
}

func (c *client) GetServer(ctx context.Context, id string) (*model.ServerDetail, error) {
	var query struct {
		Server model.ServerDetail `graphql:"server(_id: $id)"`
	}

	err := c.Query(ctx, &query, V{
		"id": ObjectID(id),
	})
	if err != nil {
		return nil, err
	}

	return &query.Server, nil
}

func (c *client) RebootServer(ctx context.Context, id string) error {
	var mutation struct {
		RebootServer bool `graphql:"rebootServer(_id: $id)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"id": ObjectID(id),
	})

	return err
}

func (c *client) ListDedicatedServerProviders(ctx context.Context) ([]model.CloudProvider, error) {
	var query struct {
		Providers []model.CloudProvider `graphql:"dedicatedServerProviders"`
	}

	err := c.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return query.Providers, nil
}

func (c *client) ListDedicatedServerRegions(ctx context.Context, provider string) ([]model.DedicatedServerRegion, error) {
	var query struct {
		Regions []model.DedicatedServerRegion `graphql:"dedicatedServerRegions(provider: $provider)"`
	}

	err := c.Query(ctx, &query, V{
		"provider": provider,
	})
	if err != nil {
		return nil, err
	}

	return query.Regions, nil
}

func (c *client) ListDedicatedServerPlans(ctx context.Context, provider, region string) (model.DedicatedServerPlans, error) {
	var query struct {
		Plans model.DedicatedServerPlans `graphql:"dedicatedServerPlans(provider: $provider, region: $region)"`
	}

	err := c.Query(ctx, &query, V{
		"provider": provider,
		"region":   region,
	})
	if err != nil {
		return nil, err
	}

	return query.Plans, nil
}

func (c *client) RentServer(ctx context.Context, provider, region, plan string) (string, error) {
	var mutation struct {
		RentServer struct {
			ID string `graphql:"_id"`
		} `graphql:"rentServer(provider: $provider, region: $region, plan: $plan)"`
	}

	err := c.Mutate(ctx, &mutation, V{
		"provider": provider,
		"region":   region,
		"plan":     plan,
	})
	if err != nil {
		return "", err
	}

	return mutation.RentServer.ID, nil
}
