# KyokusuLib

Платформа для чтения ранобэ и манги. SSR-фронтенд на Nuxt (Bun), REST API на Go, PostgreSQL.

<p align="center">
    <img src="./frontend/app/assets/images/Kyokusu/maskot-transparent.png" alt="Screenshot">
</p>


---

## Стек

| Слой | Технологии |
|------|------------|
| Фронтенд | Nuxt 4 (Vue 3, SSR), Bun, Tailwind CSS 4, Pinia, VueUse |
| Бэкенд | Go 1.26, gorilla/mux, uber-go/fx (DI), godotenv |
| База | PostgreSQL 15, pgx/v5 |
| Миграции | pressly/goose |
| Инфра | Docker Compose, Caddy (reverse proxy + HTTPS) |

---

## Функциональность

| Раздел | Возможности |
|--------|-------------|
| Новеллы | Каталог, тома, главы, изображения, описание, жанры, категории |
| Чтение | Прогресс чтения, закладки (5 категорий), лайки, рейтинги |
| Аккаунт | Регистрация, логин по email, OAuth (Google/Discord), профиль, настройки |
| Сообщество | Комментарии к новеллам, лайки комментариев, репорты |
| Команды | Издательские команды, участники, заявки на вступление, роли |
| Уведомления | SSE (Server-Sent Events) в реальном времени |
| Парсер | Импорт ранобэ с ranobehub.org |
| Роли | user → publisher → moderator → admin (иерархия прав) |

---

## Разработка

### Быстрый старт (Docker)

```bash
cp .env.example .env     # настроить секреты
docker compose up --build -d
```

Фронтенд: `http://localhost:3000`
API: `http://localhost:8080`
Парсер: `http://localhost:3005`

### Локально (без Docker)

Только база в Docker, код с hot-reload:

```bash
docker compose up db -d

# Бэкенд
cd backend && go mod download
go install github.com/air-verse/air@latest && air

# Фронтенд
cd frontend && bun install && bun run dev
```

---

## Переменные окружения

Проект использует два `.env`:

- `backend/.env` — локальный запуск бэкенда (godotenv)
- `.env` (корень) — Docker Compose (подхватывается автоматически)

| Переменная | Описание |
|-----------|----------|
| `DATABASE_URL` | Подключение к PostgreSQL |
| `JWT_SECRET` | Секрет для JWT-токенов |
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` | OAuth Google |
| `DISCORD_CLIENT_ID` / `DISCORD_CLIENT_SECRET` | OAuth Discord |
| `KYOKUSU_EMAIL_NAME` / `KYOKUSU_EMAIL_PASSWORD` | SMTP для рассылки |
| `FRONTEND_URL` | URL фронтенда (CORS) |
| `BACKEND_URL` | URL бэкенда (редиректы) |
| `NUXT_PUBLIC_API_BASE` | Базовый URL API (браузер) |
| `PARSER_URL` | URL парсера (бэкенд → парсер) |

---

## Production

### 1. Запуск

```bash
git clone <repo> && cd KyokusuLib
cp .env.example .env
# Заменить localhost на домен в OAuth callback'ах и NUXT_PUBLIC_API_BASE
docker compose up --build -d
```

### 2. Reverse Proxy (Caddy)

В корне проекта `Caddyfile` — авто-HTTPS через Let's Encrypt:

```bash
sudo apt install caddy
sudo cp Caddyfile /etc/caddy/Caddyfile   # заменить домен
sudo systemctl restart caddy
```

Caddy маршрутизирует:
- `/api/*`, `/uploads/*` → бэкенд (порт 8080)
- `/parser/*` → парсер (порт 3005)
- остальное → Nuxt SSR (порт 3000)

### 3. OAuth callback'и

В Google Cloud Console и Discord Developer Portal:

```
https://<домен>/api/auth/google/callback
https://<домен>/api/auth/discord/callback
https://<домен>/api/socials/google/callback
https://<домен>/api/socials/discord/callback
```

---

## Структура проекта

```
├── .env.example           # Шаблон переменных
├── Caddyfile              # Production reverse proxy
├── docker-compose.yml     # db + backend + frontend + parser
│
├── apps/
│   └── ranobehub-parser/  # Парсер ранобэ (Bun, порт 3005)
│       ├── src/           # handlers, parsers, mappers
│       └── Dockerfile     # Multi-stage (deps → runner)
│
├── frontend/
│   ├── app/
│   │   ├── assets/        # Стили, изображения
│   │   ├── components/    # Vue-компоненты
│   │   ├── composables/   # API-клиенты, бизнес-логика
│   │   ├── constants/     # Константы, утилиты
│   │   ├── layouts/       # Шаблоны страниц
│   │   ├── middleware/    # Auth, router guards
│   │   ├── pages/         # Роутинг Nuxt
│   │   ├── stores/        # Pinia (auth, settings, notifications)
│   │   ├── types/         # TypeScript-типы
│   │   └── utils/         # Вспомогательные функции
│   ├── public/            # Статические файлы
│   └── Dockerfile         # Multi-stage (deps → build → runner)
│
└── backend/
    ├── cmd/api/main.go    # Точка входа (uber-go/fx DI)
    ├── internal/
    │   ├── config/        # Загрузка переменных
    │   ├── handlers/      # HTTP-обработчики
    │   ├── middleware/    # CORS, JWT, auth
    │   ├── repository/    # SQL-запросы (pgx)
    │   ├── routes/        # Регистрация роутов
    │   ├── services/      # Бизнес-логика
    │   └── parse/         # Сервис парсинга ranobehub
    ├── migrations/        # SQL-миграции (goose)
    ├── uploads/           # Загруженные файлы
    └── Dockerfile         # Multi-stage (deps → builder → scratch)
```

---

## API

REST API, базовый URL: `http://localhost:8080/api`

| Эндпоинт | Описание |
|----------|----------|
| `/api/auth/` | Регистрация, логин, OAuth (Google/Discord), восстановление пароля |
| `/api/socials/` | Привязка социальных сетей |
| `/api/users/` | Профили, настройки, уровень, активность |
| `/api/novelas/` | Каталог, тома, главы, контент, рейтинги, закладки, лайки |
| `/api/authors/` | Авторы |
| `/api/teams/` | Издательские команды, участники, заявки |
| `/api/comments/` | Комментарии, лайки, репорты |
| `/api/notifications/` | Уведомления (SSE) |
| `/api/moderation/` | Модерация контента |
| `/api/parse/ranobehub/{id}` | Импорт новеллы с ranobehub |
| `/api/health` | Health check |
| `/uploads/` | Статические файлы |
| `/parser/` | Парсер ranobehub (поиск, новеллы, главы) |

---