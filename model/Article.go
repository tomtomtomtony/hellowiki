package model

import (
	"bufio"
	"github.com/blevesearch/bleve/v2"
	"hellowiki/common"
	"hellowiki/common/result"
	"hellowiki/common/utils"
	"hellowiki/config"
	"hellowiki/service/search"
	"log"
	"os"
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

func HasCategoryContent(contentName string) (bool, error) {
	return utils.HasDirectory(contentName)
}

func CreateArticle(article Article, classifiedName string) int {

	allIndexName := config.Cfg.SearchDB.Location + classifiedName

	//写入索引
	dbSearch, code := search.OpenIndex(allIndexName)
	if code != result.SUCCSE {
		return code
	}
	defer dbSearch.Close()
	docId := article.Title + common.UNDER_SCORE + strconv.FormatInt(time.Now().Unix(), 10)
	err = dbSearch.Index(docId, article)
	if err != nil {
		return result.ERROR
	}

	//写入磁盘
	dirName := config.Cfg.DirDB.Location + classifiedName
	//检查分类文件夹是否存在
	checkContentDir, err := HasCategoryContent(dirName)
	if !checkContentDir || err != nil {
		log.Printf("写入磁盘错误，未能打开{%v}:{%v}\n", dirName, err)
		return result.ERROR
	}
	contentPath := dirName + string(os.PathSeparator) + article.Title + common.UNDER_SCORE +
		strconv.FormatInt(time.Now().Unix(), 10) + common.MD_FILE_SUFFIX
	fileHandle, err := os.Create(contentPath)
	if err != nil {
		log.Fatal(err)
	}
	write := bufio.NewWriter(fileHandle)
	if _, err = write.WriteString(article.Content); err != nil {
		log.Fatalf("写入磁盘错误，未能存储指定文章:{%v}", err)
	}
	if err := write.Flush(); err != nil {
		// failed to encode
		log.Fatalf("刷入磁盘错误，未能存储指定文章:{%v}", err)

	}
	if err := fileHandle.Close(); err != nil {
		// failed to close the file
		log.Fatalf("未能正确关闭文件:{%v}", err)
	}

	return result.SUCCSE

}

func GetAllInAIndex(pageSize int, pageNum int, indexName string) (bleve.SearchResult, int) {

	allIndexName := config.Cfg.SearchDB.Location + indexName
	dbSearch, code := search.OpenIndex(allIndexName)
	if code != result.SUCCSE {
		return bleve.SearchResult{}, code
	}
	defer dbSearch.Close()
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"img", "des", "tag", "title", "content"}
	searchResult, _ := dbSearch.Search(searchRequest)
	return *searchResult, result.SUCCSE
}
