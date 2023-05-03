package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jumppad-labs/jumppad/pkg/clients"
	"github.com/jumppad-labs/jumppad/pkg/utils"
	"github.com/spf13/cobra"
)

// ErrorInvalidBlueprintURI is returned when the URI for a blueprint can not be parsed
var ErrorInvalidBlueprintURI = errors.New("error invalid Blueprint URI, blueprints should be formatted 'github.com/org/repo//blueprint'")

func newGetCmd(bp clients.Getter) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "get [remote blueprint]",
		Short: "Download the blueprint to the Shipyard config folder",
		Long:  `Download the blueprint to the Shipyard configuration folder`,
		Example: `
  # Fetch a blueprint from GitHub
  yard get github.com/shipyard-run/blueprints//vault-k8s
	`,
		Args:         cobra.ArbitraryArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// create the shipyard home
			os.MkdirAll(utils.ShipyardHome(), os.FileMode(0755))

			// check the number of args
			if len(args) != 1 {
				return fmt.Errorf("Command takes a single argument")
			}

			bp.SetForce(force)

			var err error
			dst := args[0]
			cmd.Println("Fetching blueprint from: ", dst)
			cmd.Println("")

			if utils.IsLocalFolder(dst) {
				return fmt.Errorf("Parameter is not a remote blueprint, e.g. github.com/shipyard-run/blueprints//vault-k8s")
			}

			// fetch the remote server from github
			err = bp.Get(dst, utils.GetBlueprintLocalFolder(dst))
			if err != nil {
				return fmt.Errorf("Unable to retrieve blueprint: %s", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force-update", "", false, "When set to true Shipyard will ignore cached images, or files and will download")
	return cmd
}
