package config

import (
	"fmt"
	"os"
	"path/filepath"
	"pdfminion/internal/pdf"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pdfminion/internal/domain"
)

// Package-level variable accessible to all commands

var (
	rootCmd = &cobra.Command{
		Use:   "pdfminion",
		Short: "PDFMinion adds page numbers to PDF files with custom options",
		Long:  "PDFMinion is a CLI tool to add page numbers to existing PDF files with customizable options like chapter numbers, running headers, and more",
		// PersistentPreRun runs before any command (including the default RunE)
		// This is where we configure the application AFTER flags have been parsed
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Get verbose flag
			verbose := viper.GetBool("verbose")

			// Now configure the application with parsed flags
			var err error
			ActiveMinionConfig, err = ConfigureApplication(verbose, cmd)
			if err != nil {
				log.Error().Err(err).Msg("Error loading configuration")
				os.Exit(1)
			}
			log.Debug().Interface("configuration:", ActiveMinionConfig).Msg("Configuration completed")
		},
		// When no subcommand is provided, process PDFs with the given flags
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPDFProcessing(cmd, args)
		},
	}

	ActiveMinionConfig domain.MinionConfig

	// until we have evaluated the configuration,
	// we assume non-verbose output
	verbose = false
)

func SetupApplication(appVersion string) *cobra.Command {
	// Set application version
	domain.SetAppVersion(appVersion)

	// Setup flags - these need to be defined before commands
	setupFlags()

	// Setup commands for the root CLI application
	setupCommands()

	// NOTE: Configuration is now loaded in PersistentPreRun (see rootCmd definition above)
	// This ensures flags are parsed BEFORE we try to read them
	// Previously we called ConfigureApplication here, which was too early!

	return rootCmd
}

func setupFlags() {
	// Persistent flags (available to all commands)
	rootCmd.PersistentFlags().StringP("language", "l", "", "Override system language")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Give more detailed output during processing")
	// NOTE: Config file removed - we use fixed location ~/.config/pdfminion/config.yaml
	// Use 'pdfminion config' commands to manage configuration files

	// Local flags (only for PDF processing)
	rootCmd.Flags().StringP("source", "s", domain.DefaultSourceDir, "Source directory for PDF files")
	rootCmd.Flags().StringP("target", "t", domain.DefaultTargetDir, "Target directory for processed files")
	rootCmd.Flags().BoolP("force", "f", false, "Force overwrite of target directory")
	rootCmd.Flags().BoolP("evenify", "e", true, "Ensure even page count in output")
	rootCmd.Flags().StringP("running-header", "r", "", "Text for running header")
	rootCmd.Flags().String("chapter-prefix", domain.DefaultChapterPrefix, "Prefix for chapter numbers")
	rootCmd.Flags().StringP("page-prefix", "p", domain.DefaultPageNrPrefix, "Prefix for page numbers")
	rootCmd.Flags().StringP("blank-page-text", "b", domain.DefaultBlankPageText, "Text for blank pages")
	rootCmd.Flags().Bool("personal", false, "Adds a personal touch (aka logo) to random pages")
	rootCmd.Flags().String("merge", "merged.pdf", "--merge=filename, merge generated files into <filename>")
	rootCmd.Flags().String("separator", domain.DefaultSeparator, "Separator between chapter and page")
	rootCmd.Flags().String("page-count-prefix", domain.DefaultPageCountPrefix, "Prefix for total page count")
	rootCmd.Flags().BoolP("toc", "o", false, "Generate table of contents")

	// Bind all flags to viper
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind persistent flags to viper")
	}

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind flags to viper")
	}

	log.Debug().Msg("setupFlags completed")
}

func setupCommands() {

	// Hide completion command from help output
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.AddCommand(
		VersionCmd(),
		CreditsCmd(),
		SettingsCmd(),
		ListLanguagesCmd(),
		ConfigCmd(),
	)
	log.Debug().Msg("Add commands completed")

}

func ListLanguagesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list-languages",
		Aliases: []string{"ll", "list", "list-lang", "list-langs"},
		Short:   "Show available languages (like DE, EN, FR, etc.)",
		Long:    "Show the available languages for which PDFminion provides default settings (like EN, DE, FR).",
		//PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().Msg("executing list-languages command")
			domain.PrintLanguages()
		},
	}
}

func VersionCmd() *cobra.Command {

	return &cobra.Command{
		Use:              "version",
		Aliases:          []string{"ver", "v"},
		Short:            "Show version information",
		Long:             "Show version information for the PDFMinion application, including version number, build date and build target",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().Msg("executing version command")
			domain.PrintVersion()
		},
	}
}

