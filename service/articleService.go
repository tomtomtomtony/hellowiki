package service

import (
	"fmt"
	"hellowiki/api/result"
	"hellowiki/api/v1/article/vo"
	"hellowiki/common"
	"hellowiki/model"
	"hellowiki/model/utils"
	"log"
	"strings"
)

func CreateArticle(condition vo.ConditionVO) int {
	if !HasCategoryInDBTable(condition.CategoryMenuId) {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	IndexName := model.UNCLASSIFIED_ARTICLES
	if condition.CategoryName != "" {
		IndexName = utils.ConstructStandardIndexName(condition.CategoryName, condition.CategoryMenuId)
	}
	var article model.Article
	article = voTDo(condition)

	code := model.ArticleWriteIndex(article, IndexName)
	if code != result.SUCCSE {
		return code
	}
	code = model.ArticleWriteDir(article, IndexName)
	if code != result.SUCCSE {
		return code
	}

	code = model.ArticleWriteMenu(vo2Menu(condition))
	return code
}

func vo2Menu(vo vo.ConditionVO) model.Menu {
	var Do model.Menu
	Do.Name = vo.ArticleTitle
	Do.ParentId = vo.CategoryMenuId
	Do.ParentName = vo.CategoryName
	Do.Type = common.ARTICLE_TYPE
	return Do
}

func QueryInCategory(condition vo.ConditionVO) ([]string, int) {
	categoryNameInContent := utils.ConstructStandardIndexName(condition.CategoryName, condition.CategoryMenuId)

	//检查index中分类文件夹是否存在
	checkContentDir, err := model.HasCategoryInIndexDir(categoryNameInContent)
	if !checkContentDir || err != nil {
		log.Printf("读入错误，未能找到{%v}:{%v}\n", categoryNameInContent, err)
		return []string{}, result.ERROR
	}

	resRaw, code := model.GetAllArticleTitleInCategory(categoryNameInContent)
	res := make([]string, 0, resRaw.Hits.Len())
	for _, title := range resRaw.Hits {
		fmt.Printf("{%v}", title)
		res = append(res, strings.Split(title.ID, common.UNDER_SCORE)[0])
	}
	if code != result.SUCCSE {
		return res, code
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
