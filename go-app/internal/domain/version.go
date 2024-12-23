package domain

import (
	"fmt"
	"strings"
	"time"
)

// Version and build information
var (
	// appVersion is set by main.go during intialization
	appVersion = "0.0.0"

	// buildTime and buildPlatform are set by the build process (Makefile)
	buildTime     = "unknown"
	buildPlatform = "unknown"
)

// SetAppVersion sets the application version
func SetAppVersion(v string) {
	appVersion = v
}

func PrintVersion() {
	fmt.Printf("PDFminion version %s\n", appVersion)
	fmt.Printf("Built on: %s\n", buildPlatform)
	if buildTime != "" {
		t, err := time.Parse("2006 Jan 02 15:04", buildTime)
		if err == nil {
			formattedTime := formatBuildTime(t)
			fmt.Printf("Build time: %s\n", formattedTime)
		} else {
			fmt.Printf("Build time: %s\n", buildTime)
		}
	} else {
		fmt.Println("Build time: Not available")
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
