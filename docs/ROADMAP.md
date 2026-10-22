# План развития KyokusuLib

Документ фиксирует приоритеты, критерии готовности и текущий статус. Порядок этапов — не «сделать всё», а «сначала то, что даёт максимальный эффект при минимальном риске».

## Базовые замеры

Машина: Ryzen 5 2600 (6c/12t), Go 1.27.2, Postgres 15, Redis 7 (локально).
Данные: 7 новелл, 834 главы, 1 TG-пост.

| Сценарий | До | После индекса |
|---|---:|---:|
| `GET /api/novela/7` (475 глав, 564 картинки) | 160 RPS / 124.7 мс | 618–710 RPS / 28–32 мс |
| SQL страницы новеллы | 33.8–39.7 мс | 7.8 мс |
| `Novela.GetFullByID` (Go) | 40.3 мс | 9.35 мс |
| `GET /api/health` | — | 73–80k RPS |
| `GET /api/tg/posts` | — | 25.8k RPS |
| `GET /api/novela?limit=20` | — | 2.4k RPS / 8.3 мс |
| Ресайз постера 1200×1600 | — | 51 мс / 21 MB |

| Размер ответа после 0.2 | Байт |
|---|---:|
| `GET /api/novela/7` (полный список глав, без картинок) | 67 467 |
| `GET /api/novela/7?summary=1` (страница новеллы) | 8 327 |
| `GET /api/novela/volumes/{id}/chapters` (активный том) | 1 730 |
| Итого первичная загрузка страницы | ~10 KB |

| `GetFullByID` после 0.2 | Время | Память | Аллокации |
|---|---:|---:|---:|
| Полный список глав, без картинок | 4.91 мс | 460 KB | 1238 |
| Summary для страницы | 1.41 мс | 31 KB | 172 |
| Главы одного тома | 0.21 мс | 9.3 KB | 101 |

---

## Этап 0. Производительность

| # | Задача | Зачем | Критерий готово | Статус |
|---|---|---|---|---|
| 0.1 | Индекс `novela_chapter_images(chapter_id, position)` | Убрать `Seq Scan` × 474 на странице новеллы | SQL < 10 мс, RPS ×3 | ✅ сделано |
| 0.2 | Разгрузить ответ страницы новеллы | 122 KB и 737 KB аллокаций на каждый просмотр; главы с картинками грузятся целиком | метаданные + список томов/глав < 40 KB; картинки — отдельным запросом | ✅ сделано: summary 8.3 KB + главы активного тома 1.7 KB; картинки только в `/api/novela/chapters/{id}` |
| 0.3 | Ресайз изображений в фоновую очередь | 34–51 мс и 15–21 MB блокируют HTTP-запрос | p95 upload < 200 мс, воркер масштабируется отдельно | план |
| 0.4 | Аудит индексов по `EXPLAIN ANALYZE` | Список/поиск/уведомления/активность | нет `Seq Scan` на горячих запросах | ✅ сделано: индексы и GIN добавлены; на dev-данных часть `Seq Scan` остаётся из-за размера таблиц |
| 0.5 | Cache-aside в Redis для каталога и новеллы | Redis используется только для pubsub и rate limit | cache hit > 70% на `GET /api/novela` | ✅ сделано: каталог и страница новеллы кэшируются per-user, инвалидация через generation-ключ |
| 0.6 | Cursor-пагинация лент | Offset деградирует на больших таблицах | TG/уведомления/активность на keyset | ✅ сделано: keyset по `(time, id)`, курсор в `X-Next-Cursor` |

### Замеры 0.4

Добавлены индексы: `notifications(user_id, created_at DESC)`,
`user_activities(user_id, timestamp DESC)`,
`user_novela_bookmarks(user_id, category_id, updated_at DESC)`,
`novela_genres(genre_id, novela_id)`,
`novela_categories(category_id, novela_id)`,
`novela(created_at DESC)`, `novela(views DESC)`,
частичный `novela_chapters(novela_volume_id, created_at DESC) WHERE status = 'approved'`
и GIN `pg_trgm` по выражению `novela_search_text(title, alternative_titles)`.

