package cmd

import (
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:                   "upgrade",
	Short:                 "Upgrade jumppad",
	Long:                  `Upgrade the jumppad binary, but leaves the stacks alone`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
