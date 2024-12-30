package domain

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"os"
	"testing"
)

func TestGermanSystemDefaults(t *testing.T) {
	// Backup original environment
	origLang := os.Getenv("LANG")
	origLcAll := os.Getenv("LC_ALL")
	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("LC_ALL", origLcAll)
	}()

	// Set German system locale
	os.Setenv("LANG", "de_DE.UTF-8")
	os.Setenv("LC_ALL", "de_DE.UTF-8")

	// Get the system language (this should now detect German)
	systemLang := DetectSystemLanguage()

	// Create default config with detected language
	config := NewDefaultConfig(systemLang)

	// Verify that German was detected
	assert.Equal(t, language.German, config.Language, "System language should be German")

	// Verify that German texts were used
	expectedTexts := DefaultTexts[language.German]
	assert.Equal(t, expectedTexts.ChapterPrefix, config.ChapterPrefix, "Chapter prefix should be in German")
	assert.Equal(t, expectedTexts.RunningHeader, config.RunningHeader, "Running header should be in German")
	assert.Equal(t, expectedTexts.BlankPageText, config.BlankPageText, "Blank page text should be in German")
}

func TestUnsupportedSystemLanguageFallback(t *testing.T) {
	// Backup original environment
	origLang := os.Getenv("LANG")
	origLcAll := os.Getenv("LC_ALL")
	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("LC_ALL", origLcAll)
	}()

	// Set Spanish system locale (not in our supported languages)
	os.Setenv("LANG", "es_ES.UTF-8")
	os.Setenv("LC_ALL", "es_ES.UTF-8")

	// Get the system language
	systemLang := DetectSystemLanguage()

	// Create default config with detected language
	config := NewDefaultConfig(systemLang)

	// Verify fallback to English
	assert.Equal(t, language.English, config.Language, "System language should fall back to English")

	// Verify that English texts were used
	expectedTexts := DefaultTexts[language.English]
	assert.Equal(t, expectedTexts.ChapterPrefix, config.ChapterPrefix, "Chapter prefix should be in English")
	assert.Equal(t, expectedTexts.RunningHeader, config.RunningHeader, "Running header should be in English")
	assert.Equal(t, expectedTexts.BlankPageText, config.BlankPageText, "Blank page text should be in English")
}