| Запрос | До | После |
|---|---:|---:|
| Уведомления, `user_id + created_at DESC` (тёплые) | 0.053–0.061 мс, Bitmap + Sort | 0.015–0.017 мс, Index Scan без Sort |
| Активность, `user_id=27` (форсированный индекс) | 0.111 мс, Seq Scan, 9 buffers | 0.052 мс, Bitmap Index Scan, 3 buffers |
| Каталог, сортировка `created_at DESC` | Seq Scan + Sort | `Index Scan idx_novela_created_at` |
| `last_chapter_at` / `chapters_count` | Seq Scan `novela_chapters` | `Index Only Scan idx_novela_chapters_volume_created_approved` |
| Поиск по названию/альтернативным названиям | Seq Scan, ветка `$1 ILIKE ANY(alternative_titles)` не работала | GIN `Bitmap Index Scan`; `GET /api/novela?search=monarch` находит «Я - Монарх» по `I am the Monarch` |

Данные в dev-БД маленькие (7 новелл, 573 уведомления, 159 активностей), поэтому
планировщик намеренно оставляет `Seq Scan` на таблицах в несколько страниц.
Проверка индексов выполнялась через `EXPLAIN` с `enable_seqscan=off`; на росте
данных выбираются индексы. Поиск переписан на `novela_search_text()`, иначе OR
с неиндексируемой веткой не давал использовать GIN.

### Замеры 0.5

Cache-aside для `GET /api/novela` и `GET /api/novela/{id}`. Ключи включают
`userID`, потому что ответ содержит закладку, лайк, рейтинг и прогресс чтения.
TTL: каталог 30 с, страница 60 с. Инвалидация — через общий generation-ключ
`novela:cache:version`, который инкрементируется при изменении новеллы, тома,
главы, закладки, лайка, рейтинга или прогресса чтения. Старые ключи доживают
свой TTL и не читаются.

`ab -n 1000 -c 20` для каталога и `ab -n 500 -c 10` для страницы:

| Сценарий | Без Redis | С Redis |
|---|---:|---:|
| `GET /api/novela?limit=20` | 1689 RPS / 11.8 мс | 7702 RPS / 2.6 мс |
| `GET /api/novela/7?summary=1` | 1575 RPS / 6.35 мс | 8589 RPS / 1.16 мс |

После прогрева 1500 запросов: `keyspace_hits` 3010, `keyspace_misses` 6;
в кэше один ключ каталога и два ключа страниц. Тест
`TestNovelaCacheHitAndInvalidate` проверяет miss → hit → инвалидацию для обоих
эндпоинтов.

### Замеры 0.6

Ленты переведены на keyset-пагинацию по `(created_at, id)` для TG-постов и
уведомлений и по `(timestamp, id)` для активности. На клиент возвращается
прежний массив, а позиция следующей страницы — в заголовках `X-Next-Cursor` и
`X-Has-More`. `offset` больше не используется: `offset=0` принимается для
совместимости, `offset>0` возвращает 400. Это убирает чтение и отбрасывание
всех предыдущих строк.

`EXPLAIN (ANALYZE, BUFFERS)` на dev-БД, уведомления пользователя 13:

| Запрос | План | Buffers | Время |
|---|---|---:|---:|
| `OFFSET 250 LIMIT 20` | Seq Scan + Sort, прочитано 290 строк | 37 | 0.379 мс |
| keyset `LIMIT 20` | Bitmap Index Scan по `idx_notifications_user_id_created_at_id` | 5 | 0.117 мс |

Активность пользователя 13:

| Запрос | План | Buffers | Время |
|---|---|---:|---:|
| `OFFSET 100 LIMIT 20` | Seq Scan + Sort, прочитано 137 строк | 6 | 0.239 мс |
| keyset `LIMIT 20` | Index Scan по `idx_user_activities_user_id_timestamp_id` | 3 | 0.162 мс |

