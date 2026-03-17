package portforward

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

	enable  bool
	disable bool
}

func NewCmdPortForward(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "port-forward",
		Short: "Manage port forwarding for a service",
		Long: `Manage port forwarding for a service.
example:
      zeabur service port-forward                              # show status (interactive)
      zeabur service port-forward --enable                     # enable
      zeabur service port-forward --disable                    # disable
      zeabur service port-forward --id SERVICE_ID --enable     # non-interactive
`,
		Aliases: []string{"pf"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPortForward(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().BoolVar(&opts.enable, "enable", false, "Enable port forwarding")
	cmd.Flags().BoolVar(&opts.disable, "disable", false, "Disable port forwarding")

	return cmd
}

func runPortForward(f *cmdutil.Factory, opts *Options) error {
	if opts.enable && opts.disable {
		return fmt.Errorf("cannot use both --enable and --disable")
	}

	if f.Interactive {
		return runPortForwardInteractive(f, opts)
	}
	return runPortForwardNonInteractive(f, opts)
}

func runPortForwardInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" && opts.name == "" {
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

	return runPortForwardNonInteractive(f, opts)
}

func runPortForwardNonInteractive(f *cmdutil.Factory, opts *Options) error {
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

	if opts.environmentID == "" {
		return fmt.Errorf("environment id is required (use --env-id or select interactively)")
	}

	ctx := context.Background()

	// If --enable or --disable, update the mode
	if opts.enable || opts.disable {
		mode := model.PortForwardingModeEnabled
		if opts.disable {
			mode = model.PortForwardingModeDisabled
		}

		s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(fmt.Sprintf(" Updating port forwarding mode to %s ...", mode)),
		)
		s.Start()

		err := f.ApiClient.UpdatePortForwardingMode(ctx, opts.id, opts.environmentID, mode)
		s.Stop()
		if err != nil {
			return fmt.Errorf("failed to update port forwarding mode: %w", err)
		}

		f.Log.Infof("Port forwarding mode updated to %s", mode)
		return nil
	}

	// Otherwise, show current status
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching port forwarding status ..."),
	)
	s.Start()

	mode, err := f.ApiClient.GetPortForwardingMode(ctx, opts.id, opts.environmentID)
	if err != nil {
		s.Stop()
		return err
	}

	var ports []model.ServicePort
	var host string

	if mode == model.PortForwardingModeEnabled {
		ports, err = f.ApiClient.GetServicePorts(ctx, opts.id, opts.environmentID)
		if err != nil {
			s.Stop()
			return err
		}

		host, err = f.ApiClient.GetPortForwardedHost(ctx, opts.id)
		if err != nil {
			s.Stop()
			return err
		}
	}

	s.Stop()

	if f.JSON {
		result := map[string]interface{}{
			"portForwardingMode": string(mode),
		}
		if mode == model.PortForwardingModeEnabled {
			result["portForwardedHost"] = host
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
		return f.Printer.JSON(result)
	}

	f.Log.Infof("Port forwarding: %s", mode)
	if mode == model.PortForwardingModeEnabled && host != "" {
		for _, p := range ports {
			if p.ForwardedPort != nil && (p.Type == "TCP" || p.Type == "UDP") {
				f.Log.Infof("  %s (%s %d) → %s:%d", p.ID, p.Type, p.Port, host, *p.ForwardedPort)
			}
		}
	}

	return nil
}
