package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DanilWaliev/wishlist-api/internal/config"
	"github.com/DanilWaliev/wishlist-api/internal/handlers"
	"github.com/DanilWaliev/wishlist-api/internal/middleware"
	"github.com/DanilWaliev/wishlist-api/internal/migrations"
	"github.com/DanilWaliev/wishlist-api/internal/repository"
	"github.com/DanilWaliev/wishlist-api/internal/services"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// загрузка конфига
	config, err := config.Load()
	if err != nil {
		log.Fatal("ошибка при загрузке конфига", err)
	}
	log.Println("конфиг загружен успешно")

	// подключение к БД
	db, err := OpenDB(config.DSN())
	if err != nil {
		log.Fatal("ошибка при подключении к БД", err)
	}
	log.Println("подключение к БД успешно")

	// запуск миграций
	if err := migrations.Run(config.DBURL(), "internal/migrations"); err != nil {
		log.Fatal(err)
	}
	log.Println("миграции прошли успешно")

	// инициализация зависимостей
	authRepo := repository.NewUserRepo(db)
	wishlistsRepo := repository.NewWishlistsRepo(db)

	authService := services.NewAuthService(authRepo, []byte(config.SecretKey))
	wishlistsService := services.NewWishlistsService(wishlistsRepo)

	authMW := middleware.NewAuthMiddleware(authService)
	authHandler := handlers.NewAuthHandler(authService)
	wishlistsHandler := handlers.NewWishlistsHandler(wishlistsService)

	// сбор всех обработчиков в контейнер, передача в роутер и получение mux
	h := NewHTTPHandler(authMW, authHandler, wishlistsHandler)
	mux := routes(h)

	// иницилизация структуры сервера
	server := &http.Server{
		Addr:              ":" + config.AppPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// запуск сервера
	log.Printf("Запуск сервера на %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к БД: %w\n", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка при подключении к БД: %w\n", err)
	}
	return db, err
}
