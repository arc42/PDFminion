package config

import (
	"fmt"
	"os"
	"strings"

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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		// Add RunE to the root command to handle PDF processing by default
		RunE: runPDFProcessing,
	}

	ActiveMinionConfig domain.MinionConfig

	// until we have evaluated the configuration,
	// we assume non-verbose output
	verbose = false
)

func SetupApplication(appVersion string) *cobra.Command {
	// Set application version
	domain.SetAppVersion(appVersion)

	// Add version flag to persistent flags (must be done before setupFlags)
	rootCmd.PersistentFlags().BoolP("version", "", false, "Print version information and exit")

	// Add a pre-run hook that will be executed before any command
	originalPreRunE := rootCmd.PersistentPreRunE
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Check for version flag first
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			domain.PrintVersion()
			os.Exit(0)
		}

		// Call original pre-run if it exists
		if originalPreRunE != nil {
			return originalPreRunE(cmd, args)
		}
		return nil
	}

	// Setup flags, as some commands need the flags prior to execution
	setupFlags()

	// check if verbose is set
	verbose = viper.GetBool("verbose")
	log.Debug().Bool("verbose", verbose).Msg("flag ")

	if verbose {
		fmt.Printf("verbose mode requested.")
	}
	// Setup commands for the root CLI application
	setupCommands()

	// Create a flag checker for configuration
	flagChecker := NewCobraFlagChecker(rootCmd)

	// Now that commands are set up, we can configure the application
	// using our layered approach, see ADR-0008
	var err error
	ActiveMinionConfig, err = ConfigureApplication(verbose, flagChecker)
	if err != nil {
		log.Error().Err(err).Msg("Error loading configuration")
		os.Exit(1)
	}
	log.Debug().Interface("configuration:", ActiveMinionConfig).Msg("Configuration completed ")

	return rootCmd
}

func setupFlags() {
	// Persistent flags (available to all commands)
	rootCmd.PersistentFlags().StringP("language", "l", "", "Override system language")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Give more detailed output during processing")

	// Local flags (only for PDF processing)
	rootCmd.Flags().StringP("source", "s", domain.DefaultSourceDir, "Source directory for PDF files")
	rootCmd.Flags().StringP("target", "t", domain.DefaultTargetDir, "Target directory for processed files")
	rootCmd.Flags().BoolP("force", "f", false, "Force overwrite of target directory")
	rootCmd.Flags().BoolP("evenify", "e", true, "Ensure even page count in output")
	rootCmd.Flags().StringP("running-header", "r", "", "Text for running header")
	rootCmd.Flags().StringP("chapter-prefix", "c", domain.DefaultChapterPrefix, "Prefix for chapter numbers")
	rootCmd.Flags().StringP("page-prefix", "p", domain.DefaultPageNrPrefix, "Prefix for page numbers")
	rootCmd.Flags().StringP("blank-page-text", "b", domain.DefaultBlankPageText, "Text for blank pages")
	rootCmd.Flags().Bool("personal", false, "Adds a personal touch (aka logo) to random pages")
	rootCmd.Flags().String("merge", "merged.pdf", "--merge=filename, merge generated files into <filename>")
	rootCmd.Flags().String("separator", domain.DefaultSeparator, "Separator between chapter and page")
	rootCmd.Flags().String("page-count-prefix", domain.DefaultPageCountPrefix, "Prefix for total page count")

	// Bind all flags to viper - ensure this includes persistent flags too
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

// SettingsCmd requires flags to be evaluated, so that the final configuration can be determined
func SettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "settings",
		Aliases: []string{"setting", "set"},
		Short:   "Show the final configuration",
		Long:    "Show the final configuration, after defaults, config files, and flags have been evaluated.",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().Msg("executing settings command")
			
			// Get all command-line args
			var config domain.MinionConfig
			
			// Parse language flag
			languageArg := getLanguageArg()
			if languageArg != "" {
				lang := domain.ParseLanguageCode(languageArg)
				config = domain.NewDefaultConfig(lang)
			} else {
				config = domain.NewDefaultConfig(domain.MapSystemToAppLanguage())
			}
			
			// Apply other command-line settings that might be provided
			// (e.g., --page-prefix, --chapter-prefix, etc.)
			applyCommandLineSettings(cmd, &config)
			
			// Save these settings as the new defaults
			saveSettings(config)
			
			// Print the configuration
			domain.PrintFinalConfiguration(&config)
			return nil
		},
	}
	
	// Register all config flags with the settings command
	cmd.Flags().StringP("language", "l", "", "Override system language")
	cmd.Flags().StringP("page-prefix", "p", "", "Prefix for page numbers")
	cmd.Flags().StringP("chapter-prefix", "c", "", "Prefix for chapter numbers")
	// ...other flags
	
	return cmd
}

// Apply command-line settings to the config
func applyCommandLineSettings(cmd *cobra.Command, config *domain.MinionConfig) {
	// Check for page-prefix
	if pagePrefixFlag, _ := cmd.Flags().GetString("page-prefix"); pagePrefixFlag != "" {
		config.PageNrPrefix = pagePrefixFlag
	}
	
	// Check for chapter-prefix
	if chapterPrefixFlag, _ := cmd.Flags().GetString("chapter-prefix"); chapterPrefixFlag != "" {
		config.ChapterPrefix = chapterPrefixFlag
	}
	
	// ...other flags
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

// getLanguageArg scans command line arguments for the language flag
func getLanguageArg() string {
	// Direct command-line parsing for language flag
	for i, arg := range os.Args {
		if arg == "-l" || arg == "--language" {
			if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
				fmt.Println("Debug - Found language arg:", os.Args[i+1])
				return os.Args[i+1]
			}
		} else if strings.HasPrefix(arg, "--language=") {
			parts := strings.Split(arg, "=")
			if len(parts) == 2 {
				fmt.Println("Debug - Found language arg (--lang=):", parts[1])
				return parts[1]
			}
		} else if strings.HasPrefix(arg, "-l=") {
			parts := strings.Split(arg, "=")
			if len(parts) == 2 {
				fmt.Println("Debug - Found language arg (-l=):", parts[1])
				return parts[1]
			}
		}
	}
	return ""
}

// saveSettings saves the configuration for future use
// Currently this is a stub that could be expanded to save settings to a config file
func saveSettings(config domain.MinionConfig) {
	// For now, this is a no-op stub
	// Future implementation could save to a config file
	log.Debug().Msg("Settings viewed (save not implemented yet)")
}
