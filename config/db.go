package config

import (
	"gin-full-rest/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	var err error
	dsn := "host=localhost user=postgres password=psql123 dbname=learngin port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	DB.AutoMigrate(&models.Article{})
	DB.AutoMigrate(&models.User{})
	DB.Migrator().CreateConstraint(&models.User{}, "Articles")
	DB.Migrator().CreateConstraint(&models.User{}, "fk_users_articles")

}
