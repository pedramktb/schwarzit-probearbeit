package postgres

import (
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
)

func create(postgresURI string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{})
	if err != nil {
		panic(errors.Wrap(err, "failed to connect to database"))
	}

	return db
}

func Close(db *gorm.DB) {
	if db != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}
