# 6. Use Go's Standard `flag` Package for CLI Parsing

Date: 2024-12-02
updated: 2024-12-16

## Status

superseeded by ADR-0008 (layered config handling).

## Context
PDFminion requires a mechanism to parse and process command-line arguments for a variety of commands and options, such as `--source`, `--target`, `--language`, etc. Various libraries are available for CLI parsing in Go, including the standard `flag` package and third-party libraries like `cobra` or `urfave/cli`. This decision focuses on using the standard library's `flag` package.

## Decision
We will use Go's standard `flag` package for CLI argument parsing and processing.

Update 2024-12-16: As we need file-based configuration, the flag package is not versatile enough.



### Rationale
1. **Simplicity**: The `flag` package is part of Go's standard library, meaning no external dependencies are required.
2. **Ease of Use**: It provides a straightforward API that matches the needs of PDFminion's relatively simple CLI structure.
3. **Performance**: Being part of the standard library, `flag` is lightweight and optimized.
4. **Portability**: Reduces the need for third-party dependencies, improving maintainability and compatibility with various environments.
5. **Sufficient for Current Needs**: PDFminion's CLI does not require advanced features like nested commands, making `flag` a good fit.

### Example
Here is a simple example of how the `flag` package will be used:

```go
package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define CLI flags
	source := flag.String("source", "", "Specifies the input directory for PDF files.")
	target := flag.String("target", "", "Specifies the output directory for processed files.")
	help := flag.Bool("help", false, "Displays usage information.")

	// Parse the flags
	flag.Parse()

	// Process the flags
	if *help {
		flag.Usage()
		return
	}

	fmt.Printf("Source: %s\n", *source)
	fmt.Printf("Target: %s\n", *target)
}
