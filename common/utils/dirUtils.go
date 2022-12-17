package utils

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"hellowiki/api/result"
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
	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {
			log.Printf("未能正确关闭文件夹:{%v}", err)
		}
	}(dir)
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

func GetEnforcer() (enforcer *casbin.Enforcer) {
	policyFile := config.Cfg.AuthenticationDB.AbsPath + string(os.PathSeparator) + config.Cfg.AuthenticationDB.PolicyFile
	modelFile := config.Cfg.AuthenticationDB.AbsPath + string(os.PathSeparator) + config.Cfg.AuthenticationDB.ModelFile
	csvAdapter := fileadapter.NewAdapter(policyFile)
	enforcer, err := casbin.NewEnforcer(modelFile, csvAdapter)
	if err != nil {
		log.Printf("创建鉴权器失败:{%v}", err)
		return nil
	}
	return enforcer
}
