package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var Once sync.Once

type Db struct {
	Db *gorm.DB
}

func New() *Db {
	return &Db{}
}

func (d *Db) Connect() error {

	dia := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		"postgres",
		"sunil",
		"password",
		"postgresdb",
		"5432",
	)
	db, err := gorm.Open(postgres.Open(dia), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return err
	}

	Once.Do(func() {
		d.Db = db
	})

	return nil
}
