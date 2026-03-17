package network

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/model"
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

	ctx := context.Background()

	dnsName, err := f.ApiClient.GetDNSName(ctx, opts.id)
	if err != nil {
		s.Stop()
		return err
	}

	// Fetch port forwarding info if environmentID is available
	var portForwardingMode model.PortForwardingMode
	var ports []model.ServicePort
	var portForwardedHost string

	if opts.environmentID != "" {
		portForwardingMode, err = f.ApiClient.GetPortForwardingMode(ctx, opts.id, opts.environmentID)
		if err != nil {
			s.Stop()
			return err
		}

		if portForwardingMode == model.PortForwardingModeEnabled {
			ports, err = f.ApiClient.GetServicePorts(ctx, opts.id, opts.environmentID)
			if err != nil {
				s.Stop()
				return err
			}

			portForwardedHost, err = f.ApiClient.GetPortForwardedHost(ctx, opts.id)
			if err != nil {
				s.Stop()
				return err
			}
		}
	}

	s.Stop()

	if f.JSON {
		result := map[string]interface{}{
			"dnsName": dnsName + ".zeabur.internal",
		}
		if opts.environmentID != "" {
			result["portForwardingMode"] = string(portForwardingMode)
			if portForwardingMode == model.PortForwardingModeEnabled {
				result["portForwardedHost"] = portForwardedHost
				portList := make([]map[string]interface{}, 0, len(ports))
				for _, p := range ports {
					pm := map[string]interface{}{
						"id":   p.ID,
						"port": p.Port,
						"type": p.Type,
					}
					if p.ForwardedPort != nil {
						pm["forwardedPort"] = *p.ForwardedPort
					}
					portList = append(portList, pm)
				}
				result["ports"] = portList
			}
		}
		return f.Printer.JSON(result)
	}

	f.Log.Infof("Private DNS name for %s: %s", opts.name, dnsName+".zeabur.internal")

	if opts.environmentID != "" {
		f.Log.Infof("Port forwarding: %s", portForwardingMode)
		if portForwardingMode == model.PortForwardingModeEnabled && portForwardedHost != "" {
			for _, p := range ports {
				if p.ForwardedPort != nil && (p.Type == "TCP" || p.Type == "UDP") {
					f.Log.Infof("  %s (%s %d) → %s:%d", p.ID, p.Type, p.Port, portForwardedHost, *p.ForwardedPort)
				}
			}
		}
	}

	return nil
}
