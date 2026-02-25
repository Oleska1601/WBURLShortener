package app

import (
	"fmt"

	"github.com/Oleska1601/WBURLShortener/config"
	"github.com/wb-go/wbf/dbpg"
)

func initDB(cfg *config.PostgresConfig) (*dbpg.DB, error) {
	masterDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	slavesDSN := []string{}
	options := &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}
	db, err := dbpg.New(masterDSN, slavesDSN, options)
	if err != nil {
		return nil, fmt.Errorf("create new DB instance: %w", err)
	}

	return db, nil
}
