package main

import (
	"gin-full-rest/config"
	"gin-full-rest/midleware"
	"gin-full-rest/routes"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDb()

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", routes.GetHome)

		v1.GET("/profile", midleware.IsAuth(), routes.Profile)
		auth := v1.Group("/auth")
		{
			auth.GET("/:provider", routes.RedirectHandler)
			auth.GET("/:provider/calback", routes.CallbackHandler)

		}

		articles := v1.Group("/articles")
		{
			articles.GET("/", routes.GetArticles)
			articles.GET("/tag/:tag", routes.GetArticleByTag)
			articles.PUT("/update/:id", midleware.IsAuth(), routes.UpdateArticle)
			articles.POST("/create", midleware.IsAuth(), routes.PostArticle)
			articles.GET("/:slug", routes.GetArticle)
			articles.DELETE("/delete/:id", midleware.IsAdmin(), routes.DeleteArticle)

		}

	}

	router.Run(":9000")
}
