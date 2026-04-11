package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DanilWaliev/wishlist-api/internal/config"
	"github.com/DanilWaliev/wishlist-api/internal/handlers"
	"github.com/DanilWaliev/wishlist-api/internal/middleware"
	"github.com/DanilWaliev/wishlist-api/internal/migrations"
	"github.com/DanilWaliev/wishlist-api/internal/repository"
	"github.com/DanilWaliev/wishlist-api/internal/services"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("ошибка при загрузке конфига: ", err)
	}
	log.Println("конфиг загружен успешно")

	db, err := OpenDB(cfg.DSN())
	if err != nil {
		log.Fatal("ошибка при подключении к БД: ", err)
	}
	log.Println("подключение к БД успешно")

	if err := migrations.Run(cfg.DBURL(), "internal/migrations"); err != nil {
		log.Fatal("ошибка при запуске миграций: ", err)
	}
	log.Println("миграции прошли успешно")

	authRepo := repository.NewUserRepo(db)
	wishlistsRepo := repository.NewWishlistsRepo(db)
	wlItemsRepo := repository.NewWishlistItemsRepo(db)

	authService := services.NewAuthService(authRepo, []byte(cfg.SecretKey))
	wishlistsService := services.NewWishlistsService(wishlistsRepo, wlItemsRepo)
	wlItemsService := services.NewWishlistItemsService(wlItemsRepo, wishlistsRepo)

	authMW := middleware.NewAuthMiddleware(authService)
	authHandler := handlers.NewAuthHandler(authService)
	wishlistsHandler := handlers.NewWishlistsHandler(wishlistsService)
	wlItemsHandler := handlers.NewWishlistItemsHandler(wlItemsService)

	h := NewHTTPHandler(authMW, authHandler, wishlistsHandler, wlItemsHandler)
	mux := routes(h)

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Запуск сервера на %s", server.Addr)

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("ошибка сервера: ", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("получен сигнал остановки, завершаем работу сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("ошибка graceful shutdown: ", err)

		if err := server.Close(); err != nil {
			log.Println("ошибка принудительного закрытия сервера: ", err)
		}
	}

	if err := db.Close(); err != nil {
		log.Println("ошибка при закрытии БД: ", err)
	}

	log.Println("сервер остановлен")
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к БД: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка при подключении к БД: %w", err)
	}

	return db, nil
}
