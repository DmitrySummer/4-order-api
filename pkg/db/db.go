package db

import (
	"4-order-api/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *config.Config) (*Db, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Db.Host,
		conf.Db.Port,
		conf.Db.User,
		conf.Db.Password,
		conf.Db.DbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к базе данных: %w", err)
	}

	return &Db{DB: db}, nil
}

func (db *Db) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("Не удалось получить sql.DB: %w", err)
	}
	return sqlDB.Close()
}
