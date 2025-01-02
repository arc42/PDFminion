package domain

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

const (
	MinionConfigFileName = "pdfminion.yaml"

	DefaultSourceDir       = "_pdfs"
	DefaultTargetDir       = "_target"
	DefaultPageNrPrefix    = "Page"
	DefaultChapterPrefix   = "Chapter"
	DefaultSeparator       = " - "
	DefaultPageCountPrefix = "of"
	DefaultBlankPageText   = "Intentionally left blank"
	DefaultMergeFileName   = "merged.pdf"
)

// MinionConfig holds the configuration for the PDFMinion application
// Several XYValid fields are used to check if the respective values hold valid values.
// Certain operations are possible with invalid flags, as we can fall back to defaults.
// However, some operations (like file I/O) require valid paths.
// These XYValid fields are shown by the `settings` command.
// They are set by the various validation functions.

type MinionConfig struct {
	// General settings
	ConfigFileName      string
	ConfigFileNameValid bool
	Language            language.Tag
	Verbose             bool
	SourceDir           string
	SourceDirValid      bool
	TargetDir           string
	TargetDirValid      bool
	Force               bool

	// Processing options
	Evenify       bool
	Merge         bool
	MergeFileName string

	// Page formatting
	RunningHeader   string
	ChapterPrefix   string
	Separator       string
	PagePrefix      string
	PageCountPrefix string
	BlankPageText   string

	// personal touch, adds funny logo to random pages
	// TODO
	PersonalTouch bool

	// Metadata to track which fields were explicitly set
	// This is used to determine which fields to merge
	// Note: keys are lowercase, should be converted with strings.ToLower()
	SetFields map[string]bool
}

// NewDefaultConfig creates a new configuration with default values,
// using the system language for texts
func NewDefaultConfig(systemLanguage language.Tag) *MinionConfig {
	log.Debug().Msg("Creating new default configuration")

	// Get corresponding texts for the provided language
	texts, exists := DefaultTexts[systemLanguage]
	if !exists {
		log.Debug().Str("language", systemLanguage.String()).Msg("Language not supported, falling back to English")
		systemLanguage = language.English
		texts = DefaultTexts[language.English]
	}

	return &MinionConfig{
		Verbose:       false,
		SourceDir:     DefaultSourceDir,
		TargetDir:     DefaultTargetDir,
		Force:         false,
		Evenify:       true,
		Merge:         false,
		MergeFileName: DefaultMergeFileName,

		ConfigFileName: MinionConfigFileName,
		Language:       systemLanguage,
		// Use language-specific texts
		ChapterPrefix:   texts.ChapterPrefix,
		RunningHeader:   texts.RunningHeader,
		PagePrefix:      texts.PageNumber,
		PageCountPrefix: DefaultPageCountPrefix,
		BlankPageText:   texts.BlankPageText,
		Separator:       DefaultSeparator,

		PersonalTouch: false,
		SetFields:     make(map[string]bool),
	}
}

// MergeWith merges the current config with another config, giving precedence to the other config
func (c *MinionConfig) MergeWith(other *MinionConfig) error {
	if other == nil {
		return nil
	}

	// Only override non-zero values
	if other.ConfigFileName != "" {
		c.ConfigFileName = other.ConfigFileName
	}
	if other.Language != language.Und {
		c.Language = other.Language
	}
	if other.SourceDir != "" {
		c.SourceDir = other.SourceDir
	}
	if other.TargetDir != "" {
		c.TargetDir = other.TargetDir
	}
	if other.MergeFileName != "" {
		c.MergeFileName = other.MergeFileName
	}
	if other.RunningHeader != "" {
		c.RunningHeader = other.RunningHeader
	}
	if other.ChapterPrefix != "" {
		c.ChapterPrefix = other.ChapterPrefix
	}
	if other.PagePrefix != "" {
		c.PagePrefix = other.PagePrefix
	}
	if other.BlankPageText != "" {
		c.BlankPageText = other.BlankPageText
	}
	if other.Separator != "" {
		c.Separator = other.Separator
	}

	// Boolean flags are only merged if they have been explicitly set.
	// See ADR-0009 on metadata.

	// only merge Verbose value if it has been explicitly set in "other"
	if other.SetFields["verbose"] {
		c.Verbose = other.Verbose
	}
	if other.SetFields["force"] {
		c.Force = other.Force
	}
	if other.SetFields["evenify"] {
		c.Evenify = other.Evenify
	}
	if other.SetFields["merge"] {
		c.Merge = other.Merge
	}
	if other.SetFields["personal"] {
		c.PersonalTouch = other.PersonalTouch
	}

	return nil
}
