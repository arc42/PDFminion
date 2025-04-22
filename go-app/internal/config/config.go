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
// 1. (lowest priority) default,  depending on system language,
// 2. home directory config file (if exists):
// .   2a: $HOME/pdfminion.yaml if no --config filename is given
// .   2b: $HOME/<filename>.yaml if --config filename is given
// 3. local directory config file (if exists),
// 4. (highest priority) command line flags
// In steps 2-4 we need to consider "unset" boolean values: If a boolean value is not set in the
// configuration (file or flag), it must not override the previously set value.
// Therefore, we introduced metadata ("SetFields") to ActiveMinionConfig.
func ConfigureApplication(verbose bool, flagChecker FlagChecker) (domain.MinionConfig, error) {

	if verbose {
		fmt.Println("starting to configure  application")
	}
	log.Debug().Msg("starting ConfigureApplication")

	// 0. Start with system language detection
	systemLang := domain.MapSystemToAppLanguage()
	if verbose {
		fmt.Printf("System language detected: %s\n", systemLang.String())
	}
	log.Debug().Str("language", systemLang.String()).Msg("detected")

	// Parse individual configuration sources
	// --------------------------------------

	// 1. Start with language-dependent default configuration
	minionConfig := loadDefaultConfig()

	// 2. $HOME directory configuration file
	/*	if homeDirFileConfig, err := loadHomeConfig(); err != nil {
			log.Warn().Err(err).Msg("Failed to load home directory configuration")
		}
	*/
	/*	// 3. Local directory configuration file
		if localDirFileConfig, err := localDirFileConfig(); err != nil {
			log.Warn().Err(err).Msg("Failed to load local directory configuration")
		}
	*/
	// 4. Command line flags (with the exception of --language)
	flagConfig := loadFlagConfig(flagChecker)

	// Merge configurations in priority order

	// ActiveMinionConfig.MergeWith(homeDirFileConfig)
	// ActiveMinionConfig.MergeWith(localDirFileConfig)

	if verbose {
		fmt.Println("merging default and flag configuration")
	}

	err := minionConfig.MergeWith(flagConfig)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to merge flag configuration")
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
}
