package service

import (
	"github.com/blevesearch/bleve/v2"
	"hellowiki/api/v1/category/vo"
	"hellowiki/common"
	"hellowiki/common/result"
	"hellowiki/config"
	"hellowiki/model"
	"hellowiki/service/search"
	utils2 "hellowiki/service/utils"
	"log"
	"os"
	"strconv"
)

func CreateCategory(condition vo.ConditionVO) (code int) {
	var data model.Category
	var err error
	if condition.ParentId == model.TOPLEVELCATEGORY {
		code = handleCreateRootCategory(condition)
	} else {
		code = handleCreateNonRootCategory(condition)
	}
	if code != result.SUCCSE {
		return result.ERROR
	}
	data = vo2Do(condition)
	//写入数据库
	code = model.CreateCategory(data)
	if code != result.SUCCSE {
		log.Println("新增分类写入数据库失败")
		return result.ERROR
	}
	//写入索引文件
	indexName := ConstructStandardIndexName(data.EngName, data.ID)
	tokenOpt := map[string]interface{}{
		"dicts":     config.Cfg.Analyze.Dict,
		"stop":      "",
		"opt":       "search-hmm",
		"trim":      "trim",
		"alpha":     false,
		"type":      search.TokenName,
		"tokenizer": search.TokenName,
	}
	articlesMapping := utils2.BuildArticleMapping(tokenOpt)
	index, err := bleve.New(config.Cfg.SearchDB.Location+indexName, articlesMapping)
	defer index.Close()
	if err != nil {
		log.Println("新增分类写入索引失败")
		return result.ERROR
	}

	//写入磁盘文件夹
	dirName := indexName
	err = os.Mkdir(config.Cfg.DirDB.Location+dirName, os.ModePerm)
	if err != nil {
		log.Println("新增分类写入磁盘失败")
		return result.ERROR
	}
	return result.SUCCSE
}

func handleCreateRootCategory(categoryInfo vo.ConditionVO) int {
	children := model.FindCategoryChildren(model.TOPLEVELCATEGORY)
	for _, curr := range children {
		if curr.Name == categoryInfo.Name {
			return result.ERROR_CATEGORY_EXIST
		}
	}
	return result.SUCCSE
}

func handleCreateNonRootCategory(categoryInfo vo.ConditionVO) int {
	if !HasCategory(categoryInfo.CategoryId) {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	return result.SUCCSE
}

func GetAllCategory(pageSize int, pageNum int) []model.Category {
	return model.FindAllCategory(pageSize, pageNum)
}

func DeleteCategory(category model.Category) int {
	children := model.FindCategoryChildren(category.ID)
	var newData model.Category
	tx := model.DbBase.Begin()
	//1.更新每个孩子节点的父节点Id和名称
	for _, curr := range children {
		//1.1若为待删节点为根节点，其直接子节点将成为顶级父节点
		if category.ParentId == model.TOPLEVELCATEGORY {
			newData.ParentId, newData.ParentName = model.TOPLEVELCATEGORY, curr.Name
		} else {
			newData.ParentId, newData.ParentName = category.ParentId, category.ParentName
		}
		if err := tx.Model(&model.Category{}).Where("id=?", curr.ID).Updates(newData).Error; err != nil {
			tx.Rollback()
			return result.ERROR
		}
	} //2.删除节点本身
	if err := tx.Delete(&model.Category{}, "id=?", category.ID).Error; err != nil {
		tx.Rollback()
		return result.ERROR
	}
	//3.若分类下没有文章，则对应的表也删除

	tx.Commit()
	return result.SUCCSE
}

func SetCategory(id uint, data model.Category) int {
	if !HasCategory(id) {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	return model.UpdateCategoryById(id, data)
}

func HasCategory(id uint) bool {
	return model.GetCategoryById(id) != model.Category{}
}

// 传入索引名称符合格式: article的 engName_categoryId
func HasArticleIndex(indexName string) bool {
	return model.HasCategoryDir(indexName)
}

//func repairTable(categoryId uint) int {
//	var category model.Category
//	category = model.GetCategoryById(categoryId)
//
//}

func ConstructStandardIndexName(categoryEngName string, categoryId uint) string {
	return categoryEngName + common.UNDER_SCORE + strconv.Itoa(int(categoryId))
}

func vo2Do(vo vo.ConditionVO) model.Category {
	var Do model.Category
	Do.ID = vo.CategoryId
	Do.Name = vo.Name
	Do.EngName = vo.EngName
	Do.ParentId = vo.ParentId
	Do.ParentName = vo.ParentName
	return Do
}
