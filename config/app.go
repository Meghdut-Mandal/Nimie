package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
	//rdb *redis.Client
)

func Connect() {
	d, err := gorm.Open(sqlite.Open("nimie.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d
}

func GetSqlDB() *gorm.DB {
	return db
}
