# Product Requirements Document (PRD): PDFMinion
- Version: 1.0.0
- Date: 2025-04-20
- Author: Dr Gernot Starke

## 1. Product Overview

### 1.1 Document Title and Version
- PRD: PDFMinion
- Version: 1.0.0

### 1.2 Product Summary
MyGoCLI is a command-line utility written in Go to add page numbers and running-headers on existing pdf documents, helping to produce useful handouts.

## 2. Goals

### 2.1 Business Goals
- Help create printable handouts out of existing pdf documents, where each documents is one chapter of the resulting handout, especially by adding page and chapter numbers, running headers and filler pages for pdfs with uneven page count.
- Enable to create a useful table-of-contents from the existing pdfs.
  

### 2.2 Developer Goals
- Simple, understandable codebase
- Use standard Golang coding style
- Use established open source libraries

### 2.3 Non-Goals
- No GUI interface
- No cloud storage integration

## 3. User Personas

### 3.1 User Types
- Trainers (aka Users) that create pdfs for their courses.
- Software Developers that maintain PDFMinion
- Trainees that use the resulting handouts

### 3.2 Trainer Details
Trainers are the main users of PDFMinion.
Trainers create chapter-wise pdfs with arbitrary tools that contain their course material and content. 
These pdfs have no page numbers and no running headers.
Usually they have 5-15 pdfs with strict ordering, often numbered 1-...pdf, 2-...pdf, 3-...pdf and so on.

After beeing processed by PDFMinion, these pdfs shall have consecutive page numbers over all files.

If the Trainer wants, PDFMininon can also add chapter numbers and running headers to these files. 
That will be detailed in the functional requirements, especially in the sections on commands and flags.


### 3.3 Software Developer Details
Developers of PDFMinion want a clean, highly modular codebase, using the Golang programming language.

Typical software engineering practices shall be followed,
and understandability of the codebase shall be a leading quality goal for development.

### 3.4 Trainee Details
Trainees receive their handouts either in printed and bound format or digital as either single pdf with all chapters merged into one large document or with the single page-numbered pdfs.

These handouts will contain a table-of-contents with the starting page of each chapter, so they can quickly and effortlessly navigate the complete document (be it printed or digital)




## 4 Domain Terminology

