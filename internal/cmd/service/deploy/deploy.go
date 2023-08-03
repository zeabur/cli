package deploy

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	projectID string
	template  string
	itemCode  string
	name      string
}

func NewCmdDeploy(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   "Deploy a service",
		PreRunE: util.NeedProjectContext(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeploy(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.template, "template", "", "Service template")
	cmd.Flags().StringVar(&opts.itemCode, "item-code", "", "Marketplace item code")

	return cmd
}

func runDeploy(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runDeployInteractive(f, opts)
	} else {
		return runDeployNonInteractive(f, opts)
	}
}

func runDeployNonInteractive(f *cmdutil.Factory, opts *Options) error {
	err := paramCheck(opts)
	if err != nil {
		return err
	}

	ctx := context.Background()

	opts.projectID = f.Config.GetContext().GetProject().GetID()

	if opts.template == "MARKETPLACE" {
		opts.name = opts.itemCode
		service, err := f.ApiClient.CreateServiceFromMarketplace(ctx, opts.projectID, opts.name, opts.itemCode)
		if err != nil {
			return fmt.Errorf("create service failed: %w", err)
		}

		f.Log.Infof("Service %s created", service.Name)
	}

	return nil
}

func runDeployInteractive(f *cmdutil.Factory, opts *Options) error {
	serviceTemplate, err := f.Prompter.Select("Select service template", "MARKETPLACE", []string{"MARKETPLACE", "GIT"})
	if err != nil {
		return err
	}

	ctx := context.Background()

	opts.projectID = f.Config.GetContext().GetProject().GetID()

	if serviceTemplate == 0 {
		s := spinner.New(spinner.CharSets[1], 100*time.Millisecond)
		s.Start()
		marketplaceItems, err := f.ApiClient.GetMarketplaceItems(ctx)
		if err != nil {
			return fmt.Errorf("get marketplace items failed: %w", err)
		}
		s.Stop()

		marketplaceItemsList := make([]string, len(marketplaceItems))
		for i, item := range marketplaceItems {
			marketplaceItemsList[i] = item.Name + " (" + item.Description + ")"
		}

		index, err := f.Prompter.Select("Select marketplace item", "", marketplaceItemsList)
		if err != nil {
			return fmt.Errorf("select marketplace item failed: %w", err)
		}

		opts.itemCode = marketplaceItems[index].Code
		opts.name = opts.itemCode

		s.Start()
		service, err := f.ApiClient.CreateServiceFromMarketplace(ctx, opts.projectID, opts.name, opts.itemCode)
		if err != nil {
			return fmt.Errorf("create service failed: %w", err)
		}
		s.Stop()

		f.Log.Infof("Service %s created", service.Name)
	}

	return nil
}

func paramCheck(opts *Options) error {
	if opts.template == "" {
		return fmt.Errorf("please specify service template with --template")
	}

	if opts.template == "MARKETPLACE" && opts.itemCode == "" {
		return fmt.Errorf("please specify marketplace item code with --item-code")
	}

	return nil
}
