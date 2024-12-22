package config

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"pdfminion/internal/domain"
)

// captureOutput captures stdout during a test
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestVersionCommand(t *testing.T) {
	output := captureOutput(func() {
		rootCmd.SetArgs([]string{"version"})
		rootCmd.Execute()
	})

	expected := "PDFMinion Version: " + domain.Version + "\n"
	if output != expected {
		t.Errorf("Version command output wrong.\nExpected: %q\nGot: %q", expected, output)
	}
}

func TestVersionFlag(t *testing.T) {
	// Set a known version for testing
	testVersion := "1.2.3"
	domain.SetVersion(testVersion)

	// Test with --version flag
	output := captureOutput(func() {
		rootCmd.SetArgs([]string{"--version"})
		rootCmd.Execute()
	})

	// Clean the output string
	output = strings.TrimSpace(output)
	expected := "PDFMinion Version: " + testVersion

	if output != expected {
		t.Errorf("Version flag output wrong.\nExpected: %q\nGot: %q", expected, output)
	}
}