На маленьких таблицах планировщик иногда выбирает Seq Scan; индексные планы
проверены с `enable_seqscan=off`. В новом индексе `id` добавлен как tie-breaker,
чтобы страницы не пересекались при одинаковых timestamp. Тесты
`internal/utils/cursor` и `internal/handlers/pagination_test.go` проверяют
кодирование курсора, лимиты и отказ от глубокого offset.

## Этап 1. Безопасность (P0)

| # | Задача | Файлы | Статус |
|---|---|---|---|
| 1.1 | Закрыть `PUT /api/novela/{id}`: auth + ownership/RBAC | `routes/novela.go`, `handlers/novela.go` | ✅ сделано |
| 1.2 | Cookie: `HttpOnly`, `Secure` в prod | `handlers/auth.go` | ✅ сделано |
| 1.3 | Fail-fast при дефолтных/пустых секретах | `config/config.go` | ✅ сделано |
| 1.4 | CORS из `cfg.FrontendURL`, а не хардкод localhost | `app/invokers.go` | ✅ сделано |
| 1.5 | Случайный OAuth `state` на попытку входа | `handlers/auth.go` | ✅ сделано |
| 1.6 | Ротация/отзыв JWT (короткий access + refresh) | `services/jwt_service.go` | план |
| 1.7 | Единая политика авторизации (RBAC + ownership) | routes/services | план |

## Этап 2. Ошибки и слои

| # | Задача | Статус |
|---|---|---|
| 2.1 | `apperrors`: доменная ошибка = код + HTTP-статус + сообщение | план |
| 2.2 | Error-middleware; убрать `err.Error()` из 500-ответов | план |
| 2.3 | Sentinel-ошибки вместо `errors.New("...")` | план |
| 2.4 | Транзакции — в repository (unit-of-work), не в сервисах | план |
| 2.5 | Интерфейсы на границах сервисов + моки | план |
| 2.6 | Единый auth-middleware (`required` флаг) вместо двух копий | план |
| 2.7 | `ctx` вместо `context.Background()` в сервисах | план |

## Этап 3. Качество и эксплуатация

| # | Задача | Статус |
|---|---|---|
| 3.1 | CI: `build`, `vet`, `golangci-lint`, `test`, бенчмарки | план |
| 3.2 | Первые тесты: auth, tg-delete, SSE, стратегии | план |
| 3.3 | Readiness отдельно от liveness (DB + Redis) | план |
| 3.4 | Request-id + единый `slog` | план |
| 3.5 | RED-метрики (Prometheus), затем OTel | план |
| 3.6 | Миграции/сиды отдельным job, не в старте приложения | план |
| 3.7 | Rate limit: trusted proxy / `X-Forwarded-For` | план |

## Этап 4. Масштабирование кода

| # | Задача | Статус |
|---|---|---|
| 4.1 | Распил god objects (`novela_repository` 1124 строки и др.) | план |
| 4.2 | sqlc или generics-хелперы для сканов | план |
| 4.3 | OpenAPI/контракт с фронтендом | план |
| 4.4 | Консистентная структура пакетов (feature-пакеты) | план |

---

## Метрики успеха

- `GET /api/novela/{id}`: p95 < 50 мс при 50 конкурентных, ответ < 40 KB.
- Загрузка изображения: p95 < 200 мс (обработка вне запроса).
- Ноль write-эндпоинтов, доступных без авторизации.
- Нет внутренних ошибок/SQL-деталей в HTTP-ответах.
- CI зелёный на каждый PR; покрытие `services` ≥ 30%.
- Readiness отражает состояние DB и Redis.
- `go test -bench` запускается в CI на фиксированном датасете.

## Порядок работ

1. 0.1 ✅ (индекс страницы новеллы).
2. 1.1–1.5 (закрыть дыры безопасности — P0).
3. 0.2 → 0.4 (дальше по производительности).
4. 2.1–2.4 (каркас ошибок и слоёв).
5. 3.1–3.4 (CI, тесты, наблюдаемость).
6. 0.3 и этап 4 — по мере роста нагрузки.
