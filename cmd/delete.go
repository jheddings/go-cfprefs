package cmd

import (
	"fmt"

	"github.com/jheddings/go-cfprefs"
	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <appID> <keypath>",
	Short: "Delete a preference value",
	Long: `Delete a preference value for the specified application ID.

The keypath can be a simple key name or include a JSON Pointer path 
(e.g., "config/server/port") to delete nested values.`,
	Args: cobra.ExactArgs(2),
	Run:  doDeleteCmd,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func doDeleteCmd(cmd *cobra.Command, args []string) {
	appID, keypath := args[0], args[1]

	tuiConfirm(fmt.Sprintf("Delete %s [%s]", keypath, appID))

	log.Trace().Str("app", appID).Str("keypath", keypath).Msg("Deleting preference value")

	err := cfprefs.Delete(appID, keypath)

	if err == nil {
		log.Info().Str("app", appID).Str("keypath", keypath).Msg("Value deleted successfully")
	} else {
		log.Fatal().Str("app", appID).Str("keypath", keypath).Err(err).Msg("Failed to delete value")
	}

	pterm.Success.Println("Value deleted")
}
