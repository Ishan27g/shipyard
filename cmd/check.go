package cmd

import (
	"fmt"

	"github.com/jumppad-labs/jumppad/pkg/clients"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks the system to ensure required dependencies are installed",
	Long:  `Checks the system to ensure required dependencies are installed`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		s := clients.SystemImpl{}
		o, _ := s.Preflight()

		fmt.Println("")
		fmt.Println("###### SYSTEM DIAGNOSTICS ######")
		fmt.Println(o)
	},
}