* **Source directory**: The existing pdfs already exist in the source directory.
* **Target directory**: The pdfs processed by PDFMinion shall be put in the target directory.
* **Evenify**: Chapters in technical or scientific books traditionally start on odd (right-hand) pages to ensure consistency, readability, and prominence, aligning with classic book design practices. To **evenify** a file means adding a blank page at the end of the file if the page-count is odd (1, 3, 5 or such). That means that the first page of every file in a group will always start on the front-page of paper, even in case of double-sided printing.
* **Chapter**: A single pdf file constitutes a chapter within the resulting handout.
* **Flag**: (aka configuration settings) control the behaviour of the actual processing. They are given with --, for example `pdfminion --force` or `pdfminion --source ./input`. Flags can have parameters, like the `./input` in the preceeding example of the `--source` flag. Configurations (flags) can also be set via a configuration file, either in users’ home directory or in the current working directory. The default name is `pdfminion.yaml or pdfminion.yml. Other names can be specified with the `--config` flag.
* **Command**: Commands are executed immediately. They are given without --, for example `pdfminion version`. Only the first command will be executed at a time. Subsequent commands will be ignored.
* **Running Header**: A constant text (max a single line) that is repeated on top of every page.
* **Chapter prefix**: A string constant that preceedes the chapter number. In English it is "chapter", in German "Kapitel".
* **Page prefix**: A string constant that preceedes the page number. In English it is "page", in German "Seite"
* **Separator**: A string constant (usually just one character) that comes between the chapter prefix and number and the page prefix and number. Default is "-".
* **Page count prefix**: A string constant that comes between the page number and the (optional) total page count. In English it is "of", in German "von".
* **Blank page text**: The constant string (often spanning multiple lines) that is stamped in a large font on the pages added by the evenify operation. In English the default is "deliberately left blank", in German "Diese Seite bleibt absichtlich frei".

## 5. Functional Requirements

### 5.1 Core Functions
- Use via command line: PDFMinion shall be called from the command line (on the major operating systems MacOS, Windows and Linux).
- PDFMinion supports multiple languages out of the box, especially EN (English), DE (German) and FR (French). 
    * For each of the supported languages, default values for important string constants are predefined in the source code (e.g. page-prefix, chapter-prefix, page-count-prefix)
    * The default language can sometimes be deferred from OS settings. The fallback is English (EN).
    * The language to be used can be specified by a command-line flag (`--language DE` or `--language=DE`)
  * The runtime behaviour of PDFMininion is controlled via flags and commands via the command line.

- The following tables contain all commands and flags PDFMinion shall support:

### 5.2 Commands

| **Name**  | **Long Name**  | **Shorthand** | **Description** |
|-----------|-------------------|-------------------|-----------------|
| **Help**          | `help`      | `?`| Displays a list of supported commands and their usage.<br>Example: `pdfminion --help`|
| **List Languages**| `list-languages` | `ll`, `list`   | Lists available languages for the `--language` option.<br>Example: `pdfminion list-languages` |
| **Settings**      | `settings`  |         | Prints all current settings, like page-prefix, chapter-prefix etc. <br>Example: `pdfminion settings` |
| **Version**       | `version`   |     | Displays the current version of PDFminion.<br>Example: `pdfminion version` |
| **Credits**       | `credits`   |         | Gives credit to the maintainers of several OS libraries. <br>Example: `pdfminion credits`  |

If no command is given, all flags are evaluate, validated and PDF processing is started.

### 5.3 Flags for Basic Settings

| **Name**  | **Long Name**  | **Shorthand** | **Description**    |
|-----------|-------------------|-------------------|--------------------|
| **Source Directory** | `--source <directory>` | `-s <directory>`| Specifies the input directory for PDF files. Default is `./_pdfs` Example: `pdfminion --source ./input`|
| **Target Directory** | `--target <directory>` | `-t <directory>` | Specifies the output directory for processed files. Default is `_target`. Creates the directory if it doesn’t exist. Example: `pdfminion --target ./out`|
| **Force Overwrite**  | `--force`              | `-f`    | Allows overwriting existing files in the target directory. Default: `false`. Example: `pdfminion --force` |

<!-- see ADR-0011, config files have been postponed
| **Config File**  | `--config <filename>`  | `-c <filename>` | Loads configuration from a file. It needs to be a yaml file. Example: `pdfminion --config settings.yaml`  |
-->
| **Verbose Mode**    | `--verbose`     |                | Gives detailed processing output. Default: `false`. Example: `pdfminion --verbose` |


### 5.3 Flags for Page Related Settings

Set the running head, the page- and chapter prefix etc.

| **Name**  | **Long Name**  | **Shorthand** | **Description**    |
|-----------|-------------------|-------------------|--------------------|
| **Running Head**    | `--running-head <text>`    | `-r` | Sets text for the running head at the top of each page. Default  is "" (no header). Example: `pdfminion --running-head "Document Title"`|
| **Chapter Prefix**  | `--chapter-prefix <text>`  | `-c` | Specifies prefix for chapter numbers. Default: "Chapter". Example: `pdfminion --chapter-prefix "Ch."`|
| **Page Prefix**     | `--page-prefix <text>`     | `-p` | Sets prefix for page numbers. Default: "Page". Example: `pdfminion --page-prefix "Page"` |
| **Separator**       | `--separator <symbol>`     |  | Defines the separator between chapter, page number, and total count. Default: `-`. Example: `pdfminion --separator " | "`        |
| **Page Count Prefix**  | `--page-count-prefix <text>`|  | Sets prefix for total page count. Default: "of". Example: `pdfminion --page-count-prefix "out of"` |
| **Evenify**  | `--evenify {=true\|false}`  | `-e {=true\|false}`  | Enables or disables adding blank pages for even page counts. Default: true.  Example: `pdfminion --evenify=false |
| **Personal Touch**  | `--personal {on\|off}`  |   | Adds a personal touch (aka: Our PDFminion logo) on random pages. Not yet implemented. |

