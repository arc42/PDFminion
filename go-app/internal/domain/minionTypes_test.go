package domain

import (
	"golang.org/x/text/language"
	"testing"
)

func TestMinionConfig_MergeWith(t *testing.T) {
	t.Run("merges with partial config", func(t *testing.T) {
		// Setup base config
		base := &MinionConfig{
			Language:      language.English,
			SourceDir:     "/original/source",
			TargetDir:     "/original/target",
			MergeFileName: "original.pdf",
			Verbose:       false,
			Force:         false,
			Evenify:       true,
		}

		// Setup config to merge with
		other := &MinionConfig{
			SourceDir:     "/new/source",
			MergeFileName: "new.pdf",
			Force:         true,
			// Other fields left at zero values
		}

		// Perform merge
		err := base.MergeWith(other)

		// Assert
		if err != nil {
			t.Errorf("MergeWith returned unexpected error: %v", err)
		}

		// Check that specified fields were overridden
		if base.SourceDir != "/new/source" {
			t.Errorf("SourceDir not merged correctly, got %s, want %s",
				base.SourceDir, "/new/source")
		}
		if base.MergeFileName != "new.pdf" {
			t.Errorf("MergeFileName not merged correctly, got %s, want %s",
				base.MergeFileName, "new.pdf")
		}
		if !base.Force {
			t.Error("Force flag not merged correctly, expected true")
		}

		// Check that unspecified fields retained original values
		if base.Language != language.English {
			t.Errorf("Language was modified unexpectedly, got %v, want %v",
				base.Language, language.English)
		}
		if base.TargetDir != "/original/target" {
			t.Errorf("TargetDir was modified unexpectedly, got %s, want %s",
				base.TargetDir, "/original/target")
		}
		if !base.Evenify {
			t.Error("Evenify was modified unexpectedly, expected true")
		}
	})

	t.Run("merges with nil config", func(t *testing.T) {
		// Setup base config
		base := &MinionConfig{
			Language:      language.English,
			SourceDir:     "/original/source",
			TargetDir:     "/original/target",
			MergeFileName: "original.pdf",
			Verbose:       false,
			Force:         true,
		}

		// Store original values for comparison
		originalConfig := *base

		// Perform merge with nil
		err := base.MergeWith(nil)

		// Assert
		if err != nil {
			t.Errorf("MergeWith(nil) returned unexpected error: %v", err)
		}

		// Check that no fields were modified
		if *base != originalConfig {
			t.Error("Config was modified when merging with nil")
		}
	})
}
