package main

import (
	"os"

	"github.com/jheddings/go-cfprefs/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("cfprefs.yaml")

	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Debug().Msg("Config file not found, using default values")
		} else {
			// Config file was found but another error was produced
			log.Error().Err(err).Msg("Failed to read config file")
		}
	}
}

func main() {
	cmd.Execute()
}
