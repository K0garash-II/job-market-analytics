# Job Market Analytics

**English** · [Русский](README.ru.md)

Explore IT vacancies from Remote OK and Habr Career. The React dashboard tracks postings over 7 or 30 days, breaks them down by source and work arrangement, and ranks skills and employers. Search the catalog by job title or company, apply filters, and open the original posting.

## Architecture

```text
Remote OK / Habr Career → Collector → NATS JetStream → Processor → PostgreSQL
                             ↕                                         ↑
                           Redis                                  API ← React
```

- **Collector** runs the same ingestion workflow for each feed. Source, Redis, and queue adapters implement its interfaces. Once JetStream confirms publication, the collector caches the vacancy key in Redis for 24 hours.
- **Processor** loads a skill dictionary at startup and extracts skills from each job title and description. It saves the vacancy and its skill associations in one PostgreSQL transaction, then acknowledges the message. A unique constraint on the source and external ID prevents duplicates on redelivery.
- **API** handles HTTP requests in `services/api/internal/handler`. Within the API service, `internal/market` defines the read models and repository interface, and `internal/postgres` implements the queries. Each overview or catalog request reads from a consistent database snapshot.
- **React** fetches data from the same origin. `frontend/src/api.js` handles requests and cancels them when users navigate away. The dashboard sections share list, ranking, loading, and error components.

## Run locally

Install Go 1.26.1+, Node.js 22.12+, pnpm 11, and Docker Compose. Run the following PowerShell commands from the repository root.

```powershell
Copy-Item .env.example .env
docker compose -f infra/docker-compose.yml up -d --wait
```

For a new database, apply all migrations in order. For an existing database, run only the migrations that have not yet been applied. Migration 006 adds Habr Career as a source. Migration 007 adds the vacancy retention policy.

```powershell
Get-ChildItem infra/migrations/*.up.sql | Sort-Object Name | ForEach-Object {
  Get-Content -Raw $_.FullName | docker exec -i jobmarket_pg psql -v ON_ERROR_STOP=1 -U app -d jobmarket
  if ($LASTEXITCODE) { throw "Migration failed: $($_.Name)" }
}
```

Run each Go service in a separate terminal. The collector performs one pass and exits. The processor and API keep running. When upgrading from an older version, restart the processor and use the updated collector for the next run.

```powershell
go run ./services/processor/cmd
go run ./services/collector/cmd
go run ./services/api/cmd
```

The Compose `retention` service checks the database at startup and every minute thereafter. It deletes vacancies older than 720 hours from publication, along with their skill associations. PostgreSQL skips expired postings on reimport and retains postings whose publication date is unknown. To enable retention on an existing installation, apply migration 007 and run `docker compose -f infra/docker-compose.yml up -d --no-deps retention`.

Start the frontend in another terminal.

```powershell
pnpm --dir frontend install --frozen-lockfile
pnpm --dir frontend dev
```

Open the [dashboard](http://localhost:5173). Vite proxies `/api` requests to `http://127.0.0.1:8080`. Set `API_TARGET` before starting Vite to use a different backend. For deployment, build the frontend, serve `frontend/dist` with a static web server, and proxy `/api/` to the Go API.

## API

| GET /api/v1 | Response |
|---|---|
| `/analytics/overview` | Summary metrics, source counts, and daily postings over 30 calendar days. The window ends today by default. The dashboard passes `period=latest` to end the window on the most recent publication date. |
| `/analytics/skills` | Up to 20 skills and the number of vacancies mentioning each one. |
| `/analytics/companies` | Up to 10 employers ranked by vacancy count. |
| `/analytics/salary` | Salary coverage and average lower and upper salary bounds across stored postings, without currency conversion. |
| `/vacancies` | Paginated vacancies and the total number of matches. Accepts `page`, `limit` (1–100), `q`, `location`, `remote`, and `source`. |
| `/vacancies/:id` | Details for one vacancy. |
| `/vacancies/search?q=golang` | Full-text search across job titles and descriptions. |

Market shares are calculated from collected vacancies. The repository includes an hh.ru adapter, which is disabled by default.

## Checks

```powershell
go test ./...
go vet ./...
pnpm --dir frontend build
```

To test atomic writes, set `TEST_DATABASE_URL` to a test PostgreSQL database and run `go test -count=1 ./pkg/storage`. The test uses temporary tables scoped to its connection.

The end-to-end test `TestCollectorToProcessor` uses dedicated test services and leaves a test vacancy in the test database. Apply the migrations, start the processor against the test PostgreSQL and NATS instances, then set `TEST_DATABASE_URL`, `TEST_NATS_URL`, and `TEST_REDIS_URL` and run `go test -count=1 ./services/collector/internal/collect`.
