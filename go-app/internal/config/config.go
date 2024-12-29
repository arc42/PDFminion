package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
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

// LoadConfig initializes and returns the parsed configuration
func LoadConfig() (*domain.MinionConfig, error) {
	log.Debug().Msg("Starting configuration")

	defaultConfig = domain.NewDefaultConfig()

	// Parse individual configuration sources
	if err := loadDefaultConfig(); err != nil {
		return nil, fmt.Errorf("failed to load default configuration: %w", err)
	}
	if err := loadHomeDirFileConfig(); err != nil {
		log.Warn().Err(err).Msg("Failed to load home directory configuration")
	}
	if err := loadLocalDirFileConfig(); err != nil {
		log.Warn().Err(err).Msg("Failed to load local directory configuration")
	}
	if err := loadFlagConfig(); err != nil {
		return nil, fmt.Errorf("failed to load flag configuration: %w", err)
	}

	// Merge configurations in priority order
	minionConfig = domain.NewDefaultConfig()
	minionConfig.MergeWith(defaultConfig)
	minionConfig.MergeWith(homeDirFileConfig)
	minionConfig.MergeWith(localDirFileConfig)
	minionConfig.MergeWith(flagConfig)

	return minionConfig, nil
}

// ValidateConfig checks the configuration for correctness
func ValidateConfig(config *domain.MinionConfig) error {
	return config.Validate()
}

func loadDefaultConfig() error {
	log.Debug().Msg("loading default configuration")

	defaultConfig = domain.NewDefaultConfig()
	return nil
}

func loadHomeDirFileConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFile := viper.GetString("config")
	if configFile == "" {
		configFile = domain.MinionConfigFileName
	}

	homeConfigPath := filepath.Join(homeDir, configFile)
	if _, err := os.Stat(homeConfigPath); os.IsNotExist(err) {
		homeDirFileConfig = nil
		return nil
	}

	viper.SetConfigFile(homeConfigPath)
	if err := viper.MergeInConfig(); err != nil {
		return err
	}

	homeDirFileConfig = &domain.MinionConfig{}
	if err := viper.Unmarshal(homeDirFileConfig); err != nil {
		return err
	}

	return nil
}

func loadLocalDirFileConfig() error {
	configFile := viper.GetString("config")
	if configFile == "" {
		configFile = domain.MinionConfigFileName
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		localDirFileConfig = nil
		return nil
	}

	viper.SetConfigFile(configFile)
	if err := viper.MergeInConfig(); err != nil {
		return err
	}

	localDirFileConfig = &domain.MinionConfig{}
	if err := viper.Unmarshal(localDirFileConfig); err != nil {
		return err
	}

	return nil
}

func loadFlagConfig() error {
	flagConfig = &domain.MinionConfig{
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

	if langStr := viper.GetString("language"); langStr != "" {
		lang, err := domain.ValidateLanguage(langStr)
		if err != nil {
			return fmt.Errorf("invalid language specified: %w", err)
		}
		flagConfig.Language = lang
	}

	return nil
}

func setupLogging(verbose bool) {
	level := zerolog.InfoLevel
	if verbose {
		level = zerolog.DebugLevel
	}

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
	}).Level(level).With().Timestamp().Logger()
}
