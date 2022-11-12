package service

import (
	"hellowiki/api/result"
	"hellowiki/api/v1/category/vo"
	"hellowiki/common/utils"
	"hellowiki/config"
	"hellowiki/model"
	"hellowiki/model/searchConfig"
	utils2 "hellowiki/model/utils"
	"log"
	"os"
)

func CreateCategory(condition vo.ConditionVO) (code int) {
	var data model.Category
	var err error
	code = preTreatment(condition)
	if code != result.SUCCSE {
		return code
	}
	data = vo2Do(condition)
	//写入数据库
	code, dataId := model.CategoryWriteToDBMenuTable(vo2MenuType(condition))
	if code != result.SUCCSE {
		log.Println("新增菜单写入数据库失败")
		return result.ERROR
	}

	//写入索引文件
	indexName := utils2.ConstructStandardIndexName(data.Name, dataId)

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
	code = model.WriteToCategoryIndex(indexName, articlesMapping)
	if code != result.SUCCSE {
		return code
	}

	//写入磁盘文件夹
	dirName := indexName
	err = os.Mkdir(config.Cfg.DirDB.AbsPath+string(os.PathSeparator)+dirName, os.ModePerm)
	if err != nil {
		log.Println("新增分类写入磁盘失败")
		return result.ERROR
	}

	return result.SUCCSE
}

func preTreatment(categoryInfo vo.ConditionVO) int {
	if categoryInfo.ParentId != model.TOPLEVELCATEGORY && !HasCategory(categoryInfo.ParentId) {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	children := model.FindDirectChildren(categoryInfo.ParentId)
	for _, curr := range children {
		if curr.Name == categoryInfo.Name {
			return result.ERROR_CATEGORY_EXIST
		}
	}
	return result.SUCCSE
}

func GetAllCategory(pageSize int, pageNum int) []model.Category {
	return model.FindAllCategory(pageSize, pageNum)
}

func DeleteCategory(category model.Category) int {
	children := model.FindDirectChildren(category.ID)
	var newData model.Menu
	dbBase := utils.OpenDB()

	tx := dbBase.Begin()
	//1.更新每个孩子节点的父节点Id和名称
	for _, curr := range children {
		//1.1若为待删节点为根节点，其直接子节点将成为顶级父节点
		if category.ParentId == model.TOPLEVELCATEGORY {
			newData.ParentId, newData.ParentName = model.TOPLEVELCATEGORY, curr.Name
		} else {
			newData.ParentId, newData.ParentName = category.ParentId, category.ParentName
		}
		if err := tx.Model(&model.Menu{}).Where("id=?", curr.ID).Updates(newData).Error; err != nil {
			tx.Rollback()
			return result.ERROR
		}
	} //2.删除节点本身
	if err := tx.Delete(&model.Menu{}, "id=?", category.ID).Error; err != nil {
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
	return model.GetCategoryById(id) != model.Menu{}
}

// 传入索引名称符合格式: article的 engName_categoryId
func HasCategoryInIndex(indexName string) bool {
	check, _ := model.HasCategoryInIndexDir(indexName)
	return check
}

//func repairTable(categoryId uint) int {
//	var category model.Category
//	category = model.GetCategoryById(categoryId)
//
//}

func vo2Do(vo vo.ConditionVO) model.Category {
	var Do model.Category
	Do.ID = vo.CategoryId
	Do.Name = vo.Name
	Do.ParentId = vo.ParentId
	Do.ParentName = vo.ParentName
	return Do
}

func vo2MenuType(vo vo.ConditionVO) model.Menu {
	var Do model.Menu
	Do.Name = vo.Name
	Do.ParentId = vo.ParentId
	Do.ParentName = vo.ParentName
	Do.Type = 1
	return Do
}
