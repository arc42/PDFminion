package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pdfminion/internal/domain"
)

var (
	minionConfig *domain.MinionConfig
	rootCmd      = &cobra.Command{
		Use:   "pdfminion",
		Short: "PDFMinion adds page numbers to PDF files with custom options",
		Long:  "PDFMinion is a CLI tool to add page numbers to existing PDF files with customizable options like chapter numbers, running headers, and more",
		RunE:  runPDFProcessing,
	}
)

// LoadConfig initializes and returns the complete configuration
func LoadConfig() (*domain.MinionConfig, error) {
	setupFlags()
	setupCommands()

	// Initialize default configuration
	minionConfig = domain.NewDefaultConfig()

	// Execute command processing
	if err := rootCmd.Execute(); err != nil {
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	return minionConfig, nil
}

func setupFlags() {
	// Persistent flags (available to all commands)
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to configuration file")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
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

	// Mark required flags
	rootCmd.MarkFlagRequired("source")
	rootCmd.MarkFlagRequired("target")

	// Bind all flags to viper
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindPFlags(rootCmd.Flags())
}

func setupCommands() {
	// Version command - note the PersistentPreRun override
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show application version",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Override parent PersistentPreRun
		},
		Run: func(cmd *cobra.Command, args []string) {
			domain.PrintVersion()
			os.Exit(0)
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Add all immediate commands
	var immediateCommands = map[string]*cobra.Command{
		"help": {
			Use:   "help",
			Short: "Help about any command",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				// Override parent PersistentPreRun
			},
			Run: func(cmd *cobra.Command, args []string) {
				rootCmd.Help()
				os.Exit(0)
			},
		},
		"list-languages": {
			Use:   "list-languages",
			Short: "List supported languages",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				// Override parent PersistentPreRun
			},
			Run: func(cmd *cobra.Command, args []string) {
				domain.PrintLanguages()
				os.Exit(0)
			},
		},
		"ll": {
			Use:   "ll",
			Short: "List supported languages (short form)",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				// Override parent PersistentPreRun
			},
			Run: func(cmd *cobra.Command, args []string) {
				domain.PrintLanguages()
				os.Exit(0)
			},
		},
		"settings": {
			Use:   "settings",
			Short: "Display current configuration",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				// Override parent PersistentPreRun
			},
			Run: func(cmd *cobra.Command, args []string) {
				if err := loadConfiguration(); err != nil {
					log.Error().Err(err).Msg("Failed to load configuration")
					os.Exit(1)
				}
				domain.PrintFinalConfiguration(*minionConfig)
				os.Exit(0)
			},
		},
	}

	// Add all other immediate commands to root
	for _, cmd := range immediateCommands {
		rootCmd.AddCommand(cmd)
	}

	// Add preprocessing hooks
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return loadConfiguration()
	}
}

func loadConfiguration() error {
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

	// 5. Setup logging based on final config
	setupLogging(minionConfig.Debug)

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
		Debug:         viper.GetBool("debug"),
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

func setupLogging(debug bool) {
	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
	}).Level(level).With().Timestamp().Logger()
}

func runPDFProcessing(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting PDF processing")
	// TODO: Implement actual PDF processing logic
	return nil
}
