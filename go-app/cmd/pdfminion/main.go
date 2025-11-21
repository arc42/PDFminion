package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"pdfminion/internal/config"
)

// set the application version here.
const appVersion = "0.4.1a"

// debugLogging is used for development. It traces many function calls
// and shall facilitate debugging.
const debugLogging = true

func main() {
	setupLogging(debugLogging)

	// Set up the application
	rootCmd := config.SetupApplication(appVersion)

	// Run the root command
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("Failed to execute command")
		os.Exit(1)
	}

}

func setupLogging(debug bool) {
	// Always use UTC for timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05",
			NoColor:    false, // colorize output for debug level
		}).Level(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05",
			NoColor:    true,
		}).Level(zerolog.InfoLevel)
	}
}
