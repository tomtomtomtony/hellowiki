package service

import (
	"hellowiki/api/v1/article/vo"
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

func QueryInCategory(condition vo.ConditionVO) []model.Article {
	tableName := condition.EngName + common.UNDER_SCORE + strconv.Itoa(int(condition.ID))
	if !model.Db.Migrator().HasTable(tableName) {
		return []model.Article{}
	}
	return model.GetAllInTable(condition.PageSize, condition.PageNum, tableName)
}
