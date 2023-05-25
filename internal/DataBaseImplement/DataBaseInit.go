package DataBaseImplement

import (
	"ComputerShopServer/internal/DataBaseImplement/Config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg Config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.PgHost, cfg.PgUser, cfg.PgPwd, cfg.PgDBName, cfg.PgPort,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func formatConnect(cfg Config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PgUser, cfg.PgPwd, cfg.PgHost, cfg.PgPort, cfg.PgDBName,
	)
}
