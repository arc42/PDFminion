# 10. use badges to show code status

Date: 2025-01-02

## Status

Accepted

## Context

The build process runs linter actions and automated tests, and the results are interesting
for both developers and users.
Therefore, they should be openly visible in the repository.


## Decision

Badges are a common way to show a variety of status infos in GitHub repositories.

We use:
* shields.io
* foss.com for security and license scans
* GitHub actions for CI/CD workflows
    - golangci-lint
    - go test
* [schneegans dynamic action](https://github.com/Schneegans/dynamic-badges-action) to display the number of tests

## Consequences

