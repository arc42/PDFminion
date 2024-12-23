package domain

import (
	"fmt"
)

// PrintFinalConfiguration prints the final configuration
func PrintFinalConfiguration(myConfig MinionConfig) {
	fmt.Println("Final Configuration:")
	fmt.Printf("Source directory: %s\n", myConfig.SourceDir)
	fmt.Printf("Target directory: %s\n", myConfig.TargetDir)
	fmt.Printf("Language: %s\n", myConfig.Language)
	fmt.Printf("Debug: %t\n", myConfig.Debug)
	fmt.Printf("Force: %t\n", myConfig.Force)
	fmt.Printf("Evenify: %t\n", myConfig.Evenify)
	fmt.Printf("Merge: %t\n", myConfig.Merge)
	fmt.Printf("Merge file name: %s\n", myConfig.MergeFileName)
	fmt.Printf("Running header: %s\n", myConfig.RunningHeader)
	fmt.Printf("Chapter prefix: %s\n", myConfig.ChapterPrefix)
	fmt.Printf("Separator: %s\n", myConfig.Separator)
	fmt.Printf("Page prefix: %s\n", myConfig.PagePrefix)
	fmt.Printf("Total page count prefix: %s\n", myConfig.TotalPageCountPrefix)
	fmt.Printf("Blank page text: %s\n", myConfig.BlankPageText)
}
