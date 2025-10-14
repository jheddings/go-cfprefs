package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tuiAlert  = pterm.NewStyle(pterm.BgRed, pterm.FgBlack, pterm.Bold)
	tuiPrompt = pterm.NewStyle(pterm.FgRed, pterm.Bold)
)

var rootCmd = &cobra.Command{
	Use:              "cfprefs",
	Short:            "Utility wrapper for CFPreferences on macOS",
	PersistentPreRun: initLogging,
}

func init() {
	pFlags := rootCmd.PersistentFlags()

	pFlags.BoolP("yes", "y", false, "Assume 'yes' for confirmation prompts")
	viper.BindPFlag("yes", pFlags.Lookup("yes"))

	pFlags.CountP("verbose", "v", "Increase verbosity in logging")
	viper.BindPFlag("verbose", pFlags.Lookup("verbose"))

	pFlags.BoolP("quiet", "q", false, "Only log errors and warnings (override verbose)")
	viper.BindPFlag("quiet", pFlags.Lookup("quiet"))
}

func initLogging(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetCount("verbose")
	quiet, _ := cmd.Flags().GetBool("quiet")

	if quiet {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if verbose > 2 {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else if verbose > 1 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if verbose > 0 {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute root command")
	}
}

func tuiConfirm(message string) {
	if force, _ := rootCmd.Flags().GetBool("force"); force {
		return
	}

	var response string

	if message != "" {
		tuiAlert.Println(message)
	}

	tuiPrompt.Print("Continue? [y/N]: ")
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		log.Info().Msg("Aborted by user")
		os.Exit(1)
	}
}
