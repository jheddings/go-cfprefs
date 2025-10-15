package cmd

import (
	"strconv"
	"time"

	"github.com/jheddings/go-cfprefs"
	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	writeTypeStr   bool
	writeTypeInt   bool
	writeTypeFloat bool
	writeTypeBool  bool
	writeTypeDate  bool
)

var writeCmd = &cobra.Command{
	Use:   "write <appID> <key> <value>",
	Short: "Write a preference value (supports keypaths with '/' separator)",
	Long: `Write a preference value for the specified application ID.

The key may be a keypath separated by forward slashes ("/") to set values in
nested dictionaries. For example, "settings/display/brightness" will set the
"brightness" value in the nested structure, creating intermediate dictionaries
as needed.`,
	Args: cobra.ExactArgs(3),
	Run:  doWriteCmd,
}

func init() {
	writeCmd.Flags().BoolVar(&writeTypeStr, "string", false, "Parse value as string (default)")
	writeCmd.Flags().BoolVar(&writeTypeInt, "int", false, "Parse value as integer")
	writeCmd.Flags().BoolVar(&writeTypeFloat, "float", false, "Parse value as float")
	writeCmd.Flags().BoolVar(&writeTypeBool, "bool", false, "Parse value as boolean")
	writeCmd.Flags().BoolVar(&writeTypeDate, "date", false, "Parse value as date (ISO 8601 format)")

	rootCmd.AddCommand(writeCmd)
}

func doWriteCmd(cmd *cobra.Command, args []string) {
	appID, key, valueStr := args[0], args[1], args[2]

	// make sure only one type flag is set
	typeCount := 0
	if writeTypeStr {
		typeCount++
	}
	if writeTypeInt {
		typeCount++
	}
	if writeTypeFloat {
		typeCount++
	}
	if writeTypeBool {
		typeCount++
	}
	if writeTypeDate {
		typeCount++
	}
	if typeCount > 1 {
		log.Fatal().Msg("Only one type flag may be specified")
	}

	value := parseValue(valueStr)
	log.Trace().Str("appID", appID).Str("key", key).Any("value", value).Type("type", value).Msg("Writing preference")

	err := cfprefs.Set(appID, key, value)
	if err == nil {
		log.Info().Str("app", appID).Str("key", key).Any("value", value).Msg("Value saved successfully")
	} else {
		log.Fatal().Err(err).Msg("Failed to write preference value")
	}

	pterm.Success.Println("Key written")
}

func parseValue(valueStr string) any {
	if writeTypeInt {
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse value as integer")
		}
		return value
	}

	if writeTypeFloat {
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse value as float")
		}
		return value
	}

	if writeTypeBool {
		value, err := strconv.ParseBool(valueStr)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse value as boolean")
		}
		return value
	}

	if writeTypeDate {
		formats := []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"}

		for _, format := range formats {
			value, err := time.Parse(format, valueStr)
			if err == nil {
				return value
			}
		}
		log.Fatal().Msg("Failed to parse value as date")
	}

	return valueStr
}
