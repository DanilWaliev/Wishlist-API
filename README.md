# Wishlist API

## О проекте

Wishlist API - REST-сервис для создания и управления вишлистами к праздникам и событиям. Пользователь может зарегистрироваться, создать список подарков, добавить в него позиции и поделиться публичной ссылкой. По этой ссылке другие пользователи могут просматривать вишлист и бронировать подарки без авторизации.


---

## Запуск проекта

Для запуска потребуется установленный Docker и Docker Compose.

### 1. Клонировать репозиторий

    git clone https://github.com/DanilWaliev/Wishlist-API.git
    cd Wishlist-API

### 2. Создать `.env`

    cp .env.example .env

При необходимости отредактируйте значения в `.env`.

### 3. Запустить проект

    docker-compose up --build

---

## Переменные окружения

Пример `.env.example`:

    COMPOSE_PROJECT_NAME=wishlist

    APP_PORT=8080

    DB_HOST=db
    DB_PORT=5432
    DB_NAME=wishlist_api
    DB_USER=wishlist_user
    DB_PASSWORD=wishlist_password

    SECRET_KEY=your_secret_key

---

## Эндпоинты

### Авторизация

- POST /register  
- POST /login  

### Вишлисты

- GET /wishlists  
- POST /wishlists  
- GET /wishlists/{id}  
- PUT /wishlists/{id}  
- DELETE /wishlists/{id}  

### Позиции

- POST /wishlists/{id}/items  
- GET /wishlists/{id}/items  
- PUT /wishlists/{id}/items/{itemId}  
- DELETE /wishlists/{id}/items/{itemId}  

### Публичные эндпоинты

- GET /public/{token}  
- POST /public/{token}/reserve/{itemId}  

---

## Примеры запросов

### Регистрация

    POST /register
    Content-Type: application/json

    {
      "email": "user@example.com",
      "password": "password123"
    }

---

### Вход

    POST /login
    Content-Type: application/json

    {
      "email": "user@example.com",
      "password": "password123"
    }

---

### Создание вишлиста

    POST /wishlists
    Authorization: Bearer <token>
    Content-Type: application/json

    {
      "event_name": "День рождения",
      "description": "Список подарков",
      "event_date": "2026-05-10"
    }

---

### Добавление позиции

    POST /wishlists/1/items
    Authorization: Bearer <token>
    Content-Type: application/json

    {
      "title": "Наушники",
      "description": "Беспроводные",
      "product_url": "https://example.com",
      "priority": 1
    }

---

### Публичный просмотр

    GET /public/{token}

---

### Бронирование

    POST /public/{token}/reserve/1