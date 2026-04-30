package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	uploadID string
	path     string
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "list <upload-id> [path]",
		Short:   "List files in an upload",
		Aliases: []string{"ls"},
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.uploadID = args[0]
			if len(args) > 1 {
				opts.path = args[1]
			}
			return runList(f, opts)
		},
	}

	return cmd
}

func runList(f *cmdutil.Factory, opts *Options) error {
	var pathPtr *string
	if opts.path != "" {
		pathPtr = &opts.path
	}

	files, err := f.ApiClient.ListUploadFiles(context.Background(), opts.uploadID, pathPtr)
	if err != nil {
		return fmt.Errorf("list files failed: %w", err)
	}

	if len(files) == 0 {
		if f.JSON {
			return f.Printer.JSON([]any{})
		}
		f.Log.Infof("No files found")
		return nil
	}

	if f.JSON {
		return f.Printer.JSON(files)
	}

	fmt.Println(strings.Join(files, "\n"))

	return nil
}
