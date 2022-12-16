package service

import (
	"hellowiki/api/result"
	vo2 "hellowiki/api/v1/article/vo"
	"hellowiki/api/v1/category/vo"
	"hellowiki/common/utils"
	"hellowiki/config"
	"hellowiki/model"
	"hellowiki/model/searchConfig"
	utils2 "hellowiki/model/utils"
	"io/fs"
	"log"
	"os"
	"strings"
)

func CreateCategory(condition vo.CategoryVO) (code int) {
	if 0 == len(condition.ParentPath) {
		condition.ParentPath = config.Cfg.DirDB.AbsPath
	}
	absPathInContent := condition.ParentPath + string(os.PathSeparator) + condition.Name
	//文件夹下是否已存在同名文件夹
	check, _ := utils.HasDirectoryOrFile(absPathInContent)
	if check {
		return result.ERROR_CATEGORY_EXIST
	}

	//写入磁盘文件夹
	code = model.WriteToContentDir(absPathInContent)
	if code != result.SUCCSE {
		log.Println("新增分类写入磁盘失败")
		return result.ERROR
	}

	//写入索引文件
	indexName := condition.Name
	tokenOpt := map[string]interface{}{
		"dicts":     config.Cfg.Analyze.Dict,
		"stop":      "",
		"opt":       "search-hmm",
		"trim":      "trim",
		"alpha":     false,
		"type":      searchConfig.TokenName,
		"tokenizer": searchConfig.TokenName,
	}
	articlesMapping := utils2.BuildArticleMapping(tokenOpt)
	indexParentPath := strings.Replace(absPathInContent, config.Cfg.DirDB.AbsPath, config.Cfg.SearchDB.AbsPath, -1)
	code = model.WriteToCategoryIndex(indexParentPath, indexName, articlesMapping)
	if code != result.SUCCSE {
		return code
	}

	return result.SUCCSE
}

func DeleteCategory(condition vo.CategoryVO) int {

	//2.删除txt文件
	if !utils.FoldIsEmpty(condition.Path) {
		log.Printf("content文件夹下指定文件夹{%v}不为空", condition.Path)
		return result.ERROR_CATEGORY_NOT_EMPTY
	}
	code := model.DeleteCategoryInContent(condition.Path)
	if code != result.SUCCSE {
		log.Println("content文件夹下指定文件夹删除失败")
		return result.ERROR
	}

	//3.删除索引文件
	indexName := strings.Replace(condition.Path, config.Cfg.DirDB.AbsPath, config.Cfg.SearchDB.AbsPath, -1)
	if !HasCategoryInIndex(indexName) {
		log.Println("index文件夹下指定文件夹不存在")
		return result.ERROR
	}
	code = model.DeleteCategoryInIndex(indexName)
	if code != result.SUCCSE {
		log.Println("index文件夹下指定文件夹删除失败")
		return result.ERROR
	}
	return result.SUCCSE

}

func HasCategoryInIndex(indexName string) bool {
	check, _ := model.HasCategoryInIndexDir(indexName)
	return check
}

func GetTopCategory() []vo.CategoryVO {
	resRaw := model.GetTopLevelCategory()
	var res = make([]vo.CategoryVO, 0, len(resRaw))
	for _, item := range resRaw {
		res = append(res, do2TopVo(item))
	}
	return res
}
func do2TopVo(do fs.DirEntry) vo.CategoryVO {
	var res vo.CategoryVO
	res.Name = do.Name()
	res.ParentName = ""
	if do.IsDir() {
		res.Type = "category"
	} else {
		res.Type = "article"
	}
	res.ParentPath = config.Cfg.DirDB.AbsPath
	res.Path = config.Cfg.DirDB.AbsPath + string(os.PathSeparator) + do.Name()
	return res
}
func GetDirectChildren(currentPath string) []vo.CategoryVO {
	resRaw := model.GetNextLevelCategory(currentPath)
	var res = make([]vo.CategoryVO, 0, len(resRaw))
	for _, item := range resRaw {
		res = append(res, do2Vo(item, currentPath))
	}
	return res
}
func do2Vo(do fs.DirEntry, parentPath string) vo.CategoryVO {
	var res vo.CategoryVO

	res.Name = do.Name()
	res.ParentName = ""
	if do.IsDir() {
		res.Type = "category"
	} else {
		res.Type = "article"
	}
	res.ParentPath = parentPath
	res.Path = parentPath + string(os.PathSeparator) + do.Name()
	return res
}
func GetArticleList(condition vo.CategoryVO) ([]vo2.ResultArticle, int64) {
	var res []vo2.ResultArticle
	resRaw, total := model.GetAllArticle(condition.PageSize, condition.PageNum)
	for i := 0; i < len(resRaw); i++ {
		res = append(res, do2ResultArticle(resRaw[i]))
	}

	return res, total
}

func do2ResultArticle(do model.Menu) vo2.ResultArticle {
	var res vo2.ResultArticle
	res.Title = do.Name
	res.CategoryName = do.ParentName
	res.CreateAt = do.CreatedAt.UnixMilli()
	res.UpdateAt = do.UpdatedAt.UnixMilli()
	return res
}
