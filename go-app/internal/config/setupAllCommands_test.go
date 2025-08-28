package config_test

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"pdfminion/internal/config"
	"pdfminion/internal/domain"
	"testing"
)

func TestSetupFlags(t *testing.T) {
	// Initialize the root command and setup flags
	rootCmd := config.SetupApplication("test-version")

	// Set flag values
	args := []string{
		"--language", "EN",
		"--source", "testSourceDir",
		"--target", "testTargetDir",
		"--force",
		"--evenify=false",
	}
	rootCmd.SetArgs(args)

	// Execute the command to parse flags
	err := rootCmd.Execute()
	assert.NoError(t, err, "Failed to execute command")

	// Verify the configuration
	assert.Equal(t, "EN", viper.GetString("language"))
	assert.Equal(t, "testSourceDir", viper.GetString("source"))
	assert.Equal(t, "testTargetDir", viper.GetString("target"))
	assert.True(t, viper.GetBool("force"))
	assert.False(t, viper.GetBool("evenify"))

	// Verify the MinionConfig
	minionConfig, err := config.ConfigureApplication(false)
	assert.NoError(t, err, "Failed to configure application")

	assert.Equal(t, domain.ParseLanguageCode("EN"), minionConfig.Language)
	assert.Equal(t, "testSourceDir", minionConfig.SourceDir)
	assert.Equal(t, "testTargetDir", minionConfig.TargetDir)
	assert.True(t, minionConfig.Force)
	assert.False(t, minionConfig.Evenify)
}
