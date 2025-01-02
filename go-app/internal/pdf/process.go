package pdf

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"pdfminion/internal/domain"
	"sort"
)

var appConfig domain.MinionConfig

func ProcessPDFs(cfg *domain.MinionConfig) error {
	log.Debug().Msg("Starting PDF processing") // Only shown in debug mode

	// store configuration so it becomes usable in other functions
	appConfig = *cfg

	if cfg.Verbose {
		fmt.Println("Starting PDF processing")
	}

	InitializePDFInternals()

	// TODO: remove cfg from function signature
	files, err := CollectCandidatePDFs()

	if err != nil {
		return fmt.Errorf("error collecting candidate PDFs: %w", err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	pdfFiles, nrOfValidPDFs := ValidatePDFs(files)

	err = CopyValidatedPDFs(pdfFiles, cfg.SourceDir, cfg.TargetDir, cfg.Force)
	if err != nil {
		return fmt.Errorf("error during copy: %w", err)
	}

	if cfg.Verbose {
		fmt.Printf("Found %d PDF files\n", len(files))
	}
	log.Debug().Int("fileCount", len(files)).Msg("Found files")

	Evenify(nrOfValidPDFs, pdfFiles)
	AddPageNumbersToAllFiles(nrOfValidPDFs, pdfFiles)

	return nil
}
