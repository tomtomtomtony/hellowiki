package model

import "gorm.io/gorm"

type Article struct {
	CatGory Category
	gorm.Model
	Title   string `gorm:"type:varchar(100);not null" json:"title"`
	Content string `gorm:"type:longtext" json:"content"`
	Img     string `gorm:"type:varchar(255)" json:"img"`
	//分类Id
	Cid  int    `gorm:"type:int;not null" json:"cid"`
	Desc string `gorm:"type:varchar(255);" json:"desc"`
	Tag  string `gorm:"type:varchar(255);" json:"tag"`
}
