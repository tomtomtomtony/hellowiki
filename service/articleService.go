package service

import (
	"hellowiki/common"
	"hellowiki/common/result"
	"hellowiki/model"
)

func CreateArticle(article model.Article) int {
	if model.HasCategoryById(article.Category.ID) != result.SUCCSE {
		return result.ERROR
	}
	inputTBName := model.UNCLASSIFIED_ARTICLES
	if article.Category.Name != "" {
		inputTBName = article.Category.Name + common.UNDER_SCORE + string(article.Category.ID)
	}
	var err error
	if !model.Db.Migrator().HasTable(inputTBName) {
		err = model.Db.Migrator().CreateTable(article)
		if err != nil {
			return result.ERROR
		}
	}
	return model.CreateArticle(article, inputTBName)

}
