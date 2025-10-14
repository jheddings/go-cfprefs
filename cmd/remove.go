package cmd

import (
	"fmt"

	"github.com/jheddings/go-cfprefs"
	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <appID> <key>",
	Short: "Remove a preference key",
	Args:  cobra.ExactArgs(2),
	Run:   doRemoveCmd,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func doRemoveCmd(cmd *cobra.Command, args []string) {
	appID, key := args[0], args[1]
	tuiConfirm(fmt.Sprintf("Remove %s [%s]", key, appID))

	log.Trace().Str("appID", appID).Str("key", key).Msg("Removing preference key")
	err := cfprefs.Delete(appID, key)

	if err == nil {
		log.Info().Str("app", appID).Str("key", key).Msg("Key removed successfully")
	} else {
		log.Fatal().Str("app", appID).Str("key", key).Err(err).Msg("Failed to remove key")
	}

	pterm.Success.Println("Key removed")
}
