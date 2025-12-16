# WorkPulse

A Go (Gin + GORM) + React (Ant Design) web system for OKR-aligned daily work management.

## Quick start

### 1) Database
Create a Postgres database and run the embedded migrations (versioned via golang-migrate):

```bash
make migrate
# or: go run ./cmd/migrate
```

The app automatically reads environment variables from a local `.env` file (see `.env.example`). `DB_DSN` must be set; otherwise, the command will exit with a helpful error instead of hanging on a missing database connection.

### 2) Backend
```bash
cp .env.example .env
# edit DB_DSN, JWT_SECRET
make run
```

Environment options:
- `API_VERSION` (default `v1`) controls the gateway prefix for REST routes and docs.
- `MIGRATIONS_TABLE` (default `schema_migrations`) sets the migration ledger table.
- `FEATURE_FLAGS` (JSON or env map) toggles UI features like `analytics`, `apiDocs`, `graphqlDocs`.

### 3) Frontend
```bash
cd web
npm install
npm run dev
```

Vite proxies `/api` to `http://localhost:8080`.

### API docs
- OpenAPI: `GET /api/{API_VERSION}/docs/openapi.json`
- GraphQL SDL preview: `GET /api/{API_VERSION}/docs/graphql.sdl`

## Notes
- Auth in this scaffold expects a JWT with `user_id` and `org_id` claims.
- Many endpoints are stubs; the goal is to provide a clean, evolvable skeleton.

## Operations
- CI/CD workflows live in `.github/workflows` (see `ci.yml` and `cd.yml`).
- Observability, backup, and deployment runbooks are captured in `docs/ops.md`.

## Internationalization (i18n)
- UI supports Chinese (zh) and English (en).
- Language is persisted in localStorage key `workpulse_lang`.
