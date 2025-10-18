package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/jheddings/go-cfprefs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read <appID> <key>",
	Short: "Read a preference value (supports keypaths with '/' separator)",
	Long: `Read a preference value for the specified application ID.

The key may be a keypath separated by forward slashes ("/") to traverse nested
dictionaries. For example, "settings/display/brightness" will retrieve the
"brightness" value from the nested structure.`,
	Args: cobra.MinimumNArgs(1),
	Run:  doReadCmd,
}

func init() {
	rootCmd.AddCommand(readCmd)
}

func doReadCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal().Msg("App ID is required")
	}

	if len(args) == 1 {
		doReadKeysCmd(args)
	} else {
		doReadValueCmd(args)
	}
}

func doReadKeysCmd(args []string) {
	appID := args[0]
	log.Trace().Str("app", appID).Msg("Reading keys")

	keys, err := cfprefs.GetKeys(appID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read keys")
	}

	jsonBytes, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal keys to JSON")
	}

	fmt.Println(string(jsonBytes))
}

func doReadValueCmd(args []string) {
	appID, key := args[0], args[1]
	log.Trace().Str("app", appID).Str("key", key).Msg("Reading preference")

	value, err := cfprefs.Get(appID, key)

	if err == nil {
		log.Info().Str("app", appID).Str("key", key).Type("type", value).Msg("Value read successfully")
	} else {
		log.Fatal().Err(err).Msg("Failed to read preference value")
	}

	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal value to JSON")
	}

	fmt.Println(string(jsonBytes))
}
