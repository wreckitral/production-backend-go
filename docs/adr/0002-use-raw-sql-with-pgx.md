# 0002 Use Raw SQL with pgx

## Status

Accepted

## Context

We need a database access strategy. Options considered: GORM (ORM), sqlc (codegen), raw SQL with pgx.

## Decision

Use raw SQL with `pgx/v5` directly. No ORM, no codegen step.

## Consequences

Full control over queries. No hidden N+1 queries or magic. Scan boilerplate is
manual but predictable. If query count grows significantly, migrating to sqlc
is straightforward since we already write raw SQL.
