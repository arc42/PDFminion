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
	// NOTE: Execute now triggers PersistentPreRun which loads configuration automatically
	err := rootCmd.Execute()
	assert.NoError(t, err, "Failed to execute command")

	// Verify viper got the flag values
	assert.Equal(t, "EN", viper.GetString("language"))
	assert.Equal(t, "testSourceDir", viper.GetString("source"))
	assert.Equal(t, "testTargetDir", viper.GetString("target"))
	assert.True(t, viper.GetBool("force"))
	assert.False(t, viper.GetBool("evenify"))

	// Verify the MinionConfig was populated by PersistentPreRun during Execute
	minionConfig := config.ActiveMinionConfig

	assert.Equal(t, domain.ParseLanguageCode("EN"), minionConfig.Language)
	assert.Equal(t, "testSourceDir", minionConfig.SourceDir)
	assert.Equal(t, "testTargetDir", minionConfig.TargetDir)
	assert.True(t, minionConfig.Force)
	assert.False(t, minionConfig.Evenify)
}

// TestFlagLoadingIntegration verifies that command-line flags are properly
// loaded into the configuration and override defaults
func TestFlagLoadingIntegration(t *testing.T) {
	// Test various flags to ensure they're properly loaded
	testCases := []struct {
		name        string
		args        []string
		checkConfig func(*testing.T, domain.MinionConfig)
		description string
	}{
		{
			name: "text_formatting_flags",
			args: []string{
				"--running-header", "My Custom Header",
				"--chapter-prefix", "Ch.",
				"--page-prefix", "P.",
				"--separator", " | ",
				"--page-count-prefix", "out of",
				"--blank-page-text", "This page intentionally blank",
			},
			checkConfig: func(t *testing.T, config domain.MinionConfig) {
				assert.Equal(t, "My Custom Header", config.RunningHeader, "running-header flag not loaded")
				assert.Equal(t, "Ch.", config.ChapterPrefix, "chapter-prefix flag not loaded")
				assert.Equal(t, "P.", config.PageNrPrefix, "page-prefix flag not loaded")
				assert.Equal(t, " | ", config.Separator, "separator flag not loaded")
				assert.Equal(t, "out of", config.PageCountPrefix, "page-count-prefix flag not loaded")
				assert.Equal(t, "This page intentionally blank", config.BlankPageText, "blank-page-text flag not loaded")
			},
			description: "Text formatting flags should be properly loaded",
		},
		{
			name: "processing_flags",
			args: []string{
				"--verbose",
				"--source", "/custom/source",
				"--target", "/custom/target",
				"--force",
				"--evenify=false",
			},
			checkConfig: func(t *testing.T, config domain.MinionConfig) {
				assert.True(t, config.Verbose, "verbose flag not loaded")
				assert.Equal(t, "/custom/source", config.SourceDir, "source flag not loaded")
				assert.Equal(t, "/custom/target", config.TargetDir, "target flag not loaded")
				assert.True(t, config.Force, "force flag not loaded")
				assert.False(t, config.Evenify, "evenify flag not loaded")
			},
			description: "Processing flags should be properly loaded",
		},
		{
			name: "language_flag",
			args: []string{
				"--language", "DE",
			},
			checkConfig: func(t *testing.T, config domain.MinionConfig) {
				assert.Equal(t, domain.ParseLanguageCode("DE"), config.Language, "language flag not loaded")
				// When language is set to German, German defaults should be used
				assert.Equal(t, "Kapitel", config.ChapterPrefix, "German chapter prefix not set")
				assert.Equal(t, "Seite", config.PageNrPrefix, "German page prefix not set")
			},
			description: "Language flag should load and set language-specific defaults",
		},
		{
			name: "shorthand_flags",
			args: []string{
				"-l", "FR",
				"-s", "/short/source",
				"-t", "/short/target",
				"-f",
				"-v",
			},
			checkConfig: func(t *testing.T, config domain.MinionConfig) {
				assert.Equal(t, domain.ParseLanguageCode("FR"), config.Language, "shorthand -l not working")
				assert.Equal(t, "/short/source", config.SourceDir, "shorthand -s not working")
				assert.Equal(t, "/short/target", config.TargetDir, "shorthand -t not working")
				assert.True(t, config.Force, "shorthand -f not working")
				assert.True(t, config.Verbose, "shorthand -v not working")
			},
			description: "Shorthand flags should work correctly",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a fresh root command for each test
			rootCmd := config.SetupApplication("test-version")
			rootCmd.SetArgs(tc.args)

			// Execute to parse flags
			// NOTE: Execute now triggers PersistentPreRun which loads configuration automatically
			err := rootCmd.Execute()
			assert.NoError(t, err, "Failed to execute command with args: %v", tc.args)

			// Get the configuration from the global variable set by PersistentPreRun
			minionConfig := config.ActiveMinionConfig

			// Run test-specific checks
			tc.checkConfig(t, minionConfig)
		})
	}
}
