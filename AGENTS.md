# AGENTS.md — goreact-boilerplate

Go (Fiber 3, sqlc) + React (Vite, TS strict) monorepo. One binary serves API + embedded SPA.

> This repo is a WORKING SKELETON, not sample code. Your job is to EXTEND it
> to satisfy the PRD — never to delete it and scaffold from scratch. If your
> plan starts with emptying a folder, scaffolding a new app, or rewriting
> auth/server/API-client from zero: STOP, that plan is wrong.

## 0. Boilerplate-first (non-negotiable)

1. **Reuse before you create.** Every PRD requirement maps to an existing
   capability FIRST (table in §1). Only genuinely new behavior gets new files.
2. **Never mass-delete.** No `rm -rf`, no blanking `internal/*` or
   `src/features/*`, no "clean slate". Deleting >3 boilerplate files in one
   task requires explicit user approval, stated in your reply.
3. **Never re-scaffold.** No `npm create`, `vite create`, `go mod init`,
   second HTTP server, second router, second auth system, second API client,
   new ORM/query-builder, or new state/data-fetching library.
4. **PRD wording means "use the existing one".** "Login", "auth", "session",
   "upload", "API client", "validation", "config" = the boilerplate modules
   below, not a new implementation alongside them.
5. **Boilerplate wins on conflict.** If the PRD (or your instinct) conflicts
   with the security/architecture invariants in §3–§4, keep the boilerplate
   behavior and ask the user — don't silently downgrade it.
6. **Prove you looked.** Before creating any backend context or frontend
   feature, `read_file` its template sibling (`order`) and the `shared/`
   helper you will reuse, then follow that exact pattern.

## 1. PRD workflow (mandatory when building from a PRD)

1. **Inventory:** list which PRD items are already covered by the boilerplate
   (auth/session/CSRF/CSP, OIDC/LDAP, avatar upload, config, sqlite+postgres
   db, API client, AuthProvider, Layout/RequireAuth).
2. **Map:** for each PRD item write: `PRD item → existing module` OR
   `PRD item → NEW <context>/<feature> (justified: no existing module)`.
3. **Extend:** add new code ONLY as:
   - backend: new `internal/<context>/{domain,application,infrastructure,presentation}/`
     copied from `internal/order/`; queries in `db/queries/<context>.sql` + `make sqlc`;
     tables as `CREATE TABLE IF NOT EXISTS` in `platform/db/schema.sql`;
   - frontend: new `src/features/<name>/{api.ts, pages/, components/}` copied
     from `src/features/order/`; I/O only via `shared/lib/api.ts`.
4. **Wire, don't fork:** register new routes/handlers in the EXISTING
   `platform/httpserver/server.go` and the EXISTING `src/app/App.tsx` router.
5. **Verify:** `make vet && make test` passes, all pre-existing tests stay
   green, `git status` shows additions — not mass deletions.

Reuse map (check here before writing anything new):

| Need | Use this, don't rebuild |
|---|---|
| HTTP server, routes, CSRF, CSP, static serving | `apps/api/internal/platform/httpserver/` |
| Session cookie, login/logout/me | `apps/api/internal/auth/`, `platform/session/` |
| OIDC / LDAP login | `apps/api/internal/auth/` (env-gated, disabled by default) |
| Avatar/file upload + serving | `platform/httpserver/uploads.go`, `auth/application/avatar.go` |
| DB open/migrate/seed, sqlc queries | `platform/db/`, `db/queries/*.sql` |
| Errors, config, http helpers, ids | `internal/shared/{apperr,config,httpx,id}/` |
| API calls, CSRF retry, cookies, webp | `apps/web/src/shared/lib/{api,cookies,image}.ts` |
| Auth state, route guard, layout | `features/auth/AuthProvider.tsx`, `shared/ui/{RequireAuth,Layout}.tsx` |
| API contract | `packages/api-contract/openapi.yaml` (extend, don't abandon) |

## 2. Protected vs editable

NEVER delete, blank, or replace without explicit user approval (even if the
PRD doesn't mention them — they are infrastructure, not demo content):

- `apps/api/cmd/api/`, `apps/api/internal/platform/`, `apps/api/internal/shared/`
- `apps/api/internal/auth/` (all four layers), `db/queries/users.sql`, `oidc.sql`
- `apps/web/src/shared/`, `apps/web/src/features/auth/`, `apps/web/src/app/`
- `packages/api-contract/openapi.yaml`, `Makefile`, `sqlc.yaml`, `vite.config.ts`,
  `.env.example` keys (extend, don't remove)

Reference implementation (copy the pattern, not the content):

- `apps/api/internal/order/` and `apps/web/src/features/order/` are the
  canonical templates for a new context/feature. If the PRD has no "orders",
  leave `order` in place until its replacement works AND the user approves
  removal — never delete it as step one.

## Commands

- `make setup` deps, `make dev` (API :3000 + web :5173), `make build` → `./bin/api`
- `make test` backend tests, `make vet` (gofmt + go vet + tsc), `make sqlc` after query edits
- `make smoke` probes a running :3000 (CSRF cookie, CSP header, 403 without token)

## 3. Backend (`apps/api`)

- Layout: `internal/<context>/{domain,application,infrastructure,presentation}`; `shared/` and `platform/` only for cross-cutting code.
- `domain` is pure: no fiber/sql/os imports. Interfaces live in `domain`, implemented in `infrastructure`.
- Errors: `apperr.Invalid(field, msg)` → 400, `ErrOrderNotFound` → 404, `ErrNotOpen` → 409, bad credentials → 401. Never leak internal errors (generic 500).
- SQL must run on **both** sqlite and postgres: TEXT ids, BIGINT ints, RFC3339 TEXT timestamps, `$1` placeholders. Schema auto-migrates on boot.
- Security invariants (covered by `server_test.go`/`oidc_test.go`/`uploads_test.go`, keep them green): HttpOnly `session` cookie, CSRF header on unsafe `/api/*` methods, CSP on API + static responses. OIDC/LDAP stay disabled by default; external users carry empty password hashes and can never password-login. Uploads validate webp magic + size server-side; `/uploads/*` serving is allowlisted (no raw path join) on both storage drivers.

## 4. Frontend (`apps/web`)

- Layout: `src/features/<name>/{api.ts, pages/, components/}`, shared code in `src/shared/`. Features never import each other.
- All API I/O goes through `shared/lib/api.ts` (cookies + CSRF + one retry). Validate at boundaries with zod. No `any`, no localStorage tokens.
- TS is strict (`noUncheckedIndexedAccess`, `exactOptionalPropertyTypes`): fix type errors, don't loosen config.

### Design skills (read before doing the work)

- Build/change frontend UI: read the `ui-ux-pro-max` skill first.
- Logo/brand-mark creation: use the `logo-design` skill.
- Visual polish pass: use the `emil-design-eng` skill.

## Before commit

- `make vet && make test` must pass. Pre-existing tests stay green unless the
  user explicitly approved the behavior change.
- `git status` must not show mass deletions of §2 paths. Keep `.env`, `bin/`,
  `dist/`, `*.db` out of git (already ignored).
- Match the neighboring file's pattern; don't add deps or layers the task didn't ask for.
