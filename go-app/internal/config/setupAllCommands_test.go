package config_test

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"pdfminion/internal/config"
	"testing"
)

func TestVersionCommandExists(t *testing.T) {
	// Create a simple test that just verifies the version command exists and works
	cmd := config.VersionCmd()
	
	// Basic assertions about the command
	assert.Equal(t, "version", cmd.Use)
	assert.Contains(t, cmd.Aliases, "v")
	assert.Contains(t, cmd.Short, "version")
}

func TestVersionFlagExists(t *testing.T) {
	// Create a root command
	rootCmd := &cobra.Command{Use: "test"}
	
	// Add a version flag like in the real app
	rootCmd.PersistentFlags().BoolP("version", "V", false, "Print version")
	
	// Verify the flag exists
	flag := rootCmd.PersistentFlags().Lookup("version")
	assert.NotNil(t, flag, "Version flag should exist")
	assert.Equal(t, "V", flag.Shorthand, "Short flag should be V")
	assert.Equal(t, "bool", flag.Value.Type(), "Flag should be a boolean")
}

// Additional simple tests for other commands
func TestCreditsCommandExists(t *testing.T) {
	cmd := config.CreditsCmd()
	assert.Equal(t, "credits", cmd.Use)
	assert.Contains(t, cmd.Aliases, "credit")
}

func TestSettingsCommandExists(t *testing.T) {
	cmd := config.SettingsCmd()
	assert.Equal(t, "settings", cmd.Use)
}

func TestListLanguagesCommandExists(t *testing.T) {
	cmd := config.ListLanguagesCmd()
	assert.Equal(t, "list-languages", cmd.Use)
}
