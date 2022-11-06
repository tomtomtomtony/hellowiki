package search

import (
	"github.com/blevesearch/bleve/v2"
	"hellowiki/common/result"
	"log"
)

func OpenIndex(indexName string) (bleve.Index, int) {
	dbSearch, err := bleve.Open(indexName)
	if err != nil {
		log.Printf("打开{%v}索引失败:{%v}", indexName, err)
		return nil, result.ERROR
	}
	return dbSearch, result.SUCCSE
}
