package model

import (
	"gorm.io/gorm"
	"hellowiki/common"
	"hellowiki/common/result"
)

type Article struct {
	Category Category
	gorm.Model
	Title     string `gorm:"type:varchar(100);not null" json:"title"`
	Content   string `gorm:"type:longtext" json:"content"`
	Img       string `gorm:"type:varchar(255)" json:"img"`
	Desc      string `gorm:"type:varchar(255);" json:"desc"`
	Tag       string `gorm:"type:varchar(255);" json:"tag"`
	tableName string `gorm:"-"`
}

var (
	UNCLASSIFIED_ARTICLES = "unclassified_articles"
)

func (article *Article) TableName() string {
	if article.tableName != "" {
		return article.tableName + common.UNDER_SCORE + string(article.Category.ID)
	}
	return UNCLASSIFIED_ARTICLES
}

func CreateArticle(article Article, inputTBName string) int {
	Db.AutoMigrate(article)
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
