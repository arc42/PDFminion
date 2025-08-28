package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"pdfminion/internal/domain"
)

// FlagChecker provides an interface for checking if flags have been set
type FlagChecker interface {
	HasBeenProvided(flagName string) bool
}

// CobraFlagChecker implements FlagChecker using Cobra commands
type CobraFlagChecker struct {
	cmd *cobra.Command
}

func NewCobraFlagChecker(cmd *cobra.Command) *CobraFlagChecker {
	return &CobraFlagChecker{cmd: cmd}
}

func (c *CobraFlagChecker) HasBeenProvided(flagName string) bool {
	flag := c.cmd.Flags().Lookup(flagName)
	if flag == nil {
		log.Warn().Msgf("Flag %s not found", flagName)
		return false
	}
	return flag.Changed
}

// ConfigureApplication collects configuration from all sources and merges them
// Priority order:
// 1. (lowest priority) default, depending on system language
// 2. Config file (if specified via --config)
// 3. (highest priority) command line flags
// For boolean flags: If a boolean value is not set in the configuration (file or flag),
// it must not override the previously set value.
// Therefore, we use metadata ("SetFields") to track which fields were explicitly set.
func ConfigureApplication(verbose bool, flagChecker FlagChecker) (domain.MinionConfig, error) {
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

	// 2. Load from config file if specified
	if flagChecker.HasBeenProvided("config") {
		configFile := viper.GetString("config")
		if configFile != "" {
			if fileConfig, err := loadConfigFile(configFile, verbose); err == nil {
				if verbose {
					fmt.Printf("Merging configuration from file: %s\n", configFile)
				}
				err := minionConfig.MergeWith(fileConfig)
				if err != nil {
					log.Warn().Err(err).Msg("Failed to merge file configuration")
				}
			} else {
				log.Warn().Err(err).Str("file", configFile).Msg("Failed to load config file")
			}
		}
	}

	// 3. Override with command line flags
	flagConfig := loadFlagConfig(flagChecker)
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

// loadFlagConfig loads the configuration from command line flags
func loadFlagConfig(flagChecker FlagChecker) domain.MinionConfig {
	log.Debug().Msg("loading flag configuration")

	// initialize local fconfig with empty metadata SetFields
	fconfig := domain.MinionConfig{
		SetFields: make(map[string]bool),
	}

	// to avoid the too-long-function linter error, we split the flag loading
	loadFlagProcessingConfig(fconfig, flagChecker)

	loadFlagTextOnPageConfig(fconfig, flagChecker)

	if flagChecker.HasBeenProvided("personal") {
		fconfig.PersonalTouch = viper.GetBool("personal")
		fconfig.SetFields["personal"] = true
	}

	return fconfig

}

func loadFlagTextOnPageConfig(fconfig domain.MinionConfig, flagChecker FlagChecker) {
	if flagChecker.HasBeenProvided("language") {
		fconfig.Language = domain.ParseLanguageCode(viper.GetString("language"))
		fconfig.SetFields["language"] = true
	}
	if flagChecker.HasBeenProvided("running-header") {
		fconfig.RunningHeader = viper.GetString("running-header")
		fconfig.SetFields["runningheader"] = true
	}
	if flagChecker.HasBeenProvided("chapter-prefix") {
		fconfig.ChapterPrefix = viper.GetString("chapter-prefix")
		fconfig.SetFields["chapterprefix"] = true
	}
	if flagChecker.HasBeenProvided("separator") {
		fconfig.Separator = viper.GetString("separator")
		fconfig.SetFields["separator"] = true
	}
	if flagChecker.HasBeenProvided("page-prefix") {
		fconfig.PageNrPrefix = viper.GetString("page-prefix")
		fconfig.SetFields["pageprefix"] = true
	}
	if flagChecker.HasBeenProvided("page-count-prefix") {
		fconfig.PageCountPrefix = viper.GetString("page-count-prefix")
		fconfig.SetFields["pagecountprefix"] = true
	}
	if flagChecker.HasBeenProvided("blank-page-text") {
		fconfig.BlankPageText = viper.GetString("blank-page-text")
		fconfig.SetFields["blankpagetext"] = true
	}
	// Add the TOC flag
	if flagChecker.HasBeenProvided("toc") {
		// We'll need to add this field to MinionConfig struct
		fconfig.SetFields["toc"] = true
	}
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

func loadFlagProcessingConfig(fconfig domain.MinionConfig, flagChecker FlagChecker) {
	if flagChecker.HasBeenProvided("verbose") {
		fconfig.Verbose = viper.GetBool("verbose")
		fconfig.SetFields["verbose"] = true
	}
	if flagChecker.HasBeenProvided("source") {
		fconfig.SourceDir = viper.GetString("source")
		fconfig.SetFields["sourcedir"] = true
	}
	if flagChecker.HasBeenProvided("target") {
		fconfig.TargetDir = viper.GetString("target")
		fconfig.SetFields["targetdir"] = true
	}
	if flagChecker.HasBeenProvided("force") {
		fconfig.Force = viper.GetBool("force")
		fconfig.SetFields["force"] = true
	}
	if flagChecker.HasBeenProvided("evenify") {
		fconfig.Evenify = viper.GetBool("evenify")
		fconfig.SetFields["evenify"] = true
	}
	if flagChecker.HasBeenProvided("merge") {
		fconfig.Merge = true
		fconfig.MergeFileName = viper.GetString("merge")
		fconfig.SetFields["merge"] = true
	}
	if flagChecker.HasBeenProvided("toc") {
		fconfig.TOC = viper.GetBool("toc")
		fconfig.SetFields["toc"] = true
	}
}
