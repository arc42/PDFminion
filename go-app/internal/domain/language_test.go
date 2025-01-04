package domain

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"os"
	"testing"
)

func TestDetectSystemLanguageOnGernotsComputer(t *testing.T) {
	// this test depends on a specific (German) envirionment, therefore we skip it in short mode (CI)

	t.Skip("Skipping local-only test in short mode")

	// Get the system language (this should now detect German)
	systemLang := MapSystemToAppLanguage()

	// Verify that German was detected
	assert.Equal(t, language.German, systemLang, "System language should be German on Gernots` computer")
}

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
	systemLang := MapSystemToAppLanguage()

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

func TestFrenchSystemDefaults(t *testing.T) {
	// Backup original environment
	origLang := os.Getenv("LANG")
	origLcAll := os.Getenv("LC_ALL")
	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("LC_ALL", origLcAll)
	}()

	// Set French system locale
	os.Setenv("LANG", "fr_FR.UTF-8")
	os.Setenv("LC_ALL", "fr_FR.UTF-8")

	// Get the system language (this should now detect German)
	systemLang := MapSystemToAppLanguage()

	// Create default config with detected language
	config := NewDefaultConfig(systemLang)

	// Verify that German was detected
	assert.Equal(t, language.French, config.Language, "System language should be German")

	// Verify that German texts were used
	expectedTexts := DefaultTexts[language.French]
	assert.Equal(t, expectedTexts.ChapterPrefix, config.ChapterPrefix, "Chapter prefix should be in French")
	assert.Equal(t, expectedTexts.RunningHeader, config.RunningHeader, "Running header should be in French")
	assert.Equal(t, expectedTexts.BlankPageText, config.BlankPageText, "Blank page text should be in French")
}

func TestUncommonLanguageWillMapToEnglish(t *testing.T) {
	// Backup original environment
	origLang := os.Getenv("LANG")
	origLcAll := os.Getenv("LC_ALL")
	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("LC_ALL", origLcAll)
	}()

	// Set Icelandic system locale (not in our supported languages)
	err := os.Setenv("LANG", "is_IS.UTF-8")
	if err != nil {
		log.Error().Err(err).Msg("Error setting LANG to Icelandic")
	}
	err = os.Setenv("LC_ALL", "is_IS.UTF-8")
	if err != nil {
		log.Error().Err(err).Msg("Error setting LC_AL to Icelandic")
	}
	log.Debug().Msg("Set Icelandic system locale with LANG and LC_ALL")

	systemLang := MapSystemToAppLanguage()
	assert.Equal(t, language.English, systemLang, "App language should be English, but is %v", systemLang)
}

func TestUnsupportedSystemLanguageFallback(t *testing.T) {
	// Backup original environment
	origLang := os.Getenv("LANG")
	origLcAll := os.Getenv("LC_ALL")
	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("LC_ALL", origLcAll)
	}()

	// Set Icelandic system locale (not in our supported languages)
	err := os.Setenv("LANG", "is_IS.UTF-8")
	if err != nil {
		log.Error().Err(err).Msg("Error setting LANG to Icelandig")
	}
	err = os.Setenv("LC_ALL", "is_IS.UTF-8")
	if err != nil {
		log.Error().Err(err).Msg("Error setting LC_AL to Icelandic")
	}

	// Get the system language
	systemLang := MapSystemToAppLanguage()

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

func TestIsLanguageSupported(t *testing.T) {
	tests := []struct {
		lang language.Tag
		want bool
	}{
		{language.German, true},
		{language.English, true},
		{language.French, true},
		{language.Icelandic, false},
		{language.Zulu, false},
	}

	for _, tt := range tests {
		got := IsLanguageSupported(tt.lang)
		if got != tt.want {
			t.Errorf("IsLanguageSupported(%v) = %v, want %v", tt.lang, got, tt.want)
		}
	}
}

func TestParseLanguageCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected language.Tag
	}{
		{"Empty string", "", language.English},
		{"Simple English", "en", language.English},
		{"Complex English", "en-US", language.English},
		{"Simple German", "de", language.German},
		{"Complex German", "de-DE", language.German},
		{"Simple French", "fr", language.French},
		{"Complex French", "fr-FR", language.French},
		{"Unsupported language", "is_IS", language.English},
		{"Invalid code", "xx", language.English},
		{"Garbage input", "notALanguage", language.English},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLanguageCode(tt.input)
			assert.Equal(t, tt.expected, result,
				"ParseLanguageCode(%q) = %v, expected %v",
				tt.input, result, tt.expected)
		})
	}
}
