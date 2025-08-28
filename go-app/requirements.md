## Requirements for PDFminion go app

I already have a few PDF related functions (page numbering, evenification) implemented.
The handling of commands and flags is not satisfactory, and I want to improve it.


### Functional requirements

- command line app, with flags and commands
- The following commands should be supported:
    - list-languages: list all available languages for the --language flag
    - version: prints version information
    - help
    - settings: prints all current settings, e.g. for page-prefix, chapter-prefix, etc.
    - credits: prints credits for the creators of several OS libraries used in the app

If no commands are given, all flags are evaluate, validated and PDF processing is started.

- The following flags should be supported:
- `--source`, short `-s`: Specifies the input directory for PDFs to be processed.
- `--target`, short `-t`: Specifies the output directory for generated PDFs.
- `--force`, short `-f`: Specifies whether existing files should be overwritten.
- `--evenify {=true|false}`, short `-e {=true|false}`: Enables or disables adding blank pages for even page counts. Default: `true`.
- `--toc`, short `-o`: Specifies whether a table of contents should be generated.
- `--config <filename>`, short `-c <filename>`: Specifies the configuration file to be used. It should be in YAML format.
- `--language`, short `-l`: Specifies the language for the generated PDFs.
- `--page-prefix <text>`, short `-p`: Specifies the prefix for page numbers. Default: "Page". Example: `pdfminion --page-prefix "Page"`
- `--chapter-prefix <text>`: Specifies the prefix for chapter numbers. Default: "Chapter". Example: `pdfminion --chapter-prefix "Ch."`
- `--verbose`, short `-v`: Specifies whether verbose output should be printed.
- `--running-head <text>`, short `-r`: Sets text for the running head at the top of each page. Default is "" (no header).-
- `--blankpagetext <text>`: Specifies text printed on blank pages added during evenification. Example: `pdfminion --blankpagetext "deliberately left blank"`
- `--separator <symbol>`: Defines the separator between chapter, page number, and total count. Default: `-`. Example: `pdfminion --separator " | "`
- `--page-count-prefix <text>`: Sets the prefix for the total page count. Default: "of". Example: `pdfminion --page-count-prefix "out of"`


One or several of these entries can be given in the (optional) configuration file.
A configuration file can be located in users' home directory, and/or in the current local directory.
If a flag is given on the command line, it should override all configuration files.
If a config file is given in the local directory, its settings shall override the config file settings from the users' home directory.

### Technical requirements

- current go version
- go modules
- automated tests where viable
- Makefile for building, running, linting, packaging
- generated binaries for Macos, Windows and Linux, each with specific make targets
- Dockerfiles to test the application in different environments, at least Windows and Linux
- Changes should be done on feature branches
