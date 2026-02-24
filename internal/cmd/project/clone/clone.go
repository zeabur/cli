package clone

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	ProjectID     string
	ProjectName   string
	EnvironmentID string
	Region        string
	Suspend       bool
}

func NewCmdClone(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone a project to another region",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClone(f, opts)
		},
	}

	util.AddProjectParam(cmd, &opts.ProjectID, &opts.ProjectName)
	cmd.Flags().StringVar(&opts.EnvironmentID, "env-id", "", "Source environment ID (auto-resolved if omitted)")
	cmd.Flags().StringVarP(&opts.Region, "region", "r", "", "Target region")
	cmd.Flags().BoolVar(&opts.Suspend, "suspend", false, "Suspend old project after cloning")

	return cmd
}

func runClone(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err == nil {
		return runCloneNonInteractive(f, opts)
	}

	if f.Interactive {
		return runCloneInteractive(f, opts)
	}

	return runCloneNonInteractive(f, opts)
}

func paramCheck(opts *Options) error {
	if opts.ProjectID == "" && opts.ProjectName == "" {
		return fmt.Errorf("please specify project with --id or --name")
	}
	if opts.Region == "" {
		return fmt.Errorf("please specify target region with --region")
	}
	return nil
}

func runCloneInteractive(f *cmdutil.Factory, opts *Options) error {
	// Select project if not provided
	if opts.ProjectID == "" && opts.ProjectName == "" {
		projectInfo, _, err := f.Selector.SelectProject()
		if err != nil {
			return fmt.Errorf("select project: %w", err)
		}
		opts.ProjectID = projectInfo.GetID()
	} else if opts.ProjectID == "" && opts.ProjectName != "" {
		project, err := util.GetProjectByName(f.Config, f.ApiClient, opts.ProjectName)
		if err != nil {
			return err
		}
		opts.ProjectID = project.ID
	}

	// Resolve environment ID if not provided
	if opts.EnvironmentID == "" {
		envID, err := util.ResolveEnvironmentID(f.ApiClient, opts.ProjectID)
		if err != nil {
			return err
		}
		opts.EnvironmentID = envID
	}

	// Select region if not provided
	if opts.Region == "" {
		s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Fetching available regions..."),
		)
		s.Start()
		regions, err := f.ApiClient.GetGenericRegions(context.Background())
		if err != nil {
			s.Stop()
			return err
		}
		s.Stop()

		availableRegions := make([]model.GenericRegion, 0, len(regions))
		regionOptions := make([]string, 0, len(regions))
		for _, region := range regions {
			if region.IsAvailable() {
				availableRegions = append(availableRegions, region)
				regionOptions = append(regionOptions, region.String())
			}
		}

		regionIndex, err := f.Prompter.Select("Select target region", "", regionOptions)
		if err != nil {
			return err
		}

		opts.Region = availableRegions[regionIndex].GetID()
	}

	return doClone(f, opts)
}

func runCloneNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}

	// Resolve project ID from name if needed
	if opts.ProjectID == "" && opts.ProjectName != "" {
		project, err := util.GetProjectByName(f.Config, f.ApiClient, opts.ProjectName)
		if err != nil {
			return err
		}
		opts.ProjectID = project.ID
	}

	// Resolve environment ID if not provided
	if opts.EnvironmentID == "" {
		envID, err := util.ResolveEnvironmentID(f.ApiClient, opts.ProjectID)
		if err != nil {
			return err
		}
		opts.EnvironmentID = envID
	}

	return doClone(f, opts)
}

func doClone(f *cmdutil.Factory, opts *Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Cloning project..."),
	)
	s.Start()

	result, err := f.ApiClient.CloneProject(
		context.Background(),
		opts.ProjectID,
		opts.EnvironmentID,
		opts.Region,
		opts.Suspend,
	)
	if err != nil {
		s.Stop()
		return fmt.Errorf("clone project failed: %w", err)
	}

	s.Stop()
	f.Log.Infof("Clone started, new project ID: %s", result.NewProjectID)

	// Poll for status
	seenEvents := 0
	for {
		time.Sleep(3 * time.Second)

		status, err := f.ApiClient.CloneProjectStatus(context.Background(), result.NewProjectID)
		if err != nil {
			return fmt.Errorf("query clone status failed: %w", err)
		}

		// Print new events and check for terminal states
		completed := false
		var failMsg string
		for i := seenEvents; i < len(status.Events); i++ {
			ev := status.Events[i]
			f.Log.Infof("[%s] %s", ev.Type, ev.Message)
			if ev.Type == "CloneProjectCompleted" {
				completed = true
			}
			if ev.Type == "CloneProjectFailed" {
				failMsg = ev.Message
			}
		}
		seenEvents = len(status.Events)

		if failMsg != "" {
			return fmt.Errorf("clone failed: %s", failMsg)
		}

		if status.Error != nil && *status.Error != "" {
			return fmt.Errorf("clone failed: %s", *status.Error)
		}

		if completed {
			f.Log.Infof("Project cloned successfully!")
			f.Log.Infof("New project ID: %s", result.NewProjectID)
			f.Log.Infof("Dashboard: https://dash.zeabur.com/projects/%s", result.NewProjectID)
			return nil
		}
	}
}
