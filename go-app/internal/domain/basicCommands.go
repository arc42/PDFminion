package domain

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
	"time"
)

// PrintFinalConfiguration prints the final configuration
func PrintFinalConfiguration(myConfig MinionConfig) {
	fmt.Println("Final Configuration:")
	fmt.Printf("Source directory: %s\n", myConfig.SourceDir)
	fmt.Printf("Target directory: %s\n", myConfig.TargetDir)
	fmt.Printf("Force: %t\n", myConfig.Force)
	fmt.Println(strings.Repeat("=", 20))
	fmt.Printf("Verbose: %t\n", myConfig.Verbose)
	fmt.Printf("Evenify: %t\n", myConfig.Evenify)
	fmt.Printf("Language: %s\n", myConfig.Language)
	fmt.Printf("Personal-touch: %t\n", myConfig.PersonalTouch)
	fmt.Println(strings.Repeat("=", 20))
	fmt.Printf("Running header: %s\n", myConfig.RunningHeader)
	fmt.Printf("Chapter prefix: %s\n", myConfig.ChapterPrefix)
	fmt.Printf("Separator: %s\n", myConfig.Separator)
	fmt.Printf("Page prefix: %s\n", myConfig.PagePrefix)
	fmt.Printf("Total page count prefix: %s\n", myConfig.PageCountPrefix)
	fmt.Printf("Blank page text: %s\n", myConfig.BlankPageText)
	fmt.Println(strings.Repeat("=", 20))
	fmt.Printf("Merge: %t\n", myConfig.Merge)
	fmt.Printf("Merge file name: %s\n", myConfig.MergeFileName)
}

//
// =============================================================================
//

// GiveCredits prints the credits
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

//
// =============================================================================
//

// Version and build information
var (
	// appVersion is set by main.go during initialization
	appVersion = "0.0.0"

	// buildTime and buildPlatform are set by the build process (Makefile)
	buildTime     = "unknown"
	buildPlatform = "unknown"
)

func SetAppVersion(v string) {
	appVersion = v
}

func PrintVersion() {
	fmt.Printf("PDFminion version %s\n", appVersion)
	fmt.Printf("Built on: %s\n", buildPlatform)
	t, err := time.Parse("2006 Jan 02 15:04", buildTime)
	if err == nil {
		formattedTime := formatBuildTime(t)
		fmt.Printf("Build time: %s\n", formattedTime)
	} else {
		fmt.Printf("Build time: %s\n", buildTime)
	}
}

func formatBuildTime(t time.Time) string {
	formatted := t.Format("2006 Jan 02 15:04")
	parts := strings.Split(formatted, " ")
	if len(parts) == 4 {
		day := parts[2]
		parts[2] = day + getDaySuffix(day)
		parts[3] += "h"
		return strings.Join(parts, " ")
	}
	return formatted
}

func getDaySuffix(day string) string {
	switch day {
	case "01", "21", "31":
		return "st"
	case "02", "22":
		return "nd"
	case "03", "23":
		return "rd"
	default:
		return "th"
	}
}
