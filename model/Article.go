package model

import (
	"gorm.io/gorm"
	"hellowiki/common/result"
)

type Article struct {
	gorm.Model
	Category  Category `gorm:"-" json:"category"`
	Title     string   `gorm:"type:varchar(100);not null" json:"title"`
	Content   string   `gorm:"type:longtext" json:"content"`
	Img       string   `gorm:"type:varchar(255)" json:"img"`
	Desc      string   `gorm:"type:varchar(255)" json:"desc"`
	Tag       string   `gorm:"type:varchar(255)" json:"tag"`
	tableName string   `gorm:"-"`
}

var (
	UNCLASSIFIED_ARTICLES = "unclassified_articles"
)

func (article *Article) TableName() string {
	return UNCLASSIFIED_ARTICLES
}

func CreateArticle(article Article, inputTBName string) int {
	var err error
	if len(inputTBName) != 0 {
		err = Db.Table(inputTBName).Create(&article).Error
	} else {
		err = Db.Table(UNCLASSIFIED_ARTICLES).Create(article).Error
	}
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE

}

func GetAllInTable(pageSize int, pageNum int, tableName string) []Article {
	var article []Article
	err := Db.Table(tableName).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []Article{}
	}
	return article
}
