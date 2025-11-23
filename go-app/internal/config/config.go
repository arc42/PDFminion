package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"pdfminion/internal/domain"
)

// ConfigureApplication collects configuration from all sources and merges them.
//
// SIMPLIFIED LAYERED CONFIGURATION APPROACH:
// The configuration is built in 3 simple layers:
//
//	Layer 1 (Base): System language-dependent defaults
//	  - Detects system language using go-locale library
//	  - Loads language-specific text defaults (e.g., German uses "Seite" for page)
//	  - Falls back to English if language is not supported
//
//	Layer 2 (Config File): Home directory config file (if exists)
//	  - Location: ~/.config/pdfminion/config.yaml
//	  - YAML format
//	  - Optional - if not found, no error
//	  - No --config flag needed - just put file in standard location
//
//	Layer 3 (CLI Flags): Command-line flags (highest priority)
//	  - Override all previous layers
//	  - Only flags that are explicitly set are applied
//
// METADATA TRACKING:
// The SetFields map tracks which configuration fields have been explicitly set.
// This is necessary for boolean flags to distinguish between:
//   - User explicitly set to false
//   - User didn't set at all (should use default)
//
// See ADR-0008 for architectural decision details.
func ConfigureApplication(verbose bool, cmd *cobra.Command) (domain.MinionConfig, error) {
	if verbose {
		fmt.Println("Configuring application...")
	}
	log.Debug().Msg("Starting ConfigureApplication")

	// 1. Start with language-dependent default configuration
	systemLang := domain.MapSystemToAppLanguage()
	if verbose {
		fmt.Printf("System language detected: %s\n", systemLang.String())
	}
	log.Debug().Str("language", systemLang.String()).Msg("detected")
	minionConfig := loadDefaultConfig()

	// 2. Load from home config file if it exists
	if homeConfig, err := loadHomeConfig(verbose); err == nil {
		if verbose {
			fmt.Printf("Merging configuration from: %s\n", GetConfigPath())
		}
		err := minionConfig.MergeWith(homeConfig)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to merge home config")
		}
	} else {
		log.Debug().Err(err).Msg("No home config file found (this is OK)")
	}

	// 3. Override with command line flags
	flagConfig := loadFlagConfig(cmd)
	if verbose {
		fmt.Println("Merging flag configuration")
	}
	err := minionConfig.MergeWith(flagConfig)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to merge flag configuration")
	}

	if verbose {
		domain.PrintFinalConfiguration(&minionConfig)
	}
	return minionConfig, nil
}

func loadDefaultConfig() domain.MinionConfig {
	log.Debug().Msg("loading default configuration")

	// Detect system language
	systemLang := domain.MapSystemToAppLanguage()

	// Check if detected language is supported
	tag, _, _ := domain.GetMatcher().Match(systemLang)
	if tag != systemLang {
		log.Debug().
			Str("detected", systemLang.String()).
			Msg(" language not supported, falling back to English")
		systemLang = language.English
	}

	return domain.NewDefaultConfig(systemLang)
}

// loadFlagConfig loads the configuration from command line flags.
// All flag loading is consolidated here for easier maintenance and review.
// Only flags that have been explicitly set by the user are loaded.
//
//nolint:funlen // This function is long but consolidating all flags in one place improves maintainability
func loadFlagConfig(cmd *cobra.Command) domain.MinionConfig {
	log.Debug().Msg("loading flag configuration")

	// Initialize config with empty metadata SetFields
	fconfig := domain.MinionConfig{
		SetFields: make(map[string]bool),
	}

	// Helper function to check if a flag was provided
	flagChanged := func(name string) bool {
		return cmd.Flags().Changed(name)
	}

	// ============================================================
	// PROCESSING FLAGS
	// ============================================================

	if flagChanged("verbose") {
		fconfig.Verbose = viper.GetBool("verbose")
		fconfig.SetFields["verbose"] = true
	}

	if flagChanged("source") {
		fconfig.SourceDir = viper.GetString("source")
		fconfig.SetFields["sourcedir"] = true
	}

	if flagChanged("target") {
		fconfig.TargetDir = viper.GetString("target")
		fconfig.SetFields["targetdir"] = true
	}

	if flagChanged("force") {
		fconfig.Force = viper.GetBool("force")
		fconfig.SetFields["force"] = true
	}

	if flagChanged("evenify") {
		fconfig.Evenify = viper.GetBool("evenify")
		fconfig.SetFields["evenify"] = true
	}

	// ============================================================
	// TEXT AND PAGE FORMATTING FLAGS
	// ============================================================

	if flagChanged("language") {
		fconfig.Language = domain.ParseLanguageCode(viper.GetString("language"))
		fconfig.SetFields["language"] = true
	}

	if flagChanged("running-header") {
		fconfig.RunningHeader = viper.GetString("running-header")
		fconfig.SetFields["runningheader"] = true
	}

	if flagChanged("chapter-prefix") {
		fconfig.ChapterPrefix = viper.GetString("chapter-prefix")
		fconfig.SetFields["chapterprefix"] = true
	}

	if flagChanged("separator") {
		fconfig.Separator = viper.GetString("separator")
		fconfig.SetFields["separator"] = true
	}

	if flagChanged("page-prefix") {
		fconfig.PageNrPrefix = viper.GetString("page-prefix")
		fconfig.SetFields["pageprefix"] = true
	}

	if flagChanged("page-count-prefix") {
		fconfig.PageCountPrefix = viper.GetString("page-count-prefix")
		fconfig.SetFields["pagecountprefix"] = true
	}

	if flagChanged("blank-page-text") {
		fconfig.BlankPageText = viper.GetString("blank-page-text")
		fconfig.SetFields["blankpagetext"] = true
	}

	// ============================================================
	// POST-PROCESSING FLAGS
	// ============================================================

	if flagChanged("merge") {
		fconfig.Merge = true
		fconfig.MergeFileName = viper.GetString("merge")
		fconfig.SetFields["merge"] = true
	}

	if flagChanged("toc") {
		fconfig.TOC = viper.GetBool("toc")
		fconfig.SetFields["toc"] = true
	}

	if flagChanged("personal") {
		fconfig.PersonalTouch = viper.GetBool("personal")
		fconfig.SetFields["personal"] = true
	}

	return fconfig
}

