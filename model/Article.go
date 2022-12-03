package model

import (
	"bufio"
	"container/list"
	"encoding/json"
	"github.com/blevesearch/bleve/v2"
	"hellowiki/api/result"
	"hellowiki/api/v1/article/vo"
	"hellowiki/common"
	"hellowiki/common/utils"
	"hellowiki/config"
	utils2 "hellowiki/model/utils"
	"log"
	"os"
)

type Article struct {
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Img      string    `json:"img"`
	Desc     string    `json:"desc"`
	KeyWords list.List `json:"keyWords"`
	Author   string    `json:"author"`
	CreateAt int64     `json:"createAt"`
}

var (
	UNCLASSIFIED_ARTICLES = "unclassified_articles"
)

func HasArticleInContent(categoryNameId string, articleName string) bool {
	check, _ := utils.HasMdFileInContentDir(categoryNameId, articleName)
	return check
}

func GetArticleByName(vo vo.ConditionVO) (res string, code int) {
	dirName := utils2.ConstructCategoryNameId(vo.CategoryName, vo.CategoryMenuId)
	dirAbs := config.Cfg.DirDB.AbsPath + string(os.PathSeparator) + dirName
	articleAbsPath := dirAbs + string(os.PathSeparator) + vo.ArticleTitle + common.JSON_FILE_SUFFIX
	articleContent, err := os.ReadFile(articleAbsPath)
	if err != nil {
		log.Printf("未能读取指定文件:{%v}", err)
		return res, result.ERROR
	}
	res = string(articleContent)
	return res, result.SUCCSE
}

func DeleteArticleByAbsPath(vo vo.ConditionVO) (code int) {
	dirName := utils2.ConstructCategoryNameId(vo.CategoryName, vo.CategoryMenuId)
	dirAbs := config.Cfg.DirDB.AbsPath + string(os.PathSeparator) + dirName
	articleAbsPath := dirAbs + string(os.PathSeparator) + vo.ArticleTitle + common.JSON_FILE_SUFFIX
	code = utils.DeleteFold(articleAbsPath)
	if code != result.SUCCSE {
		log.Printf("删除失败")
		return code
	}
	return result.SUCCSE
}

func DeleteArticleInIndex(vo vo.ConditionVO) (code int) {
	classifiedName := utils2.ConstructCategoryNameId(vo.CategoryName, vo.CategoryMenuId)
	//在指定索引删除文章记录
	dbSearch, code := utils.OpenIndex(classifiedName)
	if code != result.SUCCSE {
		log.Println("写入索引失败")
		return code
	}
	defer func(dbSearch bleve.Index) {
		err := dbSearch.Close()
		if err != nil {
			log.Printf("未能正确关闭索引:{%v}", err)
		}
	}(dbSearch)
	docId := vo.ArticleTitle
	err := dbSearch.Delete(docId)
	if err != nil {
		log.Printf("未能删除index中文章记录:{%v}", err)
		return result.ERROR
	}
	return result.SUCCSE
}

func ArticleWriteIndex(article Article, classifiedName string) int {
	//写入索引
	dbSearch, code := utils.OpenIndex(classifiedName)
	if code != result.SUCCSE {
		log.Println("写入索引失败")
		return code
	}
	defer func(dbSearch bleve.Index) {
		err := dbSearch.Close()
		if err != nil {
			log.Printf("未能正确关闭索引:{%v}", err)
		}
	}(dbSearch)
	docId := article.Title
	err := dbSearch.Index(docId, article)
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

func ArticleWriteDir(article Article, classifiedName string) int {
	//写入磁盘
	dirName := config.Cfg.DirDB.AbsPath + string(os.PathSeparator) + classifiedName
	//检查分类文件夹是否存在
	checkContentDir, err := utils.HasCategoryInContentDir(classifiedName)
	if !checkContentDir || err != nil {
		log.Printf("写入磁盘错误，未能找到{%v}:{%v}\n", dirName, err)
		return result.ERROR
	}
	contentPath := dirName + string(os.PathSeparator) + article.Title + common.JSON_FILE_SUFFIX
	fileHandle, err := os.Create(contentPath)
	if err != nil {
		log.Fatal(err)
	}
	write := bufio.NewWriter(fileHandle)
	articleJson, err := json.Marshal(article)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = write.WriteString(string(articleJson)); err != nil {
		log.Fatalf("写入磁盘错误，未能存储指定文章:{%v}", err)
	}
	if err := write.Flush(); err != nil {
		// failed to encode
		log.Fatalf("刷入磁盘错误，未能存储指定文章:{%v}", err)
		return result.ERROR
	}
	if err := fileHandle.Close(); err != nil {
		log.Fatalf("未能正确关闭文件:{%v}", err)
		return result.ERROR
	}
	return result.SUCCSE
}

func ArticleWriteMenu(menu Menu) int {
	//写入数据库菜单表
	dbBase := utils.OpenDB()
	err := dbBase.Create(&menu).Error
	if err != nil {
		log.Fatalf("写入数据库失败:{%v}", err)
		return result.ERROR
	}
	return result.SUCCSE
}

func GetAllInAIndex(pageSize int, pageNum int, indexName string) (bleve.SearchResult, int) {
	allIndexName := config.Cfg.SearchDB.Location + indexName
	dbSearch, code := utils.OpenIndex(allIndexName)
	if code != result.SUCCSE {
		return bleve.SearchResult{}, code
	}
	defer func(dbSearch bleve.Index) {
		err := dbSearch.Close()
		if err != nil {
			log.Printf("未能正确关闭索引:{%v}", err)
		}
	}(dbSearch)
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = pageSize
	searchRequest.From = pageNum
	searchRequest.Fields = []string{"img", "des", "tag", "title", "content"}
	searchResult, _ := dbSearch.Search(searchRequest)
	return *searchResult, result.SUCCSE
}

func GetAllArticleTitleInCategory(indexName string) (bleve.SearchResult, int) {
	dbSearch, code := utils.OpenIndex(indexName)
	if code != result.SUCCSE {
		return bleve.SearchResult{}, code
	}
	defer func(dbSearch bleve.Index) {
		err := dbSearch.Close()
		if err != nil {
			log.Printf("未能正确关闭索引:{%v}", err)
		}
	}(dbSearch)
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"id"}
	searchResult, _ := dbSearch.Search(searchRequest)
	return *searchResult, result.SUCCSE
}
