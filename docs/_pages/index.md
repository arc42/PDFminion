---
title: "Welcome to PDFminion"

layout: page

permalink: /

header:
    overlay_image: /assets/images/header.webp
    overlay_filter: rgba(0, 0, 0, 0.6)

    actions:
        - label: "GitHub Repo <i class='fab fa-github'></i>"
          url: https://github.com/arc42/pdfminion
          blank: true

excerpt: "**For all those who like handouts with page numbers and running headers!**"
---



<section>

     <p>
        PDFminion adds page numbers and running-headers on pdf documents, helping to produce useful handouts.
        <br>
        It's open-source, runs on all major platforms and is free to use.
        <br>
        <p>Brought to you by</p>
        <div class="logo">
            <a href="https://arc42.org/" target="_blank" rel="noopener noreferrer nofollow">
            <img src="assets/images/arc42-logo.png" alt="arc42">
            </a>
        </div>

    </p>

</section>




<br>
<hr class="section-sep">
<br>



<section id="features">

    <h1>Features</h1>

        <div class="box-container">

        <div class="box box--primary box-third">
            <img src="assets/images/functions/page-number.png" alt="page-numbering" class="img-half">
            <h5>Page Numbers</h5>
            <p>Adds consecutive page numbers to all documents.</p>
        </div>

        <div class="box box--primary box-third">
            <img src="assets/images/functions/running-header.png" alt="running-head" class="img-half">
            <h5>Running Header</h5>
            <p>Adds a header to the top of all pages in all documents. </p>
        </div>

        <div class="box box--primary box-third">
            <img src="assets/images/functions/mascot.png" alt="we have a sweet mascot" class="img-half">
            <h5>Add Personal Touch</h5>
            <p>Who said PDF documents couldn't have a personal touch? Thanx to <a href="https://www.sketchnotes.tech/">Lisa, @teapot418</a>.</p>
        </div>

        <div class="box box--primary box-third">
            <img src="assets/images/functions/chapter-number.png" alt="chapter numbers" class="img-half">
            <h5>Chapter Number</h5>
            <p>Every document makes a chapter within the output. This feature allows to add the chapter number to all pages.</p>
        </div>


       <div class="box box--primary box-third">
            <img src="assets/images/functions/evenify.png" alt="evenify" class="img-half">
             <h5>Evenify</h5>
             <p>Make all documents have an even number of pages by adding a blank page to files with odd page count.</p>
       </div>
 <div class="box box--primary box-third">
            <img src="assets/images/functions/privacy.png" alt="privacy" class="img-half">
             <h5>Privacy</h5>
             <p>We value your privacy: Your documents are yours, forever. We do neither collect data nor retain any information about your documents.</p>
       </div>

 <div class="box box--primary box-third">
            <img src="assets/images/functions/multi-language.png" alt="multi-language" class="img-half">
             <h5>Multi-Language</h5>
             <p>Handle a number of natural languages, starting with English (EN)m German (DE), French (FR) and Spanish (ES).</p>
       </div>

       <div class="box box--primary box-third">
            <img src="assets/images/functions/toc.png" alt="table of contents" class="img-half">
            <h5>Table of Contents</h5>
            <p>Create table-of-contents with document names, their chapter number and starting page number.</p>
        </div>

        <div class="box box--primary box-third">
            <img src="assets/images/functions/merge.png" alt="merge documents" class="img-half">
            <h5>Merge</h5>
            <p>Create a single output document from all inputs. The output file name is configurable.</p>
        </div>
    </div>

</section>



<br>
<hr class="section-sep">
<br>


<section id="installation" markdown="1">

<h1>Installation</h1>

PDFminion runs on all major (desktop) operating systems.
You can download an appropriate version and install it yourself, or you can use one of our installation options.

If you're nerdy enough, you can [fork the repo](https://github.com/arc42/pdfminion), and build your own executable version, either with `make` or `go build`.

<div class="box-container">

