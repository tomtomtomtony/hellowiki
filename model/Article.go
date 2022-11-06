package model

import (
	"github.com/blevesearch/bleve/v2"
	"hellowiki/common"
	"hellowiki/common/result"
	"hellowiki/config"
	"log"
	"strconv"
	"time"
)

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Img     string `json:"img"`
	Desc    string `json:"desc"`
	Tag     string `json:"tag"`
}

var (
	UNCLASSIFIED_ARTICLES = "unclassified_articles"
)

func openIndex(indexName string) bleve.Index {
	dbSearch, err := bleve.Open(indexName)
	if err != nil {
		log.Fatalf("打开{%v}索引失败", indexName)
	}
	return dbSearch
}

func CreateArticle(article Article, indexName string) int {

	allIndexName := config.Cfg.SearchDB.Location + indexName
	//
	//写入索引
	dbSearch := openIndex(allIndexName)
	defer dbSearch.Close()
	docId := article.Title + common.UNDER_SCORE + strconv.FormatInt(time.Now().Unix(), 10)
	err = dbSearch.Index(docId, article)
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE

}

func GetAllInAIndex(pageSize int, pageNum int, indexName string) bleve.SearchResult {

	allIndexName := config.Cfg.SearchDB.Location + indexName
	dbSearch := openIndex(allIndexName)
	defer dbSearch.Close()
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"img", "des", "tag", "title", "content"}
	searchResult, _ := dbSearch.Search(searchRequest)
	return *searchResult
}
