package domain

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

// GiveCredits prints a list of credits to the console in coloured output
func GiveCredits() {

	log.Logger = zerolog.New(coloredWriter(os.Stdout)).
		Level(zerolog.InfoLevel).
		With().Logger()

	// the real credits
	log.Info().Msg("PDFMinion (by Gernot Starke) stands on the shoulder of giants 🚀. Thanx to:")
	log.Info().Msg("Horst Rutter (pdfcpu, https://pdfcpu.io)")
	log.Info().Msg("Steve Francia (cobra + viper, https://cobra.dev)")
	log.Info().Msg("Xuanwo 👋(go-locale, https://github.com/Xuanwo/go-locale)")
	log.Info().Msg("Max Howell (Homebrew, https://brew.sh)")
	log.Info().Msg("Tom Preston-Werner & Co. (Jekyll, https://jekyllrb.com)")
	log.Info().Msg("Special thanks to the Go team (https://golang.org) and the Go community.")
	log.Info().Msg("INNOQ (https://innoq.com) for supporting my work")

}

var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
	"\033[37m", // White
}

const resetColor = "\033[0m"

// Custom writer that adds color and removes unnecessary fields
func coloredWriter(out io.Writer) zerolog.ConsoleWriter {
	colorIndex := 0
	return zerolog.ConsoleWriter{
		Out: out,
		// Format only the message
		FormatMessage: func(i interface{}) string {
			color := colors[colorIndex%len(colors)]
			colorIndex++
			return fmt.Sprintf("%s%s%s", color, i, resetColor)
		},
		// Remove log level, timestamp, and fields
		FormatLevel:      func(i interface{}) string { return "" },
		FormatFieldName:  func(i interface{}) string { return "" },
		FormatFieldValue: func(i interface{}) string { return "" },
		FormatTimestamp:  func(i interface{}) string { return "" }, // Remove the <nil>
		NoColor:          false,
	}
}
