package config

import (
	"github.com/spf13/cobra"
	"pdfminion/internal/domain"
)

func VersionCmd() *cobra.Command {
	
	return &cobra.Command{
		Use:              "version",
		Short:            "Show version information",
		Long:             "Show version information for the PDFMinion application, including version number, build date and build target",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		Run: func(cmd *cobra.Command, args []string) {
			domain.PrintVersion()
		},
	}
}
