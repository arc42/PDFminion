package domain

import (
	"fmt"
	"golang.org/x/text/language"
	"io"
	"os"
)

const (
	MinionConfigFileName = "pdfminion.yaml"

	// Default formatting settings
	DefaultSourceDir     = "_pdfs"
	DefaultTargetDir     = "_target"
	PageNrPrefix         = ""
	ChapterPrefix        = "Chapter"
	ChapterPageSeparator = " - "
)

type MinionConfig struct {
	// General settings
	ConfigFileName string
	Language       language.Tag
	Debug          bool
	SourceDir      string
	TargetDir      string
	Force          bool

	// Processing options
	Evenify       bool
	Merge         bool
	MergeFileName string

	// Page formatting
	RunningHeader        string
	ChapterPrefix        string
	Separator            string
	PagePrefix           string
	TotalPageCountPrefix string
	BlankPageText        string
}

// NewDefaultConfig creates a new configuration with default values
func NewDefaultConfig() *MinionConfig {
	return &MinionConfig{
		ConfigFileName: MinionConfigFileName,
		Language:       language.English,
		Debug:          false,
		SourceDir:      "_pdfs",
		TargetDir:      "_target",
		Force:          false,
		Evenify:        true,
		Merge:          false,
		MergeFileName:  "merged.pdf",
		Separator:      " - ",
		PagePrefix:     "",
		BlankPageText:  "",
	}
}

// MergeWith merges the current config with another config, giving precedence to the other config
func (c *MinionConfig) MergeWith(other *MinionConfig) error {
	if other == nil {
		return nil
	}

	// Only override non-zero values
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

	// Boolean flags always override
	c.Debug = other.Debug
	c.Force = other.Force
	c.Evenify = other.Evenify
	c.Merge = other.Merge

	return nil
}

func (c *MinionConfig) Validate() error {
	// Validate source directory
	if err := c.validateSourceDir(); err != nil {
		return err
	}

	// Validate target directory
	if err := c.validateTargetDir(); err != nil {
		return err
	}

	// Validate language
	if c.Language == language.Und {
		return fmt.Errorf("invalid or undefined language")
	}

	return nil
}

func (c *MinionConfig) validateSourceDir() error {
	if _, err := os.Stat(c.SourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory %q does not exist", c.SourceDir)
	}
	return nil
}

func (c *MinionConfig) validateTargetDir() error {
	if _, err := os.Stat(c.TargetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(c.TargetDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create target directory %q: %w", c.TargetDir, err)
		}
		return nil
	}

	if !c.Force {
		empty, err := isDirEmpty(c.TargetDir)
		if err != nil {
			return err
		}
		if !empty {
			return fmt.Errorf("target directory %q is not empty (use --force to override)", c.TargetDir)
		}
	}

	return nil
}

func isDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Try to read one entry
	names, err := f.Readdirnames(1)
	if err != nil && err != io.EOF {
		return false, err // Return error if it's not EOF
	}

	// Directory is empty if we got EOF (no entries)
	return len(names) == 0, nil
}
