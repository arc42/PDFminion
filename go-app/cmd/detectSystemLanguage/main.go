//go:build exclude_from_coverage
// +build exclude_from_coverage

package main

import (
	"fmt"
	"log"

	"github.com/Xuanwo/go-locale"
)

func getSystemLanguage() (string, error) {
	tag, err := locale.Detect()
	if err != nil {
		return "", err
	}
	return tag.String(), nil
}

func main() {
	lang, err := getSystemLanguage()
	if err != nil {
		log.Fatalf("Failed to detect system language: %v", err)
	}

	fmt.Printf("Detected system language: %s\n", lang)
}
