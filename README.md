# WorkPulse

A Go (Gin + GORM) + React (Ant Design) web system for OKR-aligned daily work management.

## Quick start

### 1) Database
Create a Postgres database and run the migration:

```bash
psql "$DB_DSN" -f migrations/001_init.sql
```

### 2) Backend
```bash
cp .env.example .env
# edit DB_DSN, JWT_SECRET
make run
```

### 3) Frontend
```bash
cd web
npm install
npm run dev
```

Vite proxies `/api` to `http://localhost:8080`.

## Notes
- Auth in this scaffold expects a JWT with `user_id` and `org_id` claims.
- Many endpoints are stubs; the goal is to provide a clean, evolvable skeleton.
