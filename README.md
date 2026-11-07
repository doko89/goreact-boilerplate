# Boilerplate

Monorepo: Go API (pragmatic DDD) + React SPA (feature-based). One binary serves both.

```text
apps/api/   # Fiber 3, sqlc, sqlite/postgres, session cookie + CSRF + CSP
apps/web/   # React + Vite, feature-based, same-origin fetch
packages/api-contract/openapi.yaml
```

## Quickstart

```bash
make setup     # go mod download + npm install
make dev       # API :3000 + web :5173 (proxy /api -> :3000)
```

Open http://localhost:5173, login `admin` / `admin123` (seeded from `.env`).

Production single binary:

```bash
make build     # web bundle -> embedded into bin/api
./bin/api      # serves API + SPA on :3000
```

## Backend layout (one folder = one bounded context)

```text
apps/api/
  cmd/api/main.go
  internal/
    shared/        # config (.env), apperr, httpx, id
    platform/      # db (open/migrate/seed), session (HMAC cookie), httpserver (wiring)
    auth/          # domain/ application/ infrastructure/ presentation/
    order/         # domain/ application/ infrastructure/ presentation/
  db/queries/      # sqlc queries (one set runs on both drivers)
  web/dist/        # embed target, filled by make build-web
```

Rules: `presentation -> application -> domain <- infrastructure`, interfaces in
`domain`, sqlc + bcrypt in `infrastructure`, no framework imports in `domain`.

## Frontend layout (one folder = one feature)

```text
apps/web/src/
  app/             # router
  shared/lib/      # api client (cookies + CSRF header + retry), cookie reader
  shared/ui/       # Layout, RequireAuth
  features/auth/   # api.ts (zod), AuthProvider, pages/LoginPage
  features/order/  # api.ts (zod), pages/OrdersPage
```

Features never import each other, only `shared/`. No token in localStorage:
auth state lives in the HttpOnly cookie, UI state in `AuthProvider`.

## Environment (.env)

| Key               | Default          | Notes                          |
| ----------------- | ---------------- | ------------------------------ |
| APP_ENV           | dev              | dev \| prod \| test             |
| APP_PORT          | 3000             |                                |
| DB_DRIVER         | sqlite           | sqlite \| postgres              |
| SQLITE_PATH       | ./data/app.db    | `:memory:` allowed              |
| DATABASE_URL      | —                | required for postgres           |
| SESSION_SECRET    | —                | required in prod                |
| SESSION_TTL_HOURS | 12               |                                |
| ADMIN_USER        | admin            | seed when users table empty     |
| ADMIN_PASS        | admin123         | required in prod                |

`.env` is resolved by walking upward from cwd, so both `make dev-*` and
`./bin/api` find the root file. Real env vars always override `.env`.

## Security

- **HttpOnly session**: HMAC-signed token in `session` cookie
  (`HttpOnly`, `SameSite=Lax`, `Secure` in prod). No JWT library, no storage.
- **CSRF**: Fiber csrf middleware on `/api/*`. SPA reads `csrf_token` cookie
  and sends `X-CSRF-Token` on POST/PATCH; client retries once on rotation.
- **CSP**: helmet middleware, `default-src 'self'`, `frame-ancestors 'none'`,
  re-applied on static responses. Same-origin design means no CORS at all.
- Passwords: bcrypt. Login errors are uniform (no user enumeration).

## Database

`DB_DRIVER=sqlite` (dev, zero deps) or `postgres` (prod). One sqlc query set
runs on both: TEXT ids, BIGINT ints, RFC3339 TEXT timestamps, `$1` params.
Schema auto-migrates on boot (`CREATE TABLE IF NOT EXISTS`).

```bash
make sqlc   # regenerate after editing db/queries/*.sql
```

## Commands

| Command     | Effect                                  |
| ----------- | --------------------------------------- |
| make setup  | install Go + npm deps                   |
| make dev    | API + web together                      |
| make build  | full build -> `./bin/api`               |
| make run    | build + run single binary               |
| make test   | backend tests                           |
| make vet    | gofmt + go vet + frontend typecheck     |
| make sqlc   | regenerate sqlc                         |
| make clean  | drop bin/, data/, embedded bundle       |
| make smoke  | probe running :3000 (CSRF + CSP + 403)  |
