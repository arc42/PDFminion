package config

import (
	"os"

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
		RunE:  runPDFProcessing,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().Msg("PersistentPreRunE function called")

			// Skip config loading for basic commands (like version, help, credits, list-languages)
			if isImmediateCommand(cmd) {
				log.Debug().Msg("skipping configuration, basic command detected")
				return nil
			}

			// now cmd is either settings or the default pdf processing command
			// this PersistentPreRunE function needs to return error only if the configuration is invalid
			return nil
			//return loadConfiguration(cmd)
		},
	}
)

func loadConfiguration(cmd *cobra.Command) error {
	// Initialize and populate minionConfig
	config, err := LoadConfig()
	if err != nil {
		return err
	}
	minionConfig = config // Set package variable
	return nil
}

func SetupApplication(appVersion string) *cobra.Command {
	// Set application version
	domain.SetAppVersion(appVersion)

	// Setup flags, as some commands need the flags prior to execution
	setupFlags()

	// Setup commands for the root CLI application
	setupCommands()

	// Load configuration in layers (defaults, config-file(s), flags on CLI)
	minionConfig, err := LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Error loading configuration")
		os.Exit(1)
	}
	log.Debug().Interface("config:", minionConfig).Msg("Configuration completed ")

	return rootCmd
}

func setupFlags() {
	// Persistent flags (available to all commands)
	rootCmd.Flags().String("config", "", "Path to configuration file")
	rootCmd.Flags().StringP("language", "l", "", "Override system language")

	// Local flags (only for PDF processing)
	rootCmd.Flags().StringP("source", "s", domain.DefaultSourceDir, "Source directory for PDF files")
	rootCmd.Flags().StringP("target", "t", domain.DefaultTargetDir, "Target directory for processed files")
	rootCmd.Flags().BoolP("force", "f", false, "Force overwrite of target directory")
	rootCmd.Flags().BoolP("evenify", "e", true, "Ensure even page count in output")
	rootCmd.Flags().BoolP("verbose", "v", false, "Give more detailed output during processing")
	rootCmd.Flags().StringP("running-header", "r", "", "Text for running header")
	rootCmd.Flags().StringP("chapter-prefix", "c", domain.DefaultChapterPrefix, "Prefix for chapter numbers")
	rootCmd.Flags().StringP("page-prefix", "p", domain.DefaultPageNrPrefix, "Prefix for page numbers")
	rootCmd.Flags().StringP("blank-page-text", "b", domain.DefaultBlankPageText, "Text for blank pages")
	rootCmd.Flags().Bool("personal", false, "Adds a personal touch (aka logo) to random pages")
	rootCmd.Flags().String("merge", "merged.pdf", "--merge=filename, merge generated files into <filename>")
	rootCmd.Flags().String("separator", domain.DefaultSeparator, "Separator between chapter and page")
	rootCmd.Flags().String("page-count-prefix", domain.DefaultPageCountPrefix, "Prefix for total page count")

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
	)
	log.Debug().Msg("Add commands completed")

}

func isImmediateCommand(cmd *cobra.Command) bool {
	// immediateCommands don't require config validation, but will be executed immediately

	log.Debug().Str("", cmd.Name()).Msg("Check for immediate command")

	immediateCommands := []string{"version", "credits", "help", "list-languages"}
	for _, c := range immediateCommands {
		if cmd.Name() == c {
			log.Debug().Str("", c).Msg("identified as immediate command")

			return true
		}
	}
	return false
}

func ListLanguagesCmd() *cobra.Command {
	return &cobra.Command{
		Use:              "list-languages",
		Aliases:          []string{"ll", "list", "list-lang", "list-langs"},
		Short:            "Show available languages (like DE, EN, FR, etc.)",
		Long:             "Show the available languages for which PDFminion provides default settings (like EN, DE, FR).",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
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
			domain.PrintFinalConfiguration(*minionConfig)
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

func runPDFProcessing(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting PDF processing")
	// TODO: Implement actual PDF processing logic
	return nil
}
