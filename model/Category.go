package model

import (
	"github.com/blevesearch/bleve/v2/mapping"
	utils2 "hellowiki/common/utils"
	"hellowiki/config"
	"io/fs"
	"log"
	"os"
)

// 文章分类
type Category struct {
	Name       string `json:"name"`
	ParentName string `json:"parentName"`
}

var (
	//顶级父类的parentId
	TOPLEVELCATEGORY uint = 0
)

func HasCategoryInIndexDir(path string) (bool, error) {
	return utils2.HasDirectoryOrFile(path)
}

//

// 索引写入
func WriteToCategoryIndex(parentPath string, indexName string, mapping mapping.IndexMapping) int {
	return utils2.WriteToIndexDir(parentPath, indexName, mapping)
}

// data/content下分类文件夹写入
func WriteToContentDir(parentPath string, categoryName string) int {
	return utils2.CreateFoldContent(parentPath, categoryName)
}

// 索引文件夹删除
func DeleteCategoryInIndex(path string) int {
	return utils2.DeleteFold(path)
}

func DeleteCategoryInContent(path string) int {
	return utils2.DeleteFold(path)
}

func GetTopLevelCategory() []fs.DirEntry {
	filesInfo, err := os.ReadDir(config.Cfg.DirDB.AbsPath)
	if err != nil {
		log.Printf("不能读取文件夹{%v}\n", config.Cfg.DirDB.AbsPath)
		return []fs.DirEntry{}
	}
	return filesInfo
}

func GetNextLevelCategory(parentPath string) []fs.DirEntry {
	filesInfo, err := os.ReadDir(parentPath)
	if err != nil {
		log.Printf("不能读取文件夹{%v}\n", parentPath)
		return []fs.DirEntry{}
	}
	return filesInfo
}