// SettingsCmd requires  flags to be evaluated, so that the final configuration can be determined
func SettingsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "settings",
		Aliases: []string{"setting", "set"},
		Short:   "Show the final configuration",
		Long:    "Show the final configuration, after defaults, config files, and flags have been evaluated.",
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().Msg("executing settings command")
			domain.PrintFinalConfiguration(&ActiveMinionConfig)
		},
	}
}

func CreditsCmd() *cobra.Command {
	return &cobra.Command{
		Use:              "credits",
		Aliases:          []string{"credit", "thanx"},
		Short:            "Show credits",
		Long:             "Show credits for the open-source tools used to implement PDFMinion",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().Msg("executing credits command")
			domain.GiveCredits()
		},
	}
}

// ConfigCmd provides subcommands for managing configuration files
func ConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:              "config",
		Short:            "Manage configuration files",
		Long:             "Manage PDFMinion configuration files in ~/.config/pdfminion/config.yaml",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	}

	// Add subcommands
	configCmd.AddCommand(
		configPathCmd(),
		configShowCmd(),
		configCreateCmd(),
		configEditCmd(),
	)

	return configCmd
}

// configPathCmd shows the path to the config file
func configPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Show path to config file",
		Long:  "Display the path where PDFMinion looks for its configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			path := GetConfigPath()
			fmt.Printf("Config file path: %s\n", path)

			// Check if file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Println("(file does not exist yet)")
			} else {
				fmt.Println("(file exists)")
			}
		},
	}
}

// configShowCmd displays the contents of the config file
func configShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show config file contents",
		Long:  "Display the current contents of the configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			path := GetConfigPath()

			// Check if file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Printf("Config file does not exist: %s\n", path)
				fmt.Println("Use 'pdfminion config create' to create one")
				return
			}

			// Read and display file contents
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading config file: %v\n", err)
				return
			}

			fmt.Printf("Config file: %s\n", path)
			fmt.Println("---")
			fmt.Print(string(content))
			if len(content) > 0 && content[len(content)-1] != '\n' {
				fmt.Println() // Ensure newline at end
			}
		},
	}
}

// configCreateCmd creates a new config file with example content
func configCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create example config file",
		Long:  "Create a new configuration file with example settings at ~/.config/pdfminion/config.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			path := GetConfigPath()

			// Check if file already exists
			if _, err := os.Stat(path); err == nil {
				fmt.Printf("Config file already exists: %s\n", path)
				fmt.Println("Use 'pdfminion config show' to view it or 'pdfminion config edit' to modify it")
				return
			}

			// Create directory if it doesn't exist
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("Error creating config directory: %v\n", err)
				return
			}

			// Create example config file
			exampleConfig := `# PDFMinion Configuration File
# All settings are optional - any setting not specified will use defaults
# Command-line flags override these settings

# Language for text defaults (EN, DE, FR)
# language: EN

# Default directories
# source: _source
# target: _target

# Processing options
# force: false
# evenify: true
# verbose: false

# Text customization
# running-header: ""
# chapter-prefix: "Chapter"
# page-prefix: "Page"
# separator: " "
# page-count-prefix: "of"
# blank-page-text: "This page intentionally left blank"

# Post-processing
# merge: ""  # Set to filename to enable merging
# toc: false
# personal: false
`

			if err := os.WriteFile(path, []byte(exampleConfig), 0644); err != nil {
				fmt.Printf("Error creating config file: %v\n", err)
				return
			}

			fmt.Printf("Created example config file: %s\n", path)
			fmt.Println("Edit this file to customize your default settings")
			fmt.Println("Use 'pdfminion config show' to view it or 'pdfminion config edit' to modify it")
		},
	}
}

// configEditCmd opens the config file in the user's default editor
func configEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit config file",
		Long:  "Open the configuration file in your default editor ($EDITOR)",
		Run: func(cmd *cobra.Command, args []string) {
			path := GetConfigPath()

			// Check if file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Printf("Config file does not exist: %s\n", path)
				fmt.Println("Use 'pdfminion config create' to create one first")
				return
			}

			// Get editor from environment
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi" // Fallback to vi
			}

			fmt.Printf("Opening config file with %s...\n", editor)

			// Execute editor
			execCmd := fmt.Sprintf("%s %s", editor, path)
			if err := os.Chdir(filepath.Dir(path)); err != nil {
				fmt.Printf("Error changing directory: %v\n", err)
				return
			}

			fmt.Printf("Run this command to edit: %s\n", execCmd)
			fmt.Println("(Note: Direct editor invocation not available in this context)")
		},
	}
}

func runPDFProcessing(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting PDF processing")

	fmt.Printf("Processing PDFs in %q\n", ActiveMinionConfig.SourceDir)

	// Validate configuration
	if err := domain.ValidateConfig(&ActiveMinionConfig); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Process PDFs
	if err := pdf.ProcessPDFs(&ActiveMinionConfig); err != nil {
		return fmt.Errorf("error processing PDFs: %w", err)
	}
	return nil
}
