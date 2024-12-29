package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pdfminion/internal/domain"
)

// Package-level variable accessible to all commands

var (
	rootCmd = &cobra.Command{
		Use:   "pdfminion",
		Short: "PDFMinion adds page numbers to PDF files with custom options",
		Long:  "PDFMinion is a CLI tool to add page numbers to existing PDF files with customizable options like chapter numbers, running headers, and more",
		RunE:  runPDFProcessing,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip config loading for immediate commands
			if isImmediateCommand(cmd) {
				return nil
			}
			return loadConfiguration(cmd)
		},
	}
)

func SetupApplication(appVersion string) *cobra.Command {
	// Set application version
	domain.SetAppVersion(appVersion)

	// Setup flags, as most commands need the flags prior to execution
	setupFlags()

	// Setup commands for the root CLI application
	setupCommands()

	log.Debug().Msg("Command setup complete")

	// Load configuration in layers (defaults, config-file(s), flags on CLI)
	minionConfig, err := LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Error loading configuration")
		os.Exit(1)
	}
	log.Debug().Interface("config:", minionConfig).Msg("Configuration completed ")

	// Validate configuration only for PDF processing, skip for immediate commands
	//if err := ValidateConfig(minionConfig); err != nil {
	//	log.Error().Err(err).Msg("Configuration validation failed")
	//	os.Exit(1)
	//s}

	// Return the root command for execution

	return rootCmd
}
func setupFlags() {
	// Persistent flags (available to all commands)
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to configuration file")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose mode, creates detailed processing output")
	rootCmd.PersistentFlags().StringP("language", "l", "", "Override system language")

	// Local flags (only for PDF processing)
	rootCmd.Flags().StringP("source", "s", domain.DefaultSourceDir, "Source directory for PDF files")
	rootCmd.Flags().StringP("target", "t", domain.DefaultTargetDir, "Target directory for processed files")
	rootCmd.Flags().Bool("force", false, "Force overwrite of target directory")
	rootCmd.Flags().Bool("evenify", true, "Ensure even page count in output")
	rootCmd.Flags().Bool("merge", false, "Merge all output files into one")
	rootCmd.Flags().String("merge-filename", "merged.pdf", "Name for merged output file")
	rootCmd.Flags().String("running-header", "", "Text for running header")
	rootCmd.Flags().String("chapter-prefix", domain.ChapterPrefix, "Prefix for chapter numbers")
	rootCmd.Flags().String("separator", domain.ChapterPageSeparator, "Separator between chapter and page")
	rootCmd.Flags().String("page-prefix", domain.PageNrPrefix, "Prefix for page numbers")
	rootCmd.Flags().String("total-page-count-prefix", "", "Prefix for total page count")
	rootCmd.Flags().String("blank-page-text", "", "Text for blank pages")

	// Bind all flags to viper
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind persistent flags to viper")
	}

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind flags to viper")
	}

	log.Debug().Msg("setupFlags completed")

}

func setupCommands() {
	rootCmd.AddCommand(
		VersionCmd(),
		CreditsCmd(),
		SettingsCmd(),
	)
	log.Debug().Msg("Add commands completed")

}

func isImmediateCommand(cmd *cobra.Command) bool {
	immediateCommands := []string{"version", "credits", "help"}
	for _, c := range immediateCommands {
		if cmd.Name() == c {
			return true
		}
	}
	return false
}

func loadConfiguration(cmd *cobra.Command) error {
	// 1. Start with system language detection
	systemLang := domain.DetectSystemLanguage()
	minionConfig.Language = systemLang

	// 2. Load config files
	if err := loadConfigFiles(); err != nil {
		return fmt.Errorf("failed to load config files: %w", err)
	}

	// 3. Apply command line flags (highest priority)
	if err := applyCommandLineFlags(); err != nil {
		return fmt.Errorf("failed to apply command line flags: %w", err)
	}

	// 4. Validate final configuration
	if err := minionConfig.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

func loadConfigFiles() error {
	configFile := viper.GetString("config")
	if configFile == "" {
		configFile = domain.MinionConfigFileName
	}

	// Try home directory first
	if homeConfig, err := loadHomeConfig(configFile); err == nil {
		if err := minionConfig.MergeWith(homeConfig); err != nil {
			log.Warn().Err(err).Msg("Failed to merge home directory config")
		}
	}

	// Then try current directory (takes precedence)
	if localConfig, err := loadLocalConfig(configFile); err == nil {
		if err := minionConfig.MergeWith(localConfig); err != nil {
			log.Warn().Err(err).Msg("Failed to merge local directory config")
		}
	}

	return nil
}

func loadHomeConfig(filename string) (*domain.MinionConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	viper.SetConfigFile(filepath.Join(homeDir, filename))
	if err := viper.MergeInConfig(); err != nil {
		return nil, err
	}

	var config domain.MinionConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadLocalConfig(filename string) (*domain.MinionConfig, error) {
	viper.SetConfigFile(filename)
	if err := viper.MergeInConfig(); err != nil {
		return nil, err
	}

	var config domain.MinionConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func applyCommandLineFlags() error {
	// Create a new config from flags
	flagConfig := &domain.MinionConfig{
		Verbose:       viper.GetBool("verbose"),
		SourceDir:     viper.GetString("source"),
		TargetDir:     viper.GetString("target"),
		Force:         viper.GetBool("force"),
		Evenify:       viper.GetBool("evenify"),
		Merge:         viper.GetBool("merge"),
		MergeFileName: viper.GetString("merge-filename"),
		RunningHeader: viper.GetString("running-header"),
		ChapterPrefix: viper.GetString("chapter-prefix"),
		PagePrefix:    viper.GetString("page-prefix"),
	}

	// Override language if specified
	if langStr := viper.GetString("language"); langStr != "" {
		lang, err := domain.ValidateLanguage(langStr)
		if err != nil {
			return fmt.Errorf("invalid language specified: %w", err)
		}
		flagConfig.Language = lang
	}

	return minionConfig.MergeWith(flagConfig)
}

func runPDFProcessing(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting PDF processing")
	// TODO: Implement actual PDF processing logic
	return nil
}
