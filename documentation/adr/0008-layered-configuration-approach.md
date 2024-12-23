# 8. layered-configuration-approach

Date: 2024-12-16

## Status

Accepted

## Context

PDFminion requires a mechanism to parse and process command-line arguments for a variety of commands and options, such as `--source`, `--target`, `--language`, etc.

Users shall be able to:

* save their preferred configuration in configuration files
* overwrite their saved configs with command-line flags
* rely on default settings for minimal effort

This is called a _layered configuration approach_:

### Details Requirements

#### Commands
Several commands shall be immediately executed. PDF processing shall not happen in these cases:
- version, help, ll, list-languages, settings

The default-command is pdf processing.
It shall be invoked if no other command is given.

#### Flags
The default command can be configured with numerous flags, concerning the source, target, language, overwrite target (force), merge generated files, and various page-related and language-dependend settings.
These settings can be given via command-line flags or a configuration file.
The name of this config file can be given via command-line flags (--file or --config).
The default name of this file is pdfminion.yaml.
The location of the configuration file is the current directory or the users' home directory or both.
In case both files are present, the current directory file shall have precedence.
The settings command can also accept the --config or --file flag to read the settings from a file and then display them.
The configuration shall be layered, from low to high priority:

1. Default values (hard coded in the application)
2. Configuration file in users' home directory
3. Configuration file in the current directory
4. Command-line flags

Default values depend on the language setting. The default language is English.
The application supports a number of languages (currently EN, DE and FR).
Upon startup, the application shall determine the language setting by using the library xuanwo/go-locale.


## Decision


### Define the Configuration Workflow:

* Load default values first.
* Override with values from the configuration file (if provided via --file or exists as pdfminion.cfg).
* Override further with command-line flags.
* Library Choice: Use Viper for configuration file handling and defaults.
* Use the Cobra package for command-line parsing. Cobra integrates well with Viper for CLI apps.




### Implementation Outline:





## Consequences

We test this slightly more complex configuration approach on a branch(#15-layered-configuration)

 