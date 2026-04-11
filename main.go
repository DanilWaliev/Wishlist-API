package main

import (
	"fmt"
	"log"

	"github.com/DanilWaliev/wishlist-api/config"
	"github.com/DanilWaliev/wishlist-api/migrations"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config := config.Load()
	if err := migrations.Run(config.DBURL(), "migrations"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("миграции запущены успешно")
}
