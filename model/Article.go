package model

import (
	"bufio"
	"container/list"
	"encoding/json"
	"github.com/blevesearch/bleve/v2"
	"hellowiki/api/result"
	"hellowiki/api/v1/article/vo"
	"hellowiki/common/utils"
	"hellowiki/config"
	"log"
	"os"
	"strings"
)

type Article struct {
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Img      string    `json:"img"`
	KeyWords list.List `json:"keyWords"`
	Author   string    `json:"author"`
	CreateAt int64     `json:"createAt"`
}

var (
	UNCLASSIFIED_ARTICLES = "unclassified_articles"
)

func GetArticleByName(articlePath string) (res string, code int) {

	articleContent, err := os.ReadFile(articlePath)
	if err != nil {
		log.Printf("未能读取指定文件:{%v}", err)
		return res, result.ERROR
	}
	res = string(articleContent)
	return res, result.SUCCSE
}

func DeleteArticleByAbsPath(absPath string) (code int) {

	code = utils.DeleteFold(absPath)
	if code != result.SUCCSE {
		log.Printf("删除失败")
		return code
	}
	return result.SUCCSE
}

func DeleteArticleInIndex(vo vo.ConditionVO) (code int) {
	classifiedPathInIndexDir := strings.Replace(vo.ParentPath, config.Cfg.DirDB.AbsPath, config.Cfg.SearchDB.AbsPath, -1)
	//在指定索引删除文章记录
	dbSearch, code := utils.OpenIndex(classifiedPathInIndexDir)
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
	if err == bleve.ErrorEmptyID {
		log.Printf("该文章未创建索引:{%v}", err)
	} else if err != nil && err != bleve.ErrorEmptyID {
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

func ArticleWriteDir(article Article, path string) int {

	fileHandle, err := os.Create(path)
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
func ArticleUpdateInContent(articleAbsPath string, article Article) (code int) {
	fHandle, err := os.OpenFile(articleAbsPath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("未能打开指定的文章{%v}:{%v}\n", articleAbsPath, err)
		return result.ERROR
	}
	//清空文件
	err = fHandle.Truncate(0)
	if err != nil {
		log.Printf("清空文件失败:{%v}", err)
		return result.ERROR
	}
	fHandle.Seek(0, 0)
	write := bufio.NewWriter(fHandle)
	articleJson, err := json.Marshal(article)
	if err != nil {
		log.Printf("将输入内容转换为json格式失败:{%v}", err)
		return result.ERROR
	}
	if _, err = write.WriteString(string(articleJson)); err != nil {
		log.Fatalf("写入磁盘错误，未能存储指定文章:{%v}", err)
	}
	if err := write.Flush(); err != nil {
		// failed to encode
		log.Fatalf("刷入磁盘错误，未能存储指定文章:{%v}", err)
		return result.ERROR
	}
	if err := fHandle.Close(); err != nil {
		log.Fatalf("未能正确关闭文件:{%v}", err)
		return result.ERROR
	}
	return result.SUCCSE
}
