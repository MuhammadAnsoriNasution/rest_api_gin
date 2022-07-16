package routes

import (
	"fmt"
	"gin-full-rest/models"
	"strconv"
	"time"

	"gin-full-rest/config"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func GetHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": "Welocme to gin",
	})
}

func GetArticles(c *gin.Context) {
	data := []models.Article{}
	config.DB.Find(&data)
	c.JSON(200, gin.H{
		"data": data,
	})
}
func GetArticle(c *gin.Context) {
	slug := c.Param("slug")
	var item models.Article
	err := config.DB.First(&item, "Slug = ?", slug).Error
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

func Profile(c *gin.Context) {

	var users models.User
	user_id := uint(c.MustGet("jwt_user_id").(float64))
	data := config.DB.Where("id = ?", user_id).Preload("Articles", "user_id = ?", user_id).Find(&users)

	fmt.Println(data)
	// fmp
	c.JSON(200, gin.H{
		"status": "Berhasil",
		"data":   users,
	})
}
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var item models.Article
	err := config.DB.First(&item, "id = ?", id).Error
	if err != nil {
		c.JSON(404, gin.H{"status": "error", "message": "record not found"})
		c.Abort()
		return
	}
	config.DB.Model(&item).Where("id = ?", id).Updates(models.Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Tag:   c.PostForm("tag"),
	})

	if uint(c.MustGet("jwt_user_id").(float64)) != item.UserId {
		c.JSON(403, gin.H{"status": "error", "message": "forbiden"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"status": "Berhasil",
		"data":   item,
	})
}
func PostArticle(c *gin.Context) {

	var checkArticle models.Article
	generateSlug := slug.Make(c.PostForm("title"))
	err := config.DB.Find(&checkArticle, "slug = ?", generateSlug)
	if err.RowsAffected > 0 {
		generateSlug = generateSlug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}
	item := models.Article{
		Title:  c.PostForm("title"),
		Desc:   c.PostForm("desc"),
		Slug:   generateSlug,
		Tag:    c.PostForm("tag"),
		UserId: uint(c.MustGet("jwt_user_id").(float64)),
	}
	config.DB.Create(&item)
	c.JSON(200, gin.H{
		"message": "Berhasil",
		"data":    item,
	})
}

func GetArticleByTag(c *gin.Context) {
	tag := c.Param("tag")
	var items []models.Article
	// err := config.DB.First(&item, "Slug = ?", slug).Error
	err := config.DB.Where("tag LIKE ? ", "%"+tag+"%").Find(&items).Error
	if err != nil {
		c.JSON(404, gin.H{"status": "error", "message": "record not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "Berhasil",
		"data":   items,
	})
}

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	var items models.Article
	config.DB.Where("id = ? ", id).Delete(&items)

	c.JSON(200, gin.H{
		"status": "Berhasil",
		"data":   items,
	})
}
