# 🔥 Development Setup с Hot Reload

## Быстрый старт

### 1. Запуск в режиме разработки (с hot reload):

```bash
docker compose -f docker-compose.dev.yml up --build
```

**Что это дает:**
- ✅ **Backend (Go)**: автоматически пересобирается при изменении `.go` файлов (благодаря Air)
- ✅ **Frontend (Vue 3)**: автоматически обновляется в браузере при изменениях (Vite HMR)
- ✅ Все изменения применяются **моментально**, без пересборки контейнеров

### 2. Применить миграции:

```bash
docker compose -f docker-compose.dev.yml --profile tools run --rm migrate
```

### 3. Остановить все:

```bash
docker compose -f docker-compose.dev.yml down
```

---

## Production режим (без hot reload):

```bash
# Обычный запуск для продакшена
docker compose up -d

# Миграции
docker compose --profile tools run --rm migrate
```

---

## Как работает Hot Reload

### Backend (Go + Air):

Air следит за изменениями в `.go` файлах и автоматически:
1. Пересобирает бинарник
2. Перезапускает сервер
3. Среднее время: **~2-3 секунды**

**Файлы, которые отслеживаются:**
- `cmd/**/*.go`
- `internal/**/*.go`
- Все `.go` файлы кроме `*_test.go`

**Что игнорируется:**
- `frontend/`
- `tmp/`
- `vendor/`
- `migrations/`

### Frontend (Vue 3 + Vite HMR):

Vite предоставляет Hot Module Replacement (HMR):
1. Отслеживает изменения в `frontend/src/`
2. Обновляет только измененные модули
3. **Мгновенно** (< 50ms)

---

## Структура файлов

```
.
├── docker-compose.yml          # Production
├── docker-compose.dev.yml      # Development с hot reload
├── Dockerfile                  # Production image
├── Dockerfile.dev              # Development image с Air
├── .air.toml                   # Конфигурация Air
└── frontend/
    ├── Dockerfile              # Production build (multi-stage с nginx)
    ├── Dockerfile.dev          # Development с Vite dev server
    ├── nginx.conf              # Nginx конфиг для production
    └── src/                    # Hot reload работает здесь
```

---

## Полезные команды

```bash
# Просмотр логов
docker compose -f docker-compose.dev.yml logs -f app
docker compose -f docker-compose.dev.yml logs -f frontend

# Перезапустить только backend
docker compose -f docker-compose.dev.yml restart app

# Пересобрать только frontend
docker compose -f docker-compose.dev.yml up -d --build frontend

# Запустить тесты
docker compose -f docker-compose.dev.yml exec app go test ./...
```

---

## Troubleshooting

### Backend не пересобирается при изменениях:

```bash
# Проверьте логи Air
docker compose -f docker-compose.dev.yml logs app

# Перезапустите контейнер
docker compose -f docker-compose.dev.yml restart app
```

### Frontend не обновляется:

```bash
# Проверьте, что volume смонтирован правильно
docker compose -f docker-compose.dev.yml exec frontend ls -la /app/src

# Перезапустите Vite
docker compose -f docker-compose.dev.yml restart frontend
```

### Порт уже занят:

```bash
# Найдите процесс на порту 8080
lsof -ti:8080

# Остановите conflicting контейнер
docker stop <container_name>
```

---

## Сравнение режимов

| Функция | Production | Development |
|---------|-----------|-------------|
| **Backend hot reload** | ❌ Нет | ✅ Да (Air) |
| **Frontend hot reload** | ❌ Нет | ✅ Да (Vite HMR) |
| **Размер образа** | ~15 MB | ~500 MB |
| **Время запуска** | ~2 сек | ~5 сек |
| **Использование** | Продакшен, CI/CD | Локальная разработка |

---

## Доступ

- **Frontend**: http:
- **Backend API**: http:
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

---

🎉 **Готово!** Теперь любые изменения в коде применяются автоматически!
