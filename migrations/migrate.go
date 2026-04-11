package migrations

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func Run(dbURL string, path string) error {
	m, err := migrate.New(
		"file://"+path,
		dbURL,
	)

	if err != nil {
		return fmt.Errorf("ошибка при создании миграции: %w\n", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка при запуске миграции: %w\n", err)
	}

	return nil
}
