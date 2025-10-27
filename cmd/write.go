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
	Short: "Write a preference value",
	Long: `Write a preference value for the specified application ID.

Use the --query flag to apply JSONPath queries for more precise value setting
within complex nested structures.`,
	Args: cobra.ExactArgs(3),
	Run:  doWriteCmd,
}

func init() {
	flags := writeCmd.Flags()

	flags.BoolVar(&writeTypeStr, "string", false, "Parse value as string (default)")
	flags.BoolVar(&writeTypeInt, "int", false, "Parse value as integer")
	flags.BoolVar(&writeTypeFloat, "float", false, "Parse value as float")
	flags.BoolVar(&writeTypeBool, "bool", false, "Parse value as boolean")
	flags.BoolVar(&writeTypeDate, "date", false, "Parse value as date (ISO 8601 format)")
	flags.StringP("query", "Q", "", "Apply JSONPath query for precise value setting")

	rootCmd.AddCommand(writeCmd)
}

func doWriteCmd(cmd *cobra.Command, args []string) {
	appID, key, valueStr := args[0], args[1], args[2]
	query, _ := cmd.Flags().GetString("query")

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

	log.Trace().
		Str("appID", appID).
		Str("key", key).
		Str("query", query).
		Any("value", value).
		Type("type", value).
		Msg("Writing preference")

	var err error
	if query == "" {
		err = cfprefs.Set(appID, key, value)
	} else {
		err = cfprefs.SetQ(appID, key, query, value)
	}

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
