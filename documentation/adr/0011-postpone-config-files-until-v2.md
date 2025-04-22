# 11. postpone-config-files-until-v2

Date: 2025-01-05

## Status

Accepted

## Context

Layered configuration with defaults and config files as described in ADR-0008 is tricky from an implementation viewpoint, as four different configurations have to be merged.
In addition, the file-based configuration has many edge cases that need to be handled, such as file-not-found, wrong-file-format, and so on.

## Decision

We decided to postpone the implementation of file-based configuration until the next major release, and focus on the core functionality of PDFminion instead.


## Consequences

- Need to change the public website.
- Need to change ADR-0008