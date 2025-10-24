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
	Short: "Delete a preference key",
	Long: `Delete a preference key for the specified application ID.

Use the --query flag to apply JSONPath queries for more precise deletion within
complex nested structures.`,
	Args: cobra.ExactArgs(2),
	Run:  doDeleteCmd,
}

func init() {
	deleteCmd.Flags().StringP("query", "Q", "", "Apply JSONPath query for precise deletion")
	rootCmd.AddCommand(deleteCmd)
}

func doDeleteCmd(cmd *cobra.Command, args []string) {
	appID, key := args[0], args[1]
	query, _ := cmd.Flags().GetString("query")

	tuiConfirm(fmt.Sprintf("Delete %s [%s]", key, appID))

	log.Trace().Str("app", appID).Str("key", key).Str("query", query).Msg("Deleting preference key")

	var err error
	if query == "" {
		err = cfprefs.Delete(appID, key)
	} else {
		err = cfprefs.DeleteQ(appID, key, query)
	}

	if err == nil {
		log.Info().Str("app", appID).Str("key", key).Msg("Key deleted successfully")
	} else {
		log.Fatal().Str("app", appID).Str("key", key).Err(err).Msg("Failed to delete key")
	}

	pterm.Success.Println("Key deleted")
}
