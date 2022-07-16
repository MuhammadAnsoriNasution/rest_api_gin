package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title  string
	Slug   string
	Desc   string
	Tag    string
	UserId uint
}
