# Runbook: Migration Failed

## Symptoms

- Deploy is blocked.
- `migrate` reports dirty database version.
- API fails to start with schema errors.

## Checks

1. `make migrate-version`, what version is the database at?
2. Inspect the failed migration SQL in `migrations/`.
3. Check whether objects were partially created: `make psql` then `\d` to list
tables.

## Recovery

1. Manually repair the schema if objects were partially created.
2. Force the version only after verifying: `migrate force <version> -path
migrations -database "$BLOG_DB_URL"`.
3. Re-run: `make migrate-up`.
4. Record what happened in incident notes.

## Prevention

Always test migrations against a copy of production data before deploying.
