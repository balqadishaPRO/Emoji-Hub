# Emoji-Hub

REST API for browsing, searching, and favoriting emojis. Built with Go, Gin, and PostgreSQL. Data is seeded from the public [EmojiHub API](https://emojihub.yurace.pro).

**Live API:** [https://emoji-hub-6odk.onrender.com/api/emoji](https://emoji-hub-6odk.onrender.com/api/emoji)  
**Frontend:** [Emoji-Hub-Frontend](https://github.com/balqadishaPRO/Emoji-Hub-Frontend) · [Demo](https://balqadishapro.github.io/Emoji-Hub-Frontend/)

---

## Features

- List and filter emojis by name, category, and group
- Sort results by name or category
- Emoji detail view with cached “mood” enrichment
- Session-based favorites via HttpOnly cookie (`sid`)
- One-shot importer that upserts the full upstream catalog into PostgreSQL
- CORS configured for the GitHub Pages frontend and local development

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.24.1 |
| HTTP | [Gin](https://github.com/gin-gonic/gin) v1.10 |
| Database | PostgreSQL 16 |
| Driver | [lib/pq](https://github.com/lib/pq) |
| Config | [godotenv](https://github.com/joho/godotenv) |
| Deploy | Docker · [Render](https://render.com) |

---

## Architecture

```
cmd/
  api/          → HTTP server entrypoint
  importer/     → CLI that fetches & seeds emoji data
internal/
  handler/      → HTTP handlers & route registration
  service/      → Business logic & upstream fetch client
  repo/         → PostgreSQL access
  middleware/   → Session cookie (`sid`)
  model/        → Domain types
  llm/          → Mood generation (stub; cached in DB)
migrations/     → SQL schema up/down
```

Request flow: **handler → service → repo → PostgreSQL**

---

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.24+
- [Docker](https://www.docker.com/) & Docker Compose (for local Postgres)
- `psql` or any SQL client (to apply migrations)

### 1. Clone and install dependencies

```bash
git clone https://github.com/balqadishaPRO/Emoji-Hub.git
cd Emoji-Hub
go mod download
```

### 2. Start PostgreSQL

```bash
docker compose up -d
```

This starts:

| Service | URL / Port | Credentials |
|---------|------------|-------------|
| PostgreSQL | `localhost:5432` | user / password / db: `emoji` |
| pgAdmin | [http://localhost:5050](http://localhost:5050) | `admin@local` / `admin` |

### 3. Apply migrations

```bash
psql "postgres://emoji:emoji@localhost:5432/emoji?sslmode=disable" \
  -f migrations/000001_init_schema.up.sql
```

### 4. Configure environment

Create a `.env` in the project root:

```env
DATABASE_URL=postgres://emoji:emoji@localhost:5432/emoji?sslmode=disable
PORT=8080
```

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | — | PostgreSQL connection string |
| `PORT` | No | `8080` | HTTP listen port |

### 5. Seed emoji data

```bash
go run ./cmd/importer
```

Fetches all emojis from `https://emojihub.yurace.pro/api/all` and upserts them into the database.

### 6. Run the API

```bash
go run ./cmd/api
```

Server listens on `http://localhost:8080`.

---

## API Reference

Base path: `/api`  
All routes set a session cookie `sid` (UUID, 7-day lifetime, HttpOnly).

### List emojis

```http
GET /api/emoji?search=&category=&group=&sort=name
```

| Query | Description |
|-------|-------------|
| `search` | Filter by name (substring) |
| `category` | Filter by category |
| `group` | Filter by group |
| `sort` | `name` (default) or `category` |

**Response** `200` — array of emoji objects:

```json
[
  {
    "id": "...",
    "name": "grinning face",
    "category": "smileys and people",
    "group": "face positive",
    "htmlCode": ["&#128512;"],
    "unicode": ["U+1F600"]
  }
]
```

### Get emoji detail

```http
GET /api/emoji/:id
```

**Response** `200` — emoji plus `mood` (from cache or generated stub):

```json
{
  "id": "...",
  "name": "grinning face",
  "category": "smileys and people",
  "group": "face positive",
  "htmlCode": ["&#128512;"],
  "unicode": ["U+1F600"],
  "mood": "Cheerful and friendly vibes!"
}
```

**Response** `404` — `{"error":"..."}`

### List favorites

```http
GET /api/favorites
```

Uses the `sid` cookie. Returns `[]` when empty.

### Add favorite

```http
POST /api/favorites/:id
```

**Response** `204 No Content`

### Remove favorite

```http
DELETE /api/favorites/:id
```

**Response** `204 No Content`

### CORS

Allowed origins:

- `https://balqadishapro.github.io`
- `http://localhost:5173`

Credentials are enabled (`AllowCredentials: true`). Clients must send `credentials: 'include'` with fetch/XHR.

---

## Database Schema

| Table | Purpose |
|-------|---------|
| `emoji` | Catalog (`id`, `name`, `category`, `group`, `html_code`, `unicode`) |
| `favorites` | Session ↔ emoji pairs (`session_id`, `emoji_id`) |
| `llm_cache` | Cached mood strings per emoji |

Rollback: `migrations/000001_init_schema.down.sql`

---

## Docker

Build and run the API image (Postgres must be reachable via `DATABASE_URL`):

```bash
docker build -t emoji-hub .
docker run --rm \
  -e DATABASE_URL="postgres://emoji:emoji@host.docker.internal:5432/emoji?sslmode=disable" \
  -e PORT=8080 \
  -p 8080:8080 \
  emoji-hub
```

> `docker-compose.yml` runs **Postgres + pgAdmin only**. The API is started with `go run` or the Docker image above.

---

## Project Layout

```
.
├── cmd/
│   ├── api/main.go           # API server
│   └── importer/main.go      # Data importer CLI
├── internal/
│   ├── handler/              # Gin handlers
│   ├── service/              # Business logic & fetch client
│   ├── repo/                 # DB repository
│   ├── middleware/           # Session middleware
│   ├── model/                # Emoji / EmojiDetail types
│   └── llm/                  # Mood stub
├── migrations/               # Schema migrations
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## Related Repositories

| Repo | Role |
|------|------|
| [Emoji-Hub](https://github.com/balqadishaPRO/Emoji-Hub) | This backend API |
| [Emoji-Hub-Frontend](https://github.com/balqadishaPRO/Emoji-Hub-Frontend) | Static frontend (HTML/JS/CSS) on GitHub Pages |

---

## Notes

- Mood generation in `internal/llm` is currently a stub; results are cached in `llm_cache`.
- The production session cookie domain is pinned to the Render host; local cross-origin favorites may need that setting adjusted for your environment.
- Render Free Tier may cold-start after idle periods (~15 minutes).

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes
4. Open a pull request against `master`

Bug reports and feature ideas are welcome via GitHub Issues.

---

## License

No license file is included in this repository. All rights reserved by the author unless otherwise stated.
