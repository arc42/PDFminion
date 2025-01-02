package domain

import (
	"fmt"
	"github.com/Xuanwo/go-locale"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
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
		RunningHeader   string
		ChapterPrefix   string
		PageCountPrefix string
		PageNumber      string
		BlankPageText   string
	}{
		language.German: {
			RunningHeader:   "",
			ChapterPrefix:   "Kapitel",
			PageCountPrefix: "von",
			PageNumber:      "Seite",
			BlankPageText:   "Diese Seite bleibt absichtlich leer",
		},
		language.English: {
			RunningHeader:   "",
			ChapterPrefix:   DefaultChapterPrefix,
			PageCountPrefix: "of",
			PageNumber:      DefaultPageNrPrefix,
			BlankPageText:   DefaultBlankPageText,
		},
		language.French: {
			RunningHeader:   "",
			ChapterPrefix:   "Chapitre",
			PageCountPrefix: "sur",
			PageNumber:      "Page",
			BlankPageText:   "Cette page est intentionnellement laissée vide",
		},
	}
)

// GetMatcher returns the language matcher for supported languages
func GetMatcher() language.Matcher {
	return matcher
}

// MapSystemToAppLanguage detects the system language and maps to a supported language
// if the detected language is not supported, English is returned!
func MapSystemToAppLanguage() language.Tag {
	tag, err := locale.Detect()
	log.Debug().Str("found ", tag.String()).Msg("Detecting system language")
	if err == nil {
		// Check if detected language is supported
		tag, _, _ = matcher.Match(tag)
		if base, conf := tag.Base(); conf != language.No {
			log.Debug().Str("language", base.String()).Msg("Mapped detected language to English")
			return language.Make(base.String())
		}
	}

	// Default to English if all else fails
	log.Debug().Msg("Failed detecting system language, falling back to English")
	return language.English
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

// ParseLanguageCode takes a language string (e.g., "en", "en-US", "de-DE")
// and returns the corresponding base language Tag.
// Returns English tag if the input is invalid or unsupported.
func ParseLanguageCode(langStr string) language.Tag {
	if langStr == "" {
		return language.English
	}

	parsedLang, err := ValidateLanguage(langStr)
	if err != nil {
		log.Debug().Str("input", langStr).Msg("Invalid language code, falling back to English")
		return language.English
	}

	// Extract base language (e.g., "en" from "en-US")
	baseLang, conf := parsedLang.Base()
	if conf == language.No {
		log.Debug().Str("input", langStr).Msg("Could not determine base language, falling back to English")
		return language.English
	}

	return language.Make(baseLang.String())
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
	currentLanguage := MapSystemToAppLanguage()
	currentCode, currentNameInOriginal, currentNameInEnglish := GetLanguageName(currentLanguage)

	fmt.Println("Supported Languages:")
	for _, lang := range ListAvailableLanguages() {
		fmt.Printf("Code %s (%s, %s)\n",
			lang[0], lang[1], lang[2])
	}

	fmt.Printf("Current system language: %s (%s, %s)\n",
		currentCode, currentNameInOriginal, currentNameInEnglish)
}
