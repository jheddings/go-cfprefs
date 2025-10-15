package cmd

import (
	"fmt"

	"github.com/jheddings/go-cfprefs"
	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <appID> <key>",
	Short: "Delete a preference key (supports keypaths with '/' separator)",
	Long: `Delete a preference key for the specified application ID.

The key may be a keypath separated by forward slashes ("/") to delete values
from nested dictionaries. For example, "settings/display/brightness" will delete
the "brightness" key from the nested structure while preserving parent dictionaries.`,
	Args: cobra.ExactArgs(2),
	Run:  doDeleteCmd,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func doDeleteCmd(cmd *cobra.Command, args []string) {
	appID, key := args[0], args[1]
	tuiConfirm(fmt.Sprintf("Delete %s [%s]", key, appID))

	log.Trace().Str("app", appID).Str("key", key).Msg("Deleting preference key")
	err := cfprefs.Delete(appID, key)

	if err == nil {
		log.Info().Str("app", appID).Str("key", key).Msg("Key deleted successfully")
	} else {
		log.Fatal().Str("app", appID).Str("key", key).Err(err).Msg("Failed to delete key")
	}

	pterm.Success.Println("Key deleted")
}
