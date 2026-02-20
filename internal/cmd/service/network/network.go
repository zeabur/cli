package network

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
)

type Options struct {
	id            string
	name          string
	environmentID string
}

func NewCmdPrivateNetwork(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "network",
		Short:   "Network information for service",
		Long:    `Network information for service`,
		Aliases: []string{"net"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNetwork(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)

	return cmd
}

func runNetwork(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive && opts.id == "" && opts.name == "" {
		zctx := f.Config.GetContext()
		if _, err := f.ParamFiller.ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions{
			ProjectCtx:    zctx,
			ServiceID:     &opts.id,
			ServiceName:   &opts.name,
			EnvironmentID: &opts.environmentID,
			CreateNew:     false,
		}); err != nil {
			return err
		}
	}

	if opts.id == "" && opts.name != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = service.ID
	}

	if opts.id == "" {
		return fmt.Errorf("service id or name is required")
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Fetching network information of service %s ...", opts.name)),
	)
	s.Start()

	dnsName, err := f.ApiClient.GetDNSName(context.Background(), opts.id)
	if err != nil {
		s.Stop()
		return err
	}

	s.Stop()

	f.Log.Infof("Private DNS name for %s: %s", opts.name, dnsName+".zeabur.internal")

	return nil
}
