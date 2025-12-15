# Operational playbook

This document captures day-2 guidance for keeping WorkPulse healthy. It pairs with the automation under `.github/workflows` and the telemetry hooks in the API.

## CI: lint/test/build
- **Workflow:** `.github/workflows/ci.yml` triggers on pushes and PRs.
- **Backend:** `go vet`, `go test`, and `go build` run against the whole module to block regressions.
- **Frontend:** `npm ci` + `npm run build` validates the Vite/TypeScript bundle using the locked dependencies in `web/package-lock.json`.
- **Local parity:** run `go vet ./... && go test ./... && go build ./...` and `cd web && npm ci && npm run build` before opening a PR.

## CD: staging + production canary
- **Workflow:** `.github/workflows/cd.yml` is triggered manually via `workflow_dispatch` with inputs `environment` (`staging`/`production`) and optional `canary_percentage`.
- **Build artifacts:** the workflow produces a `workpulse-release.tar.gz` containing `dist/workpulse-api` and `web/dist`. Artifacts are uploaded for downstream jobs.
- **Staging rollout:** when `environment=staging`, the `staging` job downloads the artifact and provides a hook to push to the pre-prod host (replace the commented `scp/systemctl` lines with your deploy commands).
- **Production canary:** when `environment=production`, the `canary` job is a guarded step for partial traffic. Wire it to your ingress or service mesh (e.g., `kubectl apply -f k8s/canary.yaml`) using the provided percentage input.
- **Promotion:** after canary, the `production` job represents full rollout. Swap the placeholder commands with your orchestrator (Kubernetes, Nomad, etc.).

## Observability (OpenTelemetry + Prometheus/Grafana)
- **Config:** set `SERVICE_NAME` (default `workpulse-api`), `ENV`, `OTEL_EXPORTER_OTLP_ENDPOINT`, optional `OTEL_EXPORTER_OTLP_HEADERS`, and `METRICS_PATH` in the environment. See `.env.example`.
- **Tracing:** the server wires `otelgin` middleware; if `OTEL_EXPORTER_OTLP_ENDPOINT` is set, spans are shipped via OTLP/HTTP. Headers allow authenticated collectors.
- **Metrics:** a Prometheus exporter is mounted at `METRICS_PATH` (default `/metrics`) on the main HTTP server. Scrape it from Prometheus and visualize in Grafana; common dashboards include request latency, status code counts, and DB pool metrics once instrumented.
- **Shutdown:** telemetry providers flush on process exit to minimize span loss.

## Logging
- Gin access logs remain enabled for edge visibility. Collector-side enrichment is recommended: run Fluent Bit/Vector to ship logs alongside trace IDs from the `traceparent` header emitted by the middleware.

## Backups and export interfaces
- **Database backups:**
  - Nightly logical dumps (`pg_dump`) to object storage (e.g., S3) with retention tiers (daily for 14 days, weekly for 12 weeks). Include schema + data.
  - Optional point-in-time recovery by archiving WAL segments with `wal-g` or native `archive_command` to the same bucket.
  - Tag backup objects with commit SHA and migration version for traceability.
- **Disaster recovery drills:** restore quarterly into a staging database and run smoke tests (API health + representative queries) to validate the process.
- **Data export API:** keep `/api/{API_VERSION}/reports/export` for user-driven exports; add rate limiting and audit logging via middleware. For larger exports, surface signed URLs referencing generated CSV/Parquet dumps stored in object storage.

## Full-text search (optional)
- **Choice:** OpenSearch/Elasticsearch for scale or lightweight Meilisearch for starter deployments.
- **Sync strategy:** emit change events from GORM hooks or a queue processor to keep the search index consistent with PostgreSQL.
- **Query surface:** expose a `/search` endpoint that fan-outs to the search backend while enforcing RBAC filters already present in middleware.

## Environment reset (local/staging)
- **Stop running processes:** halt `make run`/`npm run dev` and tear down any local containers to avoid open connections during cleanup.
- **Drop + recreate the database:** with Postgres running locally, execute `psql -c "DROP DATABASE IF EXISTS workpulse;"` followed by `psql -c "CREATE DATABASE workpulse;"` (or your custom name). Re-run migrations with `make migrate` to repopulate schema and seed tables tracked in the migrations ledger.
- **Purge local caches:** `go clean -cache -modcache` removes compiled artifacts and downloaded modules; `rm -rf web/node_modules && npm ci` resets the frontend dependency tree.
- **Config reset:** delete or recreate `.env` (copy from `.env.example`) to clear credentials and feature flags. This ensures backend and frontend pick up a clean set of environment variables on the next start.
