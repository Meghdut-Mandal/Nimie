package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Connect() {
	d, err := gorm.Open(sqlite.Open("nimie.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}
