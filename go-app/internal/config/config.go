package config

import (
	"golang.org/x/text/language"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"pdfminion/internal/domain"
)

var (
	defaultConfig      *domain.MinionConfig
	homeDirFileConfig  *domain.MinionConfig
	localDirFileConfig *domain.MinionConfig
	flagConfig         *domain.MinionConfig
	minionConfig       *domain.MinionConfig
)

// LoadConfig collects configuration from all sources and merges them
// Priority order:
// 1. (lowest priority) default,  depending on system language,
// 2. home directory config file (if exists):
// .   2a: $HOME/pdfminion.yaml if no --config filename is given
// .   2b: $HOME/<filename>.yaml if --config filename is given
// 3. local directory config file (if exists),
// 4. (highest priority) command line flags
//
// In steps 2-4 we need to consider "unset" boolean values: If a boolean value is not set in the
// configuration (file or flag), it must not override the previously set value.
// Therefore we introduced metadata ("SetFields") to MinionConfig.
func LoadConfig() (*domain.MinionConfig, error) {

	log.Debug().Msg("starting LoadConfig")

	// 0. Start with system language detection
	systemLang := domain.MapSystemToAppLanguage()

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
	flagConfig := loadFlagConfig()

	// Merge configurations in priority order

	//minionConfig.MergeWith(homeDirFileConfig)
	//minionConfig.MergeWith(localDirFileConfig)
	err := minionConfig.MergeWith(flagConfig)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to merge flag configuration")
	}

	return minionConfig, nil
}

func loadDefaultConfig() *domain.MinionConfig {
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

// loadFlagConfig loads the configuration from command line flags
func loadFlagConfig() *domain.MinionConfig {
	log.Debug().Msg("loading flag configuration")

	// initialize local fconfig with empty metadata SetFields
	fconfig := domain.MinionConfig{
		SetFields: make(map[string]bool),
	}

	if flagHasBeenProvided("verbose") {
		fconfig.Verbose = viper.GetBool("verbose")
		fconfig.SetFields["verbose"] = true
	}
	if flagHasBeenProvided("source") {
		fconfig.SourceDir = viper.GetString("source")
		fconfig.SetFields["sourcedir"] = true
	}
	if flagHasBeenProvided("target") {
		fconfig.TargetDir = viper.GetString("target")
		fconfig.SetFields["targetdir"] = true
	}
	if flagHasBeenProvided("force") {
		fconfig.Force = viper.GetBool("force")
		fconfig.SetFields["force"] = true
	}
	if flagHasBeenProvided("evenify") {
		fconfig.Evenify = viper.GetBool("evenify")
		fconfig.SetFields["evenify"] = true
	}
	if flagHasBeenProvided("merge") {
		fconfig.Merge = true
		fconfig.MergeFileName = viper.GetString("merge")
		fconfig.SetFields["merge"] = true
	}
	if flagHasBeenProvided("config") {
		fconfig.ConfigFileName = viper.GetString("config")
		fconfig.SetFields["configfilename"] = true
	}
	if flagHasBeenProvided("language") {
		fconfig.Language = domain.ParseLanguageCode(viper.GetString("language"))
		fconfig.SetFields["language"] = true
	}
	if flagHasBeenProvided("running-header") {
		fconfig.RunningHeader = viper.GetString("running-header")
		fconfig.SetFields["runningheader"] = true
	}
	if flagHasBeenProvided("chapter-prefix") {
		fconfig.ChapterPrefix = viper.GetString("chapter-prefix")
		fconfig.SetFields["chapterprefix"] = true
	}
	if flagHasBeenProvided("separator") {
		fconfig.Separator = viper.GetString("separator")
		fconfig.SetFields["separator"] = true
	}
	if flagHasBeenProvided("page-prefix") {
		fconfig.PagePrefix = viper.GetString("page-prefix")
		fconfig.SetFields["pageprefix"] = true
	}
	if flagHasBeenProvided("page-count-prefix") {
		fconfig.PageCountPrefix = viper.GetString("page-count-prefix")
		fconfig.SetFields["pagecountprefix"] = true
	}
	if flagHasBeenProvided("blank-page-text") {
		fconfig.BlankPageText = viper.GetString("blank-page-text")
		fconfig.SetFields["blankpagetext"] = true
	}
	if flagHasBeenProvided("personal") {
		fconfig.PersonalTouch = viper.GetBool("personal")
		fconfig.SetFields["personal"] = true
	}

	return &fconfig

}

// flagHasBeenProvided checks if a flag has been explicitly provided on the command line
func flagHasBeenProvided(flagName string) bool {
	flag := rootCmd.Flags().Lookup(flagName)
	if flag == nil {
		log.Warn().Msgf("Flag %s not found", flagName)
		return false
	}
	return flag.Changed
}

/*func setupLogging(verbose bool) {
	level := zerolog.InfoLevel
	if verbose {
		level = zerolog.DebugLevel
	}

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
	}).Level(level).With().Timestamp().Logger()
}
*/
