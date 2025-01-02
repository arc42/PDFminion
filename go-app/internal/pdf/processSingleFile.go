package pdf

import (
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/rs/zerolog/log"
	"pdfminion/internal/util"
	"strconv"
)

// TODO: Split Evenify function into loop-over-all-files and evenify-single-file

func Evenify(nrOfValidPDFs int, pdfFiles []SingleFileToProcess) {
	relaxedConf := model.NewDefaultConfiguration()
	relaxedConf.ValidationMode = model.ValidationRelaxed

	for i := 0; i < nrOfValidPDFs; i++ {
		if !util.IsEven(pdfFiles[i].PageCount) {
			// add single blank page at the end of the file
			_ = api.InsertPagesFile(pdfFiles[i].Filename, "", []string{strconv.Itoa(pdfFiles[i].PageCount)}, false, relaxedConf)

			pdfFiles[i].PageCount++

			onTop := true
			update := false

			wm, err := api.TextWatermark(appConfig.BlankPageText, "font:Helvetica, points:48, col: 0.5 0.6 0.5, rot:45, sc:1 abs",
				onTop, update, types.POINTS)
			if err != nil {
				log.Printf("Error creating watermark configuration %v: %v\n", wm, err)
			} else {

				err = api.AddWatermarksFile(pdfFiles[i].Filename, "", []string{strconv.Itoa(pdfFiles[i].PageCount)}, wm,
					relaxedConf)

				if err != nil {
					log.Printf("error stamping blank page in file %v: %v\n", pdfFiles[i].Filename, err)
				}

			}
			if appConfig.Verbose {
				fmt.Printf("File %s was evenified\n", pdfFiles[i].Filename)
			}
			log.Debug().Str("File %s\n", pdfFiles[i].Filename).Msg("was evenified")
		}
	}
}
