package main

import (
	"gin-full-rest/routes"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title string
	Slug  string
	Desc  string
}

var DB *gorm.DB

func maina() {
	var err error
	dsn := "host=localhost user=postgres password=psql123 dbname=learngin port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	DB.AutoMigrate(&Article{})

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", routes.GetHome)

		auth := v1.Group("/auth")
		{
			auth.GET("/:provider", routes.RedirectHandler)
			auth.GET("/:provider/calback", routes.CallbackHandler)

		}
		articles := v1.Group("/articles")
		{
			articles.GET("/", getArticles)
			articles.GET("/:slug", getArticle)
			articles.POST("/", postArticle)
		}

	}

	router.Run(":3000")
}

func getHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": "Welocme to gin",
	})
}

func getArticles(c *gin.Context) {
	data := []Article{}
	DB.Find(&data)
	c.JSON(200, gin.H{
		"data": data,
	})
}
func getArticle(c *gin.Context) {
	slug := c.Param("slug")
	var item Article
	err := DB.First(&item, "Slug = ?", slug).Error
	if err != nil {
		c.JSON(404, gin.H{"status": "error", "message": "record not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "Berhasil",
		"data":   item,
	})
}

func postArticle(c *gin.Context) {

	item := Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Slug:  slug.Make(c.PostForm("title")),
	}
	DB.Create(&item)
	c.JSON(200, gin.H{
		"message": "Berhasil",
		"data":    item,
	})
}
