package service

import (
	"hellowiki/common"
	"hellowiki/common/result"
	"hellowiki/model"
	"strconv"
)

func CreateArticle(article model.Article) int {
	if model.HasCategoryById(article.Category.ID) == result.SUCCSE {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	TBName := model.UNCLASSIFIED_ARTICLES
	if article.Category.Name != "" {
		TBName = article.Category.EngName + common.UNDER_SCORE + strconv.Itoa(int(article.Category.ID))
	}
	if !model.Db.Migrator().HasTable(TBName) {
		return result.ERROR
	}
	return model.CreateArticle(article, TBName)

}
