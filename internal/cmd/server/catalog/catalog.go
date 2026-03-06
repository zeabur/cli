package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	provider  string
	country   string
	minCPU    int
	minMemory int
	gpu       bool
}

type catalogOutput struct {
	Providers []providerOutput `json:"providers"`
}

type providerOutput struct {
	Code    string         `json:"code"`
	Name    string         `json:"name"`
	Regions []regionOutput `json:"regions"`
}

type regionOutput struct {
	ID        string                      `json:"id"`
	Name      string                      `json:"name"`
	City      string                      `json:"city"`
	Country   string                      `json:"country"`
	Continent string                      `json:"continent"`
	Plans     []model.DedicatedServerPlan `json:"plans"`
}

func NewCmdCatalog(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "List all available providers, regions, and plans for renting a server",
		Long:  "List all available providers, regions, and plans. Supports filters to narrow results.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCatalog(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.provider, "provider", "", "Filter by provider code (e.g. hetzner, vultr)")
	cmd.Flags().StringVar(&opts.country, "country", "", "Filter by country code (e.g. US, DE, JP)")
	cmd.Flags().IntVar(&opts.minCPU, "min-cpu", 0, "Filter plans with at least this many CPU cores")
	cmd.Flags().IntVar(&opts.minMemory, "min-memory", 0, "Filter plans with at least this much memory (MB)")
	cmd.Flags().BoolVar(&opts.gpu, "gpu", false, "Only show plans with GPU")

	return cmd
}

func runCatalog(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()

	providers, err := f.ApiClient.ListDedicatedServerProviders(ctx)
	if err != nil {
		return fmt.Errorf("list providers failed: %w", err)
	}

	// filter providers
	if opts.provider != "" {
		var filtered []model.CloudProvider
		for _, p := range providers {
			if strings.EqualFold(p.Code, opts.provider) {
				filtered = append(filtered, p)
			}
		}
		providers = filtered
	}

	output := catalogOutput{
		Providers: make([]providerOutput, len(providers)),
	}

	if len(providers) == 0 {
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for i, p := range providers {
		output.Providers[i] = providerOutput{Code: p.Code, Name: p.Name}

		wg.Add(1)
		go func(idx int, providerCode string) {
			defer wg.Done()

			regions, err := f.ApiClient.ListDedicatedServerRegions(ctx, providerCode)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("list regions for %s failed: %w", providerCode, err)
				}
				mu.Unlock()
				return
			}

			// filter regions by country
			if opts.country != "" {
				var filtered []model.DedicatedServerRegion
				for _, r := range regions {
					if strings.EqualFold(r.Country, opts.country) {
						filtered = append(filtered, r)
					}
				}
				regions = filtered
			}

			regionOutputs := make([]regionOutput, len(regions))
			var regionWg sync.WaitGroup

			for j, r := range regions {
				regionWg.Add(1)
				go func(rIdx int, region model.DedicatedServerRegion) {
					defer regionWg.Done()

					plans, err := f.ApiClient.ListDedicatedServerPlans(ctx, providerCode, region.ID)
					if err != nil {
						mu.Lock()
						if firstErr == nil {
							firstErr = fmt.Errorf("list plans for %s/%s failed: %w", providerCode, region.ID, err)
						}
						mu.Unlock()
						return
					}

					// filter plans
					var filtered []model.DedicatedServerPlan
					for _, plan := range plans {
						if opts.minCPU > 0 && plan.CPU < opts.minCPU {
							continue
						}
						if opts.minMemory > 0 && plan.Memory < opts.minMemory {
							continue
						}
						if opts.gpu && plan.GPU == nil {
							continue
						}
						filtered = append(filtered, plan)
					}

					regionOutputs[rIdx] = regionOutput{
						ID:        region.ID,
						Name:      region.Name,
						City:      region.City,
						Country:   region.Country,
						Continent: region.Continent,
						Plans:     filtered,
					}
				}(j, r)
			}

			regionWg.Wait()

			// remove regions with no matching plans
			var nonEmpty []regionOutput
			for _, r := range regionOutputs {
				if len(r.Plans) > 0 {
					nonEmpty = append(nonEmpty, r)
				}
			}

			mu.Lock()
			output.Providers[idx].Regions = nonEmpty
			mu.Unlock()
		}(i, p.Code)
	}

	wg.Wait()

	if firstErr != nil {
		return firstErr
	}

	// remove providers with no matching regions
	nonEmpty := make([]providerOutput, 0, len(output.Providers))
	for _, p := range output.Providers {
		if len(p.Regions) > 0 {
			nonEmpty = append(nonEmpty, p)
		}
	}
	output.Providers = nonEmpty

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal output failed: %w", err)
	}

	fmt.Println(string(data))

	return nil
}
