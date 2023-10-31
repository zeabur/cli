package create

import (
	"context"
	"fmt"
	"github.com/zeabur/cli/pkg/zcontext"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	ProjectRegion string
	ProjectName   string
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(f, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.ProjectName, "name", "n", "", "Project name")
	cmd.Flags().StringVarP(&opts.ProjectRegion, "region", "r", "", "Project region")

	return cmd
}

func runCreate(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err == nil {
		return runCreateNonInteractive(f, opts)
	}

	if f.Interactive {
		return runCreateInteractive(f, opts)
	}

	return runCreateNonInteractive(f, opts)
}

func runCreateInteractive(f *cmdutil.Factory, opts *Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching available project regions..."),
	)
	s.Start()
	regions, err := f.ApiClient.GetRegions(context.Background())
	if err != nil {
		return err
	}
	s.Stop()

	regionIDs := make([]string, 0, len(regions))
	for _, region := range regions {
		regionIDs = append(regionIDs, region.ID)
	}

	projectRegionIndex, err := f.Prompter.Select("Select project region", "", regionIDs)
	if err != nil {
		return err
	}

	projectRegion := regions[projectRegionIndex].ID

	s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Creating project..."),
	)
	s.Start()
	if err := createProject(f, projectRegion, &opts.ProjectName); err != nil {
		return err
	}
	s.Stop()

	return nil
}

func runCreateNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}

	err := createProject(f, opts.ProjectRegion, &opts.ProjectName)
	if err != nil {
		f.Log.Error(err)
		return err
	}

	return nil
}

func createProject(f *cmdutil.Factory, projectRegion string, projectName *string) error {
	project, err := f.ApiClient.CreateProject(context.Background(), projectRegion, projectName)
	if err != nil {
		f.Log.Error(err)
		return err
	}

	f.Log.Infof("Project %s created", project.Name)
	err = setProject(f, project.ID, project.Name)
	if err != nil {
		f.Log.Error(err)
		return err
	}
	return nil
}

func paramCheck(opts *Options) error {
	if opts.ProjectRegion == "" {
		return fmt.Errorf("please specify project region with --region")
	}

	return nil
}

func setProject(f *cmdutil.Factory, id, name string) error {
	if id == "" && name == "" {
		return fmt.Errorf("context auto-switching failed, check if the project was created successfully")
	}
	f.Config.GetContext().SetProject(zcontext.NewBasicInfo(id, name))

	// clear environment and service context when project context is set
	f.Config.GetContext().ClearService()
	f.Config.GetContext().ClearEnvironment()

	f.Log.Infof("Project context is set to <%s>", name)

	return nil
}
