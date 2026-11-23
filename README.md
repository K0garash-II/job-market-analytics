# Job Market Analytics

Аналитика IT-вакансий из Remote OK и Хабр Карьеры. React-интерфейс показывает динамику публикаций за 7 и 30 дней, распределение по источникам и формату работы, востребованные навыки и рейтинг работодателей. В каталоге доступны поиск по названию или компании, фильтры и переход к исходной публикации.

## Архитектура

```text
Remote OK / Хабр Карьера → Collector → NATS JetStream → Processor → PostgreSQL
                             ↕                                         ↑
                           Redis                               API ← React
```

- **Collector** запускает единый сценарий сбора для всех фидов. Адаптеры источников, Redis и очереди реализуют его интерфейсы. После подтверждения публикации в JetStream ключ вакансии попадает в Redis на 24 часа.
- **Processor** загружает словарь при старте и извлекает навыки из заголовка и описания. PostgreSQL сохраняет вакансию и связи в одной транзакции. Уникальность пары источника и внешнего ID позволяет повторить доставку при сбое. После успешной записи обработчик подтверждает сообщение.
- **API** принимает HTTP-запросы в `services/api/internal/handler`. Модели чтения и контракт находятся в `internal/market`, реализация запросов — в `internal/postgres`. Обзор и страницы каталога читаются в согласованном снимке базы.
- **React** получает данные по HTTP с того же origin. `frontend/src/api.js` отвечает за запросы и отмену при смене страницы. Компоненты разделов используют общие представления списков, рейтингов, ошибок и загрузки.

## Запуск

Нужны Go 1.26.1+, Node.js 22.12+, pnpm 11 и Docker Compose. Команды для PowerShell выполняются из корня репозитория.

```powershell
Copy-Item .env.example .env
docker compose -f infra/docker-compose.yml up -d --wait
```

На новой базе примените миграции по порядку. Если база уже существует, выполните только ещё не применённые файлы. Миграция 006 добавляет источник Хабр Карьеры.

```powershell
Get-ChildItem infra/migrations/*.up.sql | Sort-Object Name | ForEach-Object {
  Get-Content -Raw $_.FullName | docker exec -i jobmarket_pg psql -v ON_ERROR_STOP=1 -U app -d jobmarket
  if ($LASTEXITCODE) { throw "Migration failed: $($_.Name)" }
}
```

Запустите сервисы в отдельных терминалах. Сборщик выполняет один проход, API и обработчик остаются запущенными. При переходе со старой версии перезапустите оба сервиса очереди до следующего сбора.

```powershell
go run ./services/processor/cmd
go run ./services/collector/cmd
go run ./services/api/cmd
```

```powershell
pnpm --dir frontend install --frozen-lockfile
pnpm --dir frontend dev
```

Откройте [дашборд](http://localhost:5173). Vite передаёт запросы `/api` на `http://127.0.0.1:8080`. Для другого адреса задайте `API_TARGET` перед запуском Vite. Для размещения соберите `frontend/dist`, раздайте каталог статическим сервером и проксируйте путь `/api/` на Go-сервис.

## API

| GET /api/v1 | Данные |
|---|---|
| `/analytics/overview` | Общие показатели, источники и ежедневные публикации за 30 календарных дней. По умолчанию период заканчивается сегодня. Дашборд передаёт `period=latest`, чтобы показать данные по дату последней публикации |
| `/analytics/skills` | До 20 навыков и количество вакансий с каждым из них |
| `/analytics/companies` | До 10 работодателей по числу вакансий |
| `/analytics/salary` | Покрытие зарплатами и исторические средние границы без конвертации валют |
| `/vacancies` | Каталог и общее количество совпадений. Параметры `page`, `limit` (1–100), `q`, `location`, `remote`, `source` |
| `/vacancies/:id` | Подробности вакансии |
| `/vacancies/search?q=golang` | Полнотекстовый поиск по заголовку и описанию |

Доли рынка рассчитываются по собранным вакансиям. Источник hh.ru имеет отдельный адаптер и по умолчанию в сборе не участвует.

## Проверки

```powershell
go test ./...
go vet ./...
pnpm --dir frontend build
```

Для интеграционного теста атомарной записи задайте `TEST_DATABASE_URL` с адресом тестового PostgreSQL и выполните `go test -count=1 ./pkg/storage`. Тест использует временные таблицы своего соединения.

Сквозной тест `TestCollectorToProcessor` использует отдельные тестовые сервисы и оставляет контрольную вакансию в тестовой базе. Примените миграции, запустите Processor с адресами тестовых PostgreSQL и NATS, затем задайте `TEST_DATABASE_URL`, `TEST_NATS_URL` и `TEST_REDIS_URL` и выполните `go test -count=1 ./services/collector/internal/collect`.
