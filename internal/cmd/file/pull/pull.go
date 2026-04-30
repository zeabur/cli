package pull

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

// NewCmdPull creates the file pull command.
func NewCmdPull(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull <upload-id> [target-dir]",
		Short: "Download uploaded project files to local directory",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			uploadID := ""
			targetDir := "."
			if len(args) > 0 {
				uploadID = args[0]
			}
			if len(args) > 1 {
				targetDir = args[1]
			}
			return runPull(cmd, f, uploadID, targetDir)
		},
	}

	return cmd
}

func runPull(cmd *cobra.Command, f *cmdutil.Factory, uploadID string, targetDir string) error {
	if uploadID == "" {
		if !f.Interactive {
			return fmt.Errorf("upload-id is required")
		}
		id, err := f.Prompter.Input("Enter upload ID:", "")
		if err != nil {
			return err
		}
		uploadID = strings.TrimSpace(id)
		if uploadID == "" {
			return fmt.Errorf("upload-id is required")
		}
	}

	count, err := f.ApiClient.PullUploadFiles(cmd.Context(), uploadID, targetDir)
	if err != nil {
		return fmt.Errorf("pull files failed: %w", err)
	}

	f.Log.Infof("Pulled %d files to %s", count, targetDir)

	return nil
}
