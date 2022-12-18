package service

import (
	"fmt"
	"hellowiki/api/result"
	"hellowiki/api/v1/article/vo"
	"hellowiki/common"
	utils2 "hellowiki/common/utils"
	"hellowiki/config"
	"hellowiki/model"
	"log"
	"os"
	"strings"
	"time"
)

func DeleteArticle(condition vo.ConditionVO) (code int) {
	//1.删除content下的文件
	//检查分类文件夹是否存在
	checkContentDir, err := utils2.HasDirectoryOrFile(condition.ParentPath)
	if !checkContentDir || err != nil {
		log.Printf("未能找到{%v}:{%v}\n", condition.ParentPath, err)
		return result.ERROR
	}
	//检查文章是否已存在
	checkContentDir, err = utils2.HasDirectoryOrFile(condition.Path)
	if !checkContentDir || os.IsNotExist(err) {
		log.Printf("文章不存在，未能找到{%v}:{%v}\n", condition.Path, err)
		return result.ERROR
	}

	code = model.DeleteArticleByAbsPath(condition.Path)
	if code != result.SUCCSE {
		log.Println("content文件夹下指定文件夹删除失败")
		return result.ERROR
	}

	//2.删除索引文件
	articleCategoryIndexName := strings.Replace(condition.ParentPath, config.Cfg.DirDB.AbsPath, config.Cfg.SearchDB.AbsPath, -1)
	//检查分类文件夹是否存在
	checkContentDir, err = utils2.HasDirectoryOrFile(articleCategoryIndexName)
	if !checkContentDir || err != nil {
		log.Printf("未能找到{%v}:{%v}\n", condition.ParentPath, err)
		return result.ERROR
	}

	code = model.DeleteArticleInIndex(condition)
	if code != result.SUCCSE {
		log.Println("index文件夹下指定文章删除失败")
		return result.ERROR
	}

	return result.SUCCSE
}

func UpdateArticle(condition vo.ConditionVO) (code int) {
	//检查分类文件夹是否存在
	checkContentDir, err := utils2.HasDirectoryOrFile(condition.ParentPath)
	if !checkContentDir || err != nil {
		log.Printf("写入磁盘错误，未能找到{%v}:{%v}\n", condition.ParentPath, err)
		return result.ERROR
	}
	checkContentDir, err = utils2.HasDirectoryOrFile(condition.Path)
	if !checkContentDir || os.IsNotExist(err) {
		log.Printf("文章不存在，未能找到{%v}:{%v}\n", condition.Path, err)
		return result.ERROR
	}

	article := voTDo(condition)
	code = model.ArticleUpdateInContent(condition.Path, article)
	if code != result.ERROR {
		log.Println("写入磁盘错误")
		return code
	}
	return result.SUCCSE
}

func GetArticle(condition vo.ConditionVO) (res string, code int) {
	//检查分类文件夹是否存在
	checkContentDir, err := utils2.HasDirectoryOrFile(condition.ParentPath)
	if !checkContentDir || err != nil {
		log.Printf("未能找到{%v}:{%v}\n", condition.ParentPath, err)
		return res, result.ERROR
	}
	//检查文章是否已存在
	checkContentDir, err = utils2.HasDirectoryOrFile(condition.Path)
	if !checkContentDir || os.IsNotExist(err) {
		log.Printf("文章不存在，未能找到{%v}:{%v}\n", condition.Path, err)
		return res, result.ERROR
	}
	res, code = model.GetArticleByName(condition.Path)
	if code != result.SUCCSE {
		return res, code
	}
	code = result.SUCCSE
	return res, code
}

func CreateArticle(condition vo.ConditionVO) int {
	articlePath := condition.ParentPath + string(os.PathSeparator) + condition.ArticleTitle + common.JSON_FILE_SUFFIX
	//检查分类文件夹是否存在
	checkContentDir, err := utils2.HasDirectoryOrFile(condition.ParentPath)
	if !checkContentDir || err != nil {
		log.Printf("写入磁盘错误，未能找到{%v}:{%v}\n", condition.ParentPath, err)
		return result.ERROR
	}
	//检查文章是否已存在
	checkContentDir, err = utils2.HasDirectoryOrFile(articlePath)
	if checkContentDir || !os.IsNotExist(err) {
		log.Printf("文章已存在，未能找到{%v}:{%v}\n", articlePath, err)
		return result.ERROR
	}
	var article model.Article
	article = voTDo(condition)
	article.CreateAt = time.Now().Unix()
	//code := model.ArticleWriteIndex(article, IndexName)
	//if code != result.SUCCSE {
	//	return code
	//}

	code := model.ArticleWriteDir(article, articlePath)
	if code != result.SUCCSE {
		return code
	}

	return code
}

func QueryInCategory(condition vo.ConditionVO) ([]string, int) {
	categoryNameInContent := condition.CategoryName

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
	article.Title = conditionVO.ArticleTitle
	article.Content = conditionVO.ArticleContent
	article.Author = conditionVO.Author
	article.KeyWords = conditionVO.Keywords
	return article
}