Please note: Most of these processing defaults are language-specific: The German language, for example, uses "Seite" for "Page" and "Kapitel" for "Chapter".

If you set a language (e.g German, DE), then the defaults of that language will be used.
You can still override some (or all) of these settings, though (see [example 6](#examples)).


### 5.4 Flags for Multi-Language Support

PDFMinion provides defaults for page processing for several languages.
With these commands you can change these defaults and provide your own values.


| **Name**| **Long Name**  | **Shorthand** | **Description** |
|-----------|-------------------|-------------------|-----------------|
| **Language**      | `--language <code>`        | `-l <code>`     | Sets the language for stamped text. Currently supports `EN` (English) and `DE` (German). Default: `EN`. Example: `pdfminion --language DE`. You get all supported languages with the `--list-languages` command. |
| **Blank Page Text** | `--blankpagetext <text>`   | `-b <text>`     | Specifies text printed on blank pages added during evenification. Example: `pdfminion --blankpagetext "deliberately left blank"`|




### 5.5 Flags for File-Related Settings

After all files have been processed, PDFminion can perform some post-processing.

| **Name**  | **Long Name**  | **Shorthand** | **Description** |
|-----------|-------------------|-------------------|-----------------|
| **Merge** | `--merge <filename>`       | `-m <filename>` | Merges input files into a single PDF. Uses default name if `<filename>` not provided. Not yet implemented. Example: `pdfminion --merge combined.pdf`   |
| **Table of Contents**  | `--toc`   |  | Generates a table-of-contents PDF. Not yet implemented. Example: `pdfminion --toc`|
 
### 5.6 Other Flags

Some commands can alternatively be invoked via flags, for those users who cannot remember the syntax.
See section 5.2 (Commands) above.

| **Name**  | **Long Name**  | **Shorthand** | **Description** |
|-----------|-------------------|-------------------|-----------------|
| **Help** | `--help`       | `-h` | displays the help text (Commends)  |
| **List Languages**| `--list-languages` | `--ll`, `--list`   | Lists available languages for the `--language` option. |
| **Settings**      | `--settings`  |         | Prints all current settings, like page-prefix, chapter-prefix etc. <br>Example: `pdfminion --settings` |
| **Version**       | `--version`   |`-v`     | Displays the current version of PDFminion.|
| **Credits**       | `--credits`   |         | Gives credit to the maintainers of several OS libraries. |


## 6. Additional Requirements
- The name of the binary shall be `pdfminion`.

### 6.1 Build process
- For development, `make` shall be used for compilation, test, local deployment and other development activities.
- `make mac` shall create the MacOS binary version
- `make windows` shall create the Microsoft-Windows binary version
- `make linux` shall create a Linux binary version

Other make targets shall be required later. The existing Makefile is good enough for now.

## 7. User Experience

### 7.1 Installation
- Developer installs via `make install`
- Trainers (users) can install via a package manager, depending on their operating system.

### 7.2 Core Experience

## 8. User Stories



## 9. Technical Considerations

- No external dependencies except Go standard library and well-known and established open-source libraries.
- Only libraries shall be used that are actively maintained 

- The whole PDFMinion project is maintained in a single GitHub repository
- The website (not described here) is located under "/docs".
- Technical documentation and architectural decisions under "/documentation"
- The Golang source code is located under "go-app".

- Ensure that human-readlable error messages are given by PDFMinion.
  

### 9.1 Instructions

- think before you answer.
- read and consume external documentation.
- always try to find the latest version of external documentation or libraries.

- The PDF manipulation code already

## 10. Current Tasks

0. Analyze the existing code and suggest improvements towards the requirements given here.
1. Ensure that handling flags and commands works flawlessly. 
2. Assume that the pdf handling code already exists and works as expected.
