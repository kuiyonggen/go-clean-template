package app

import (
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
        "github.com/kuiyonggen/go-clean-template/config"
        "github.com/kuiyonggen/go-clean-template/migrations"
        "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func initMig(cfg *config.Config) {
        databaseURL := cfg.PG.URL + "?sslmode=disable"

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

        // wrap assets into Resource
        s := bindata.Resource(migrations.AssetNames(),
        func(name string) ([]byte, error) {
            return migrations.Asset(name)
        })

        d, err := bindata.WithInstance(s)
        if err != nil {
            log.Fatalf("failed create bindata instance, err: %s.", err)
            return
        }

	for attempts > 0 {
                m, err = migrate.NewWithSourceInstance("go-bindata", d, databaseURL)
                if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}
