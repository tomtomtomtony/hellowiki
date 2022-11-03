package model

import (
	"gorm.io/gorm"
	"hellowiki/common"
	"hellowiki/common/result"
)

type Article struct {
	gorm.Model
	Category   Category `gorm:"-" json:"category"`
	Title      string   `gorm:"type:varchar(100);not null" json:"title"`
	Content    string   `gorm:"type:longtext" json:"content"`
	Img        string   `gorm:"type:varchar(255)" json:"img"`
	Desc       string   `gorm:"type:varchar(255)" json:"desc"`
	Tag        string   `gorm:"type:varchar(255)" json:"tag"`
	CategoryId uint     `gorm:"type:int;not null" json:"categoryId"`
}

var (
	UNCLASSIFIED_ARTICLES = "unclassified_articles"
)

func (article *Article) TableName() string {
	if article.Category.Name != "" {
		return article.Category.Name + common.UNDER_SCORE + string(article.Category.ID)
	}
	return UNCLASSIFIED_ARTICLES
}

func CreateArticle(article Article, inputTBName string) int {
	var err error
	if len(inputTBName) != 0 {
		err = Db.Table(inputTBName).Create(article).Error
	} else {
		err = Db.Table(UNCLASSIFIED_ARTICLES).Create(article).Error
	}
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE

}
