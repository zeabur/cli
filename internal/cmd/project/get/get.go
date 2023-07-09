package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	ID string

	OwnerName   string
	ProjectName string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get project",
		Long:  "Get project, use --id or both --owner and --project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.ID, "id", "", "Project ID")

	cmd.Flags().StringVarP(&opts.OwnerName, "owner", "o", "", "Owner name")
	cmd.Flags().StringVarP(&opts.ProjectName, "project", "p", "", "Project name")

	return cmd
}

// Note: don't import other packages directly in this function, or it will be hard to mock and test
// If you want to add new dependencies, please add them in the Options struct

func runGet(f *cmdutil.Factory, opts Options) error {
	if opts.ID == "" || (opts.OwnerName == "" && opts.ProjectName == "") {
		return fmt.Errorf("please specify --id or both --owner and --project")
	}

	var (
		objID primitive.ObjectID
		err   error
	)

	if opts.ID != "" {
		objID, err = primitive.ObjectIDFromHex(opts.ID)
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}
	}

	project, err := f.ApiClient.GetProject(context.Background(), objID, opts.OwnerName, opts.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// todo: pretty print
	f.Log.Info(project)

	return nil
}
