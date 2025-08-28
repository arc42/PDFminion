package domain_test

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"io"
	"os"
	"pdfminion/internal/domain"
	"testing"
)

// captureOutputAndCallFunction captures the output of the PrintFinalConfiguration function
// that writes to stdout. It returns the output as a string.
func captureOutputAndCallFunction(config *domain.MinionConfig) (string, error) {
	// Capture the output
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	defer func() {
		os.Stdout = old
		if err := w.Close(); err != nil {
			log.Error().Err(err).Msg("error closing pipe writer")
		}
	}()

	// Call the function
	domain.PrintFinalConfiguration(config)

	// Read the output
	if err := w.Close(); err != nil {
		log.Error().Err(err).Msg("error closing pipe writer")
		return "", err
	}
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), nil
}

func TestPrintFinalConfigurationWithDefaultSourceDir(t *testing.T) {
	myConfig := domain.NewDefaultEnglishConfig()

	output, err := captureOutputAndCallFunction(&myConfig)
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	// Test the output for standard source directory

	// Test the output for standard source directory
	expectedSource := domain.DefaultSourceDir
	assert.Contains(t, output, expectedSource, "expected output to contain '%s'", expectedSource)
}

func TestPrintFinalConfigurationWithOtherSourceDir(t *testing.T) {
	myConfig := domain.NewDefaultEnglishConfig()
	otherSource := "abcotherSource42"
	myConfig.SourceDir = otherSource

	output, err := captureOutputAndCallFunction(&myConfig)
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	// Test the output for standard source directory

	// Test the output for standard source directory
	assert.Contains(t, output, otherSource, "expected output to contain '%s'", otherSource)
}

func TestPrintFinalConfigurationWithOtherLanguage(t *testing.T) {
	myConfig := domain.NewDefaultEnglishConfig()
	otherLanguage := language.Icelandic
	myConfig.Language = otherLanguage

	output, err := captureOutputAndCallFunction(&myConfig)
	if err != nil {
		t.Fatalf("error capturing output: %v", err)
	}
	// Test the output for standard source directory

	// Test the output for standard source directory
	assert.Contains(t, output, otherLanguage, "expected output to contain '%s'", otherLanguage)
}
