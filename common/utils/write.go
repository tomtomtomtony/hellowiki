package utils

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"hellowiki/api/result"
	"hellowiki/config"
	"log"
	"os"
)

// 写入data/index
func WriteToIndexDir(indexName string, mapping mapping.IndexMapping) int {
	index, err := bleve.New(config.Cfg.SearchDB.AbsPath+string(os.PathSeparator)+indexName, mapping)
	if err != nil {
		log.Println("新增分类写入索引失败")
		return result.ERROR
	}
	if err := index.Close(); err != nil {
		log.Fatalf("未能正确关闭文件")
	}
	return result.SUCCSE
}

// 写入data/content
func WriteToFS() {

}
