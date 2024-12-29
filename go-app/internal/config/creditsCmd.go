package config

import (
	"github.com/spf13/cobra"
	"pdfminion/internal/domain"
)

func CreditsCmd() *cobra.Command {
	return &cobra.Command{
		Use:              "credits",
		Aliases:          []string{"credit", "thanx"},
		Short:            "Show credits",
		Long:             "Show credits for the open-source tools used to implement PDFMinion",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		Run: func(cmd *cobra.Command, args []string) {
			domain.GiveCredits()
		},
	}
}
