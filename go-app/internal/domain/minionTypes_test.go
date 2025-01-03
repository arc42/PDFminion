package domain

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestDefaultConfigIsEnglish(t *testing.T) {
	base := NewDefaultConfig(language.English)

	// the language specific values should be set to the English values
	assert.Equal(t, DefaultPageNrPrefix, base.PagePrefix)
	assert.Equal(t, DefaultChapterPrefix, base.ChapterPrefix)
	assert.Equal(t, DefaultPageCountPrefix, base.PageCountPrefix)
	assert.Equal(t, DefaultBlankPageText, base.BlankPageText)
}

func TestDefaultConfigLanguageGerman(t *testing.T) {
	base := NewDefaultConfig(language.German)

	// the language specific values should be set to the German values
	assert.Equal(t, DefaultTexts[language.German].PageNumber, base.PagePrefix)
	assert.Equal(t, DefaultTexts[language.German].ChapterPrefix, base.ChapterPrefix)
	assert.Equal(t, DefaultTexts[language.German].PageCountPrefix, base.PageCountPrefix)
	assert.Equal(t, DefaultTexts[language.German].BlankPageText, base.BlankPageText)
}

func TestDefaultConfigLanguageFrench(t *testing.T) {
	base := NewDefaultConfig(language.French)

	assert.Equal(t, DefaultTexts[language.French].PageNumber, base.PagePrefix)
	assert.Equal(t, DefaultTexts[language.French].ChapterPrefix, base.ChapterPrefix)
	assert.Equal(t, DefaultTexts[language.French].PageCountPrefix, base.PageCountPrefix)
	assert.Equal(t, DefaultTexts[language.French].BlankPageText, base.BlankPageText)
}

func TestDefaultConfigLanguageUnknown(t *testing.T) {
	// Set up a default config with an unknown language "Zulu"
	base := NewDefaultConfig(language.Zulu)

	// the language specific values should be set to the English values
	assert.Equal(t, DefaultTexts[language.English].PageNumber, base.PagePrefix)
	assert.Equal(t, DefaultTexts[language.English].ChapterPrefix, base.ChapterPrefix)
	assert.Equal(t, DefaultTexts[language.English].PageCountPrefix, base.PageCountPrefix)
	assert.Equal(t, DefaultTexts[language.English].BlankPageText, base.BlankPageText)
}

// TestDefaultConfigMergeWithFlagEnglish merges an English default config with
// --language German. The language specific values should be set to the German (default)
func TestDefaultConfigMergeWithFlagGerman(t *testing.T) {
	base := NewDefaultConfig(language.English)

	german := &MinionConfig{
		Language:  language.German,
		SetFields: map[string]bool{"language": true},
	}

	// can merge without error
	assert.NoError(t, base.MergeWith(german), "MergeWith should not return an error")

	// the language specific values should be set to the German values
	// the language specific values should be set to the German values
	assert.Equal(t, DefaultTexts[language.German].PageNumber, base.PagePrefix)
	assert.Equal(t, DefaultTexts[language.German].ChapterPrefix, base.ChapterPrefix)
	assert.Equal(t, DefaultTexts[language.German].PageCountPrefix, base.PageCountPrefix)
	assert.Equal(t, DefaultTexts[language.German].BlankPageText, base.BlankPageText)
}

// TestDefaultConfigMergeWithFlagEnglish merges an English default config with
// --language French. The language specific values should be set to the German (default)
func TestDefaultConfigMergeWithFlagFrench(t *testing.T) {
	base := NewDefaultConfig(language.English)

	french := &MinionConfig{
		Language:  language.French,
		SetFields: map[string]bool{"language": true},
	}

	// can merge without error
	assert.NoError(t, base.MergeWith(french), "MergeWith should not return an error")

	// the language specific values should be set to the German values
	// the language specific values should be set to the German values
	assert.Equal(t, DefaultTexts[language.French].PageNumber, base.PagePrefix)
	assert.Equal(t, DefaultTexts[language.French].ChapterPrefix, base.ChapterPrefix)
	assert.Equal(t, DefaultTexts[language.French].PageCountPrefix, base.PageCountPrefix)
	assert.Equal(t, DefaultTexts[language.French].BlankPageText, base.BlankPageText)
}

