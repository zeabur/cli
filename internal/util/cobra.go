package util

import (
	"github.com/spf13/cobra"
)

type CobraRunE func(cmd *cobra.Command, args []string) error

// RunEChain runs a chain of CobraRunE functions
func RunEChain(chains ...CobraRunE) CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		for _, chain := range chains {
			if err := chain(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}
