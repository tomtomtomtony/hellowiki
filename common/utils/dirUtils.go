package utils

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"hellowiki/api/result"
	"hellowiki/config"
	"io"
	"log"
	"os"
)

func HasDirectory(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FoldIsEmptyInContent(foldName string) bool {
	absPath := config.Cfg.DirDB.AbsPath + string(os.PathSeparator) + foldName
	dir, err := os.Open(absPath)
	if err != nil {
		log.Printf("不能打开文件夹{%v}\n", absPath)
		return false
	}
	defer dir.Close()
	_, err = dir.Readdirnames(1)
	return err == io.EOF
}

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

// 删除data/index指定目录
func DeleteFold(absPath string) int {
	err := os.RemoveAll(absPath)
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 在data/content创建文件夹
func CreateFoldContent(foldName string) int {
	err := os.Mkdir(config.Cfg.DirDB.AbsPath+string(os.PathSeparator)+foldName, os.ModePerm)
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

func OpenIndex(indexName string) (bleve.Index, int) {
	_, err := HasCategoryIndex(indexName)
	if err != nil {
		return nil, result.ERROR
	}

	dbSearch, err := bleve.Open(config.Cfg.SearchDB.AbsPath + string(os.PathSeparator) + indexName)
	if err != nil {
		log.Printf("打开{%v}索引失败:{%v}", indexName, err)
		return nil, result.ERROR
	}
	return dbSearch, result.SUCCSE
}

func HasCategoryIndex(indexName string) (bool, error) {
	return HasDirectory(config.Cfg.SearchDB.AbsPath + string(os.PathSeparator) + indexName)
}
