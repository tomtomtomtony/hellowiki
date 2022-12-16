package utils

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"hellowiki/api/result"
	"hellowiki/common"
	"hellowiki/config"
	"io"
	"log"
	"os"
)

func HasDirectoryOrFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, err
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

func FoldIsEmpty(absPath string) bool {
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
func WriteToIndexDir(parentPath string, indexName string, mapping mapping.IndexMapping) int {
	index, err := bleve.New(parentPath+string(os.PathSeparator)+indexName, mapping)
	if err != nil {
		log.Println("新增分类写入索引失败")
		return result.ERROR
	}
	if err := index.Close(); err != nil {
		log.Fatalf("未能正确关闭文件")
	}
	return result.SUCCSE
}

// 删除指定目录
func DeleteFold(absPath string) int {
	err := os.RemoveAll(absPath)
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 在data/content创建文件夹
func CreateFoldContent(parentPath string, foldName string) int {
	err := os.Mkdir(parentPath+string(os.PathSeparator)+foldName, os.ModePerm)
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

func OpenIndex(absPath string) (bleve.Index, int) {
	_, err := HasDirectoryOrFile(absPath)
	if err != nil {
		return nil, result.ERROR
	}

	dbSearch, err := bleve.Open(absPath)
	if err != nil {
		log.Printf("打开{%v}索引失败:{%v}", absPath, err)
		return nil, result.ERROR
	}
	return dbSearch, result.SUCCSE
}

func HasCategoryInIndexDir(indexName string) (bool, error) {
	return HasDirectoryOrFile(config.Cfg.SearchDB.AbsPath + string(os.PathSeparator) + indexName)
}

func HasCategoryInContentDir(categoryName string) (bool, error) {
	return HasDirectoryOrFile(config.Cfg.DirDB.AbsPath + string(os.PathSeparator) + categoryName)
}

func HasMdFileInContentDir(categoryNameId string, mdFileName string) (bool, error) {
	return HasDirectoryOrFile(config.Cfg.DirDB.AbsPath + string(os.PathSeparator) + categoryNameId + string(os.PathSeparator) + mdFileName + common.JSON_FILE_SUFFIX)

}
