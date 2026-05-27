# Runbook: Database Unavailable

## Symptoms

- All API requests return 500.
- `/readyz` returns `{"status":"not ready"}`.
- Logs show `connection refused` or `dial timeout`.

## Checks

1. `docker compose ps` — is the postgres container running?
2. `make psql` — can you connect manually?
3. Check `BLOG_DB_URL` in `.env` — correct host, port, credentials?
4. Check disk space: `df -h` — Postgres stops if disk is full.

## Recovery

1. If container is down: `docker compose up -d`.
2. If disk is full: free space, then restart: `docker compose restart postgres`.
3. If credentials are wrong: update `.env`, restart the API.
4. Verify recovery: `curl http://localhost:8080/readyz`.
