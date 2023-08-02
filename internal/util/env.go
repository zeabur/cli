package util

import "github.com/spf13/cobra"

// AddEnvParam todo: support name
func AddEnvParam(cmd *cobra.Command, id *string) {
	cmd.Flags().StringVar(id, "id", "", "Environment ID")
}

func AddEnvOfServiceParam(cmd *cobra.Command, id *string) {
	cmd.Flags().StringVar(id, "env-id", "", "Environment ID of service")
}

func AddEnvOfProjectParam(cmd *cobra.Command, id *string) {
	cmd.Flags().StringVar(id, "env-id", "", "Environment ID of project")
}
