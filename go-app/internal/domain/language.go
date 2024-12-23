package domain

import (
	"github.com/Xuanwo/go-locale"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"os"
	"strings"
)

var (
	supportedLanguages = []language.Tag{
		language.German,
		language.English,
		language.French,
	}

	matcher = language.NewMatcher(supportedLanguages)

	// DefaultTexts holds localized UI texts
	DefaultTexts = map[language.Tag]struct {
		ChapterPrefix string
		RunningHeader string
		PageFooter    string
		PageNumber    string
		BlankPageText string
	}{
		language.German: {
			ChapterPrefix: "Kapitel",
			RunningHeader: "Seite",
			PageFooter:    "Seite %d von %d",
			PageNumber:    "Seite %d",
			BlankPageText: "Diese Seite bleibt absichtlich leer",
		},
		language.English: {
			ChapterPrefix: "Chapter",
			RunningHeader: "Page",
			PageFooter:    "Page %d of %d",
			PageNumber:    "Page %d",
			BlankPageText: "deliberately left blank",
		},
		language.French: {
			ChapterPrefix: "Chapitre",
			RunningHeader: "Page",
			PageFooter:    "Page %d sur %d",
			PageNumber:    "Page %d",
			BlankPageText: "Cette page est intentionnellement laissée vide",
		},
	}
)

// DetectSystemLanguage detects the system language
func DetectSystemLanguage() language.Tag {
	// Check environment variables first
	for _, envVar := range []string{"LANG", "LC_ALL", "LC_MESSAGES"} {
		if tag := getLanguageFromEnv(envVar); tag != language.Und {
			return tag
		}
	}

	// Fallback to go-locale
	tag, err := locale.Detect()
	if err == nil {
		tag, _, _ = matcher.Match(tag)
		if base, conf := tag.Base(); conf != language.No {
			return language.Make(base.String())
		}
	}

	// Default to English if all else fails
	return language.English
}

// getLanguageFromEnv attempts to extract a language tag from an environment variable
func getLanguageFromEnv(envVar string) language.Tag {
	if sysLocale := os.Getenv(envVar); sysLocale != "" {
		// Strip charset if present (e.g., "en_US.UTF-8" -> "en_US")
		langPart := strings.Split(sysLocale, ".")[0]
		if tag, err := language.Parse(langPart); err == nil {
			tag, _, _ = matcher.Match(tag)
			if base, conf := tag.Base(); conf != language.No {
				return language.Make(base.String())
			}
		}
	}
	return language.Und
}

// ValidateLanguage checks if a language string is valid and supported
func ValidateLanguage(lang string) (language.Tag, error) {
	tag, err := language.Parse(lang)
	if err != nil {
		return language.Und, err
	}

	tag, _, _ = matcher.Match(tag)
	return tag, nil
}

// GetLanguageName returns the name of the language in its native form and English
func GetLanguageName(tag language.Tag) (code, nameInOriginal, nameInEnglish string) {
	base, _ := tag.Base()
	code = base.String()
	nameInOriginal = display.Self.Name(tag)
	nameInEnglish = display.English.Tags().Name(tag)
	return
}

// ListAvailableLanguages returns a list of supported languages with their codes and names
func ListAvailableLanguages() [][]string {
	var languages [][]string
	for _, tag := range supportedLanguages {
		code, nameInOriginal, nameInEnglish := GetLanguageName(tag)
		languages = append(languages, []string{code, nameInOriginal, nameInEnglish})
	}
	return languages
}

// PrintLanguages displays the list of supported languages and the current system language
func PrintLanguages() {
	currentLanguage := DetectSystemLanguage()
	currentCode, currentNameInOriginal, currentNameInEnglish := GetLanguageName(currentLanguage)

	log.Info().Msg("Supported Languages:")
	for _, lang := range ListAvailableLanguages() {
		log.Info().
			Str("code", lang[0]).
			Str("name", lang[1]).
			Str("english_name", lang[2]).
			Msg("Language")
	}

	log.Info().
		Str("code", currentCode).
		Str("native_name", currentNameInOriginal).
		Str("english_name", currentNameInEnglish).
		Msg("Current system language")
}
