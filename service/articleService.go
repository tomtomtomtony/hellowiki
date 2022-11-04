package service

import (
	"hellowiki/api/v1/article/vo"
	"hellowiki/common/result"
	"hellowiki/model"
)

func CreateArticle(condition vo.ConditionVO) int {
	if !HasCategory(condition.CategoryId) {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	TBName := model.UNCLASSIFIED_ARTICLES
	if condition.CategoryName != "" {
		TBName = ConstructStandardTableName(condition.CategoryEngName, condition.CategoryId)
	}
	if !model.Db.Migrator().HasTable(TBName) {
		return result.ERROR_ARTICLE_DATABASE_TABLE_NOT_FOUND
	}
	var article model.Article
	article = voTDo(condition)
	return model.CreateArticle(article, TBName)
}

func QueryInCategory(condition vo.ConditionVO) []model.Article {
	tableName := ConstructStandardTableName(condition.CategoryEngName, condition.CategoryId)
	if !model.Db.Migrator().HasTable(tableName) {
		return []model.Article{}
	}
	return model.GetAllInTable(condition.PageSize, condition.PageNum, tableName)
}

func voTDo(conditionVO vo.ConditionVO) model.Article {
	var article model.Article
	article.Category.Name = conditionVO.CategoryName
	article.Title = conditionVO.ArticleTitle
	article.Content = conditionVO.ArticleContent
	return article
}
