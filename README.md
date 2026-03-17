# Camp Booking REST API

Backend-сервис для бронирования детских лагерей.

## Функциональность

| Регистрация и авторизация пользователей (JWT)
| Управление профилями детей
| Просмотр лагерей
| Добавление лагерей в избранное
| Бронирование лагеря для ребёнка
| Middleware (логирование, recovery, auth)
| Docker-окружение
| Управление миграциями базы данных

---

## Стек технологий

- Go
- PostgreSQL
- JWT
- Docker
- golang-migrate

---

## Основные endpoint'ы

### Auth

- POST `/auth/register`
- POST `/auth/login`
- GET `/auth/me`

### Children

- POST `/children`
- GET `/children`
- PUT `/children/{id}`

### Camps

- GET `/camps`
- GET `/camps/{id}`

### Favorites

- POST `/favorites/{campId}`
- DELETE `/favorites/{campId}`
- GET `/favorites`

### Bookings

- POST `/bookings`
- GET `/bookings`

---

## ⚙️ Запуск проекта (Docker)

### 1. Поднять базу данных

```bash
docker compose up -d db
```

### 2. Применить миграции

```bash
docker compose run --rm migrate up
```

### 3. Запустить приложение

```bash
docker compose up --build app
```

---

## Управление миграциями

Применить миграции:

```bash
docker compose run --rm migrate up
```

Откатить одну миграцию:

```bash
docker compose run --rm migrate down 1
```

Проверить версию:

```bash
docker compose run --rm migrate version
```

---
