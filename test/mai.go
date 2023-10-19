package main

import (
	"HiChat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/hiChat?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		&models.UserBasic{}, 
		&models.Relation{},
		&models.Community{},
		&models.Message{},
	)
	if err != nil {
		panic(err)
	}

}
