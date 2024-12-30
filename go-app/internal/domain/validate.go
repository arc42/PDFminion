package domain

import (
	"fmt"
	"golang.org/x/text/language"
	"io"
	"os"
)

// ValidateConfig checks the configuration for correctness
func ValidateConfig(config *MinionConfig) error {
	return config.Validate()
}

// Validate validates the configuration for semantic correctness:
// - Source directory must exist
// - Target directory shall be empty if it exists (unless --force is set), otherwise it will be created
// - a given language must be available in the provided languages
// - and more
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
