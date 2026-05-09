# 🚀 Development Setup with Hot Reload

Complete development environment setup guide with automatic code reloading.

## Quick Start

### 1. Start development mode (with hot reload):

```bash
docker compose -f docker-compose.dev.yml up --build
```

**What you get:**
- ✅ **Backend (Go)**: Auto-rebuilds on `.go` file changes (powered by Air)
- ✅ **Frontend (Vue 3)**: Auto-reloads in browser on changes (Vite HMR)
- ✅ All changes apply **instantly** without rebuilding containers

### 2. Run database migrations:

```bash
docker compose -f docker-compose.dev.yml --profile tools run --rm migrate
```

### 3. Stop all services:

```bash
docker compose -f docker-compose.dev.yml down
```

---

## Production Mode (without hot reload)

```bash
# Standard production launch
docker compose up -d

# Run migrations
docker compose --profile tools run --rm migrate
```

---

## How Hot Reload Works

### Backend (Go + Air)

Air watches for changes in `.go` files and automatically:
1. Rebuilds the binary
2. Restarts the server
3. Average time: **~2-3 seconds**

**Watched files:**
- `cmd/**/*.go`
- `internal/**/*.go`
- All `.go` files except `*_test.go`

**Ignored:**
- `frontend/`
- `tmp/`
- `vendor/`
- `migrations/`

### Frontend (Vue 3 + Vite HMR)

Vite provides Hot Module Replacement (HMR):
1. Watches changes in `frontend/src/`
2. Updates only modified modules
3. **Instant** (< 50ms)

---

## Project Structure

```
.
├── docker-compose.yml          # Production setup
├── docker-compose.dev.yml      # Development with hot reload
├── Dockerfile                  # Production image
├── Dockerfile.dev              # Development image with Air
├── .air.toml                   # Air configuration
└── frontend/
    ├── Dockerfile              # Production build (multi-stage with nginx)
    ├── Dockerfile.dev          # Development with Vite dev server
    ├── nginx.conf              # Nginx config for production
    └── src/                    # Hot reload works here
```

---

## Useful Commands

```bash
# View logs
docker compose -f docker-compose.dev.yml logs -f app
docker compose -f docker-compose.dev.yml logs -f frontend

# Restart backend only
docker compose -f docker-compose.dev.yml restart app

# Rebuild frontend only
docker compose -f docker-compose.dev.yml up -d --build frontend

# Run tests
docker compose -f docker-compose.dev.yml exec app go test ./...

# Run linter
docker compose -f docker-compose.dev.yml exec app golangci-lint run ./...

# Access database
docker compose -f docker-compose.dev.yml exec postgres psql -U postgres -d anagrams

# Access Redis CLI
docker compose -f docker-compose.dev.yml exec redis redis-cli
```

---

## Troubleshooting

### Backend not rebuilding on changes:

```bash
# Check Air logs
docker compose -f docker-compose.dev.yml logs app

# Restart container
docker compose -f docker-compose.dev.yml restart app
```

### Frontend not updating:

```bash
# Verify volume is mounted correctly
docker compose -f docker-compose.dev.yml exec frontend ls -la /app/src

# Restart Vite
docker compose -f docker-compose.dev.yml restart frontend
```

### Port already in use:

```bash
# Find process using port 8080
lsof -ti:8080

# Stop conflicting container
docker stop <container_name>
```

### Clean rebuild (nuclear option):

```bash
# Stop and remove all containers, volumes, and networks
docker compose -f docker-compose.dev.yml down -v

# Rebuild from scratch
docker compose -f docker-compose.dev.yml up --build
```

---

## Mode Comparison

| Feature | Production | Development |
|---------|-----------|-------------|
| **Backend hot reload** | ❌ No | ✅ Yes (Air) |
| **Frontend hot reload** | ❌ No | ✅ Yes (Vite HMR) |
| **Image size** | ~15 MB | ~500 MB |
| **Startup time** | ~2 sec | ~5 sec |
| **Use case** | Production, CI/CD | Local development |

---

## Access Points

After running `docker compose -f docker-compose.dev.yml up`:

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

---

## Development Workflow

1. **Make changes** to backend (`internal/`, `cmd/`) or frontend (`frontend/src/`) code
2. **Save the file** - changes auto-apply:
   - Backend: Wait ~2-3 seconds for Air to rebuild
   - Frontend: Instant browser update via HMR
3. **Check logs** if something goes wrong:
   ```bash
   docker compose -f docker-compose.dev.yml logs -f
   ```

---

## Environment Variables

Create a `.env` file in the project root (optional - defaults are set):

```env
# Database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=anagrams

# Backend
DATABASE_URL=postgres://postgres:postgres@postgres:5432/anagrams?sslmode=disable
REDIS_URL=redis:6379
JWT_SECRET=your-secret-key-change-in-production
GOOGLE_CLIENT_ID=your-google-oauth-client-id
GOOGLE_CLIENT_SECRET=your-google-oauth-secret

# Frontend
VITE_API_URL=http://localhost:8080
```

---

🎉 **Ready!** All code changes now apply automatically!
