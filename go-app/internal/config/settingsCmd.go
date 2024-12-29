package config

import (
	"github.com/spf13/cobra"
	"pdfminion/internal/domain"
)

func SettingsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "settings",
		Aliases: []string{"setting", "set"},
		Short:   "Show the final configuration",
		Long:    "Show the final configuration, after defaults, config files, and flags have been evaluated.",
		Run: func(cmd *cobra.Command, args []string) {
			domain.PrintFinalConfiguration(*minionConfig)
		},
	}
}
