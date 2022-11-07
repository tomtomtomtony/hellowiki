package service

import (
	"github.com/blevesearch/bleve/v2"
	"hellowiki/api/result"
	"hellowiki/api/v1/article/vo"
	"hellowiki/model"
)

func CreateArticle(condition vo.ConditionVO) int {
	if !HasCategory(condition.CategoryId) {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	IndexName := model.UNCLASSIFIED_ARTICLES
	if condition.CategoryName != "" {
		IndexName = ConstructStandardIndexName(condition.CategoryEngName, condition.CategoryId)
	}
	var article model.Article
	article = voTDo(condition)
	return model.CreateArticle(article, IndexName)
}

func QueryInCategory(condition vo.ConditionVO) (bleve.SearchResult, int) {
	tableName := ConstructStandardIndexName(condition.CategoryEngName, condition.CategoryId)
	if !HasArticleIndex(tableName) {
		return bleve.SearchResult{}, result.ERROR
	}
	res, code := model.GetAllInAIndex(condition.PageSize, condition.PageNum, tableName)
	if code != result.SUCCSE {
		return bleve.SearchResult{}, code
	}
	return res, code
}

func voTDo(conditionVO vo.ConditionVO) model.Article {
	var article model.Article
	//article.Category.Name = conditionVO.CategoryName
	article.Title = conditionVO.ArticleTitle
	article.Content = conditionVO.ArticleContent
	return article
}
