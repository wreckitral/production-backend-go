# 0001 Use Vertical Slices

## Status

Accepted

## Context

We need a project layout that is easy to navigate and keeps feature changes
local. Layer-first layouts (controllers/, services/, repositories/) scatter a
single feature across multiple directories.

## Decision

Organize feature code under `internal/<feature>` with handler, service, repo,
messages, and routes together in one package.

## Consequences

Feature work stays in one place. Cross-cutting concerns (middleware, platform,
config) still live outside slices to avoid duplication.
