# 0004 Migration Strategy

## Status

Accepted

## Context

We need a way to manage database schema changes across environments.

## Decision

Use `golang-migrate` with sequential numbered SQL files in `migrations/`. Up
and down migrations are both required. Migrations run manually
via `make migrate-up` and in CI before tests.

## Consequences

Schema changes are versioned and auditable. Rollback is possible via
`make migrate-down`. Dirty state requires manual intervention, see the
migration failed runbook.
