package read

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	uploadID string
	path     string
}

// NewCmdRead creates the file read command.
func NewCmdRead(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read <upload-id> <path>",
		Short: "Read a file from an upload",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &Options{}
			if len(args) > 0 {
				opts.uploadID = args[0]
			}
			if len(args) > 1 {
				opts.path = args[1]
			}
			return runRead(cmd, f, opts)
		},
	}

	return cmd
}

func runRead(cmd *cobra.Command, f *cmdutil.Factory, opts *Options) error {
	if opts.uploadID == "" {
		if !f.Interactive {
			return fmt.Errorf("upload-id is required")
		}
		id, err := f.Prompter.Input("Enter upload ID:", "")
		if err != nil {
			return err
		}
		opts.uploadID = id
	}

	if opts.path == "" {
		if !f.Interactive {
			return fmt.Errorf("path is required")
		}
		p, err := f.Prompter.Input("Enter file path:", "")
		if err != nil {
			return err
		}
		opts.path = p
	}

	content, err := f.ApiClient.ReadUploadFile(cmd.Context(), opts.uploadID, opts.path)
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
