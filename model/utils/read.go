package utils

import (
	"github.com/blevesearch/bleve/v2"
	"hellowiki/api/result"
	"hellowiki/config"
	"log"
	"os"
)

func OpenIndex(indexName string) (bleve.Index, int) {
	if !HasCategoryIndex(indexName) {
		return nil, result.ERROR
	}
	dbSearch, err := bleve.Open(config.Cfg.SearchDB.AbsPath + string(os.PathSeparator) + indexName)
	if err != nil {
		log.Printf("打开{%v}索引失败:{%v}", indexName, err)
		return nil, result.ERROR
	}
	return dbSearch, result.SUCCSE
}

func HasCategoryIndex(indexName string) bool {
	_, err := os.OpenFile(config.Cfg.SearchDB.AbsPath+indexName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return false
	}
	return true
}