<div class="box box--primary box-third" markdown="1">
<h4><i class="fab fa-apple"></i> MacOS</h4>
<p>How to install on MacOS</p>
PDFminion will be installable with [Homebrew](https://brew.sh), the package manager for MacOS.
But currently, this installation method is not available...

</div>

<div class="box box--primary box-third" markdown="1">
<h4><i class="fab fa-windows"></i> Windows</h4>
<p>How to install on Windows</p>
PDFminion will be installable with [Chocolatey](https://chocolatey.org/), the package manager for Windows.
But currently, this installation method is not available...
</div>

<div class="box box--primary box-third" markdown="1">
<h4><i class="fab fa-linux"></i> Linux</h4>
<p>How to install on Linux</p>
PDFminion will be installable with [Snapcraft](https://snapcraft.io/snapcraft), the package manager for Linux, that works for most distributions.
But currently, this installation method is not available...
</div>

</div>



</section>


<br>
<hr class="section-sep">
<br>


<section id="terminology"></section>


<h1>Terminology</h1>

<img src="assets/images/page-terminology.png" alt="page terminology"/>

<h2>Evenify</h2>

To evenify a file means adding a blank page at the end of the file if the page-count is odd (1, 3, 5 or such).
That means that the first page of every file in a group will always start on the front-page of paper,
even in case of double-sided printing.

Chapters in technical or scientific books traditionally start on odd (right-hand) pages to ensure consistency,
readability, and prominence, aligning with classic book design practices.


<img src="assets/images/page-terminology-evenify.png" alt="evenify"/>


<br>
<hr class="section-sep">
<br>


<section id="usage"></section>


<h1>Usage</h1>

<h2>Terminology and Conventions</h2>

<div markdown="1" class="left-aligned-section">
* **Commands** are executed immediately. They are given without `--`, for example `pdfminion version`
* Only the first command will be executed at a time.
* **Flags** (_configuration settings_) control the behaviour of the actual processing. They are given with `--`, for example `pdfminion --force` or `pdfminion --source ./input`
* Configurations (_flags_) can also be set via a configuration file, either in users' home directory or in the current working directory.
    The default name is `pdfminion.yaml` or `pdfminion.yml`. Other names can be specified with the `--config` flag.
</div>

<h2>Commands</h2>

| **Name**  | **Long Name**  | **Shorthand** | **Description** |
|-----------|-------------------|-------------------|-----------------|
| **Help**          | `help`      | `?`| Displays a list of supported commands and their usage.<br>Example: `pdfminion --help`|
| **List Languages**| `list-languages` | `ll`, `list`   | Lists available languages for the `--language` option.<br>Example: `pdfminion list-languages` |
| **Settings**      | `settings`  |         | Prints all current settings, like page-prefix, chapter-prefix etc. <br>Example: `pdfminion settings` |
| **Version**       | `version`   |     | Displays the current version of PDFminion.<br>Example: `pdfminion version` |
| **Credits**       | `credits`   |         | Gives credit to the maintainers of several OS libraries. <br>Example: `pdfminion credits`  |

If no command is given, all flags are evaluate, validated and PDF processing is started.

<h2>Basic Settings</h2>

| **Name**  | **Long Name**  | **Shorthand** | **Description**    |
|-----------|-------------------|-------------------|--------------------|
| **Source Directory** | `--source <directory>` | `-s <directory>`| Specifies the input directory for PDF files. Default is `./_pdfs` Example: `pdfminion --source ./input`|
| **Target Directory** | `--target <directory>` | `-t <directory>` | Specifies the output directory for processed files. Default is `_target`. Creates the directory if it doesn’t exist. Example: `pdfminion --target ./out`|
| **Force Overwrite**  | `--force`              | `-f`    | Allows overwriting existing files in the target directory. Default: `false`. Example: `pdfminion --force` |
| **Config File**  | `--config <filename>`  | `-c <filename>` | Loads configuration from a file. It needs to be a yaml file. Example: `pdfminion --config settings.yaml`  |
| **Verbose Mode**    | `--verbose`     |                | Gives detailed processing output. Default: `false`. Example: `pdfminion --verbose` |


<h2>Page Related Settings</h2>

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


<h2>Multi-Language Support</h2>

PDFMinion provides defaults for page processing for several languages.
With these commands you can change these defaults and provide your own values.


| **Name**| **Long Name**  | **Shorthand** | **Description** |
|-----------|-------------------|-------------------|-----------------|
| **Language**      | `--language <code>`        | `-l <code>`     | Sets the language for stamped text. Currently supports `EN` (English) and `DE` (German). Default: `EN`. Example: `pdfminion --language DE`. You get all supported languages with the `--list-languages` command. |
| **Blank Page Text** | `--blankpagetext <text>`   | `-b <text>`     | Specifies text printed on blank pages added during evenification. Example: `pdfminion --blankpagetext "deliberately left blank"`|




<h2>File-Related Settings</h2>
After all files have been processed, PDFminion can perform some post-processing.

| **Name**  | **Long Name**  | **Shorthand** | **Description** |
|-----------|-------------------|-------------------|-----------------|
| **Merge** | `--merge <filename>`       | `-m <filename>` | Merges input files into a single PDF. Uses default name if `<filename>` not provided. Not yet implemented. Example: `pdfminion --merge combined.pdf`   |
| **Table of Contents**  | `--toc`   |  | Generates a table-of-contents PDF. Not yet implemented. Example: `pdfminion --toc`|



<br>
<hr class="section-sep">
<br>



<section id="examples" class="left-aligned-section">

<h1 markdown="1">Examples</h1>

<div markdown="1" >



> **Example 1**: Default processing: Add page numbers and running headers to all PDF files in the `input` directory and save the processed files in the `output` directory. Requires the `output` to be empty!
>
>`$ pdfminion --source ./input --target ./output`
>
<hr>

> **Example 2**: Force overwrite of existing files in the `output` directory
>
>`$ pdfminion --force --source ./input --target ./output`
>

> **Example 3**: Show all current configuration items
>
>`$ pdfminion settings`
>

<hr>
> **Example 4**: Show the version info
>
>`$ pdfminion version`
>
<hr>

> **Example 5**: List available languages for the `--language` flag
>
>`$ pdfminion ll`
>`$ pdfminion list-languages`
>
<hr>

> **Example 6**: Set a language with the `--language` flag and override one of the defaults
>
>`$ pdfminion --language=EN --page-prefix=pg.`
>
<hr>

> **Example 41**: Gives credit to the maintainers of the open-source libraries used by PDFminion, e.g. [pdfcpu](https://pdfcpu.io/) and a few others..
>
>`$ pdfminion thanx`
>

PDFminion was created on the shoulder of giants...
```

<br>
<a href="/downloads/pdfminion-config.yaml" download="pdfminion-config.yaml" class="btn btn--primary">Download example config file</a>




</div>

</section>


<br>
<hr class="section-sep">
<br>
<section id="credits" markdown="1" class="left-aligned-section">
<h1>Credits</h1>

PDFminion uses numerous open-source libraries, and wish to thank the maintainers of the following projects:


* Horst Rutter for [pdfcpu](https://pdfcpu.io/), all things related to PDF processing.
* Hao Ding for [go-locale](https://github.com/Xuanwo/go-locale) to handle tricky locale settings.
* [Steve Francia](https://spf13.com) for a lot of stuff, including [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper), the command-line and configuration libraries.
* [Max Howell](https://brew.sh) for creating Homebrew, the package manager for MacOS.
* [Tom Preston-Werner & Co.](https://jekyllrb.com) for creating Jekyll, the static site generator.
* [Simon Schneegans](https://schneegans.github.io/tutorials/2022/04/18/badges) for the dynamic badge action.
* [INNOQ](https://www.innoq.com) for generous support of the arc42 open source initiative.
* and, of course, the [Go team](https://golang.org/) for creating the language that compiles to various operating systems, and refrains from adding fancy features every 6 month (pun intended).
</section>