func TestMinionConfig_MergeWithMinimal(t *testing.T) {
	minimalBaseConfig := &MinionConfig{SourceDir: "/original/source"}
	other := &MinionConfig{
		SourceDir: "/new/source",
		SetFields: map[string]bool{"source": true},
	}
	assert.NoError(t, minimalBaseConfig.MergeWith(other), "MergeWith should not return an error")
	assert.Equal(t, "/new/source", minimalBaseConfig.SourceDir)
}

func TestMinionConfig_MergeMiniWithEmpty(t *testing.T) {
	minimalBaseConfig := &MinionConfig{
		SourceDir: "/original/source",
		TargetDir: "/original/target",
		Verbose:   true,
		Force:     false,
	}
	other := &MinionConfig{
		SetFields: map[string]bool{},
	}

	assert.NoError(t, minimalBaseConfig.MergeWith(other), "MergeWith should not return an error")
	assert.Equal(t, "/original/source", minimalBaseConfig.SourceDir)
	assert.Equal(t, "/original/target", minimalBaseConfig.TargetDir)
	assert.Equal(t, true, minimalBaseConfig.Verbose)
	assert.Equal(t, false, minimalBaseConfig.Force)
}

// TestMinionConfig_MergeWithEmpty tests that merging booleans values
// overwritten in the other config.
func TestMinionConfig_MergeBooleanOverwrite(t *testing.T) {
	minimalBaseConfig := &MinionConfig{
		Verbose: true,
		Force:   false,
	}
	other := &MinionConfig{
		Verbose:   false,
		Force:     true,
		SetFields: map[string]bool{"verbose": true, "force": true},
	}

	assert.NoError(t, minimalBaseConfig.MergeWith(other), "MergeWith should not return an error")
	assert.Equal(t, false, minimalBaseConfig.Verbose)
	assert.Equal(t, true, minimalBaseConfig.Force)
}

// TestMinionConfig_MergeWithPartialSuperset: A few fields are overwritten in the other config, one field (merge) was unset in base and is set in other.
// One boolean field in other overwrites the value in base.
func TestMinionConfig_MergeWithPartialSuperset(t *testing.T) {
	// Setup base config
	base := &MinionConfig{
		Language:  language.English,
		SourceDir: "/original/source",
		TargetDir: "/original/target",
		Verbose:   false,
		Force:     false,
		Evenify:   false,
	}

	// Setup config to merge with
	other := &MinionConfig{
		Language:      language.German,
		SourceDir:     "/new/source",
		Verbose:       true,
		Merge:         true,
		MergeFileName: "new.pdf",
		SetFields: map[string]bool{
			"source":        true,
			"language":      true,
			"verbose":       true,
			"merge":         true,
			"mergeFileName": true,
		},
	}

	assert.NoError(t, base.MergeWith(other), "MergeWith should not return an error")

	assert.Equal(t, language.German, base.Language, "Language should have been overwritten and set to German")
	assert.Equal(t, true, base.Verbose, "Verbose should have been overwritten and set to true")
	assert.Equal(t, "/new/source", base.SourceDir, "SourceDir should have been overwritten and set to /new/source")
	assert.Equal(t, true, base.Merge, "Merge should have been set to true")
	assert.Equal(t, "new.pdf", base.MergeFileName, "MergeFileName should have been set to new.pdf")
}

// TestMinionConfig_MergeWithSubset: A few fields are overwritten in the other config,
func TestMinionConfig_MergeWithSubset(t *testing.T) {
	// Setup base config
	base := &MinionConfig{
		Language:  language.English,
		SourceDir: "/original/source",
		TargetDir: "/original/target",
		Verbose:   false,
		Force:     false,
		Evenify:   false,
	}

	// Setup config to merge with
	other := &MinionConfig{
		SourceDir: "/new/source",
		Verbose:   true,
		SetFields: map[string]bool{
			"source":  true,
			"verbose": true,
		},
	}

	assert.NoError(t, base.MergeWith(other), "MergeWith should not return an error")

	assert.Equal(t, true, base.Verbose, "Verbose should have been overwritten and set to true")
	assert.Equal(t, "/new/source", base.SourceDir, "SourceDir should have been overwritten and set to /new/source")

}

// TestMinionConfig_MergeWithSubset: A few fields are overwritten in the other config,
func TestMinionConfig_MergeWithKnownOtherLanguage(t *testing.T) {
	// Setup base config
	base := NewDefaultConfig(language.English)

	// Setup config to merge with
	other := &MinionConfig{
		SetFields: map[string]bool{
			"language": true,
			"verbose":  true,
		},
	}

	assert.NoError(t, base.MergeWith(other), "MergeWith should not return an error")

}
