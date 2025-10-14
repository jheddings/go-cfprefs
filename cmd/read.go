package cmd

import (
	"fmt"

	"github.com/jheddings/go-cfprefs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read <appID> <key>",
	Short: "Read a preference value",
	Args:  cobra.ExactArgs(2),
	Run:   doReadCmd,
}

func init() {
	rootCmd.AddCommand(readCmd)
}

func doReadCmd(cmd *cobra.Command, args []string) {
	appID, key := args[0], args[1]
	log.Trace().Str("app", appID).Str("key", key).Msg("Reading preference")

	value, err := cfprefs.GetStr(appID, key)

	if err == nil {
		log.Info().Str("app", appID).Str("key", key).Str("value", value).Msg("Value read successfully")
	} else {
		log.Fatal().Err(err).Msg("Failed to read preference value")
	}

	fmt.Println(value)
}
