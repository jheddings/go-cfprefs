package cmd

import (
	"github.com/jheddings/go-cfprefs"
	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var writeCmd = &cobra.Command{
	Use:   "write <appID> <key> <value>",
	Short: "Write a preference value",
	Args:  cobra.ExactArgs(3),
	Run:   doWriteCmd,
}

func init() {
	rootCmd.AddCommand(writeCmd)
}

func doWriteCmd(cmd *cobra.Command, args []string) {
	appID, key, value := args[0], args[1], args[2]
	log.Trace().Str("appID", appID).Str("key", key).Str("value", value).Msg("Writing preference")

	err := cfprefs.Set(appID, key, value)

	if err == nil {
		log.Info().Str("app", appID).Str("key", key).Str("value", value).Msg("Value saved successfully")
	} else {
		log.Fatal().Err(err).Msg("Failed to write preference value")
	}

	pterm.Success.Println("Key written")
}