// GetConfigPath returns the path to the home config file
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get user home directory")
		return ""
	}
	return filepath.Join(homeDir, ".config", "pdfminion", "config.yaml")
}

// loadHomeConfig loads configuration from ~/.config/pdfminion/config.yaml
// Returns an error if the file doesn't exist (which is OK - not all users need a config file)
func loadHomeConfig(verbose bool) (domain.MinionConfig, error) {
	configPath := GetConfigPath()

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return domain.MinionConfig{}, err // File not found - this is OK
	}

	if verbose {
		fmt.Printf("Found config file: %s\n", configPath)
	}
	log.Debug().Str("file", configPath).Msg("Loading home config file")

	return loadConfigFile(configPath, verbose)
}

// loadConfigFile loads configuration from the specified YAML file
func loadConfigFile(configFile string, verbose bool) (domain.MinionConfig, error) {
	log.Debug().Str("file", configFile).Msg("Loading config file")

	// Create a new viper instance for file configuration
	v := viper.New()
	v.SetConfigFile(configFile)

	if err := v.ReadInConfig(); err != nil {
		return domain.MinionConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}

	// Initialize config with empty metadata
	config := domain.MinionConfig{
		SetFields: make(map[string]bool),
	}

	// Map file values to config struct
	if v.IsSet("language") {
		config.Language = domain.ParseLanguageCode(v.GetString("language"))
		config.SetFields["language"] = true
	}

	if v.IsSet("source") {
		config.SourceDir = v.GetString("source")
		config.SetFields["sourcedir"] = true
	}

	if v.IsSet("target") {
		config.TargetDir = v.GetString("target")
		config.SetFields["targetdir"] = true
	}

	if v.IsSet("force") {
		config.Force = v.GetBool("force")
		config.SetFields["force"] = true
	}

	if v.IsSet("evenify") {
		config.Evenify = v.GetBool("evenify")
		config.SetFields["evenify"] = true
	}

	if v.IsSet("verbose") {
		config.Verbose = v.GetBool("verbose")
		config.SetFields["verbose"] = true
	}

	if v.IsSet("running-header") {
		config.RunningHeader = v.GetString("running-header")
		config.SetFields["runningheader"] = true
	}

	if v.IsSet("chapter-prefix") {
		config.ChapterPrefix = v.GetString("chapter-prefix")
		config.SetFields["chapterprefix"] = true
	}

	if v.IsSet("page-prefix") {
		config.PageNrPrefix = v.GetString("page-prefix")
		config.SetFields["pageprefix"] = true
	}

	if v.IsSet("page-count-prefix") {
		config.PageCountPrefix = v.GetString("page-count-prefix")
		config.SetFields["pagecountprefix"] = true
	}

	if v.IsSet("blank-page-text") {
		config.BlankPageText = v.GetString("blank-page-text")
		config.SetFields["blankpagetext"] = true
	}

	if v.IsSet("separator") {
		config.Separator = v.GetString("separator")
		config.SetFields["separator"] = true
	}

	if v.IsSet("merge") {
		config.Merge = true
		config.MergeFileName = v.GetString("merge")
		config.SetFields["merge"] = true
	}

	if v.IsSet("personal") {
		config.PersonalTouch = v.GetBool("personal")
		config.SetFields["personal"] = true
	}

	if v.IsSet("toc") {
		config.TOC = v.GetBool("toc")
		config.SetFields["toc"] = true
	}

	return config, nil
}
