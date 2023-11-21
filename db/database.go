package db

import (
	"admin-api/entity"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	name := os.Getenv("MYSQL_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return
		}
		sqlDB.Close()
	}
}

func GetDB() *gorm.DB {
	return db
}

func Migrate() error {
	return db.AutoMigrate(
		&entity.AdminUser{},
	)
}
