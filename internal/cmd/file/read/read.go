package read

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	uploadID string
	path     string
}

func NewCmdRead(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "read <upload-id> <path>",
		Short: "Read a file from an upload",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.uploadID = args[0]
			opts.path = args[1]
			return runRead(f, opts)
		},
	}

	return cmd
}

func runRead(f *cmdutil.Factory, opts *Options) error {
	content, err := f.ApiClient.ReadUploadFile(context.Background(), opts.uploadID, opts.path)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(map[string]string{
			"path":    opts.path,
			"content": content,
		})
	}

	fmt.Print(content)

	return nil
}
