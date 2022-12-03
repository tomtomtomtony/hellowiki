package service

import (
	"hellowiki/api/result"
	"hellowiki/api/v1/category/vo"
	"hellowiki/common"
	"hellowiki/common/utils"
	"hellowiki/config"
	"hellowiki/model"
	"hellowiki/model/searchConfig"
	utils2 "hellowiki/model/utils"
	"log"
)

func CreateCategory(condition vo.ConditionVO) (code int) {
	var data = vo2Category(condition)
	code = preTreatment(condition)
	if code != result.SUCCSE {
		return code
	}

	//写入数据库
	code, dataId := model.CategoryWriteToDBMenuTable(vo2MenuType(condition))
	if code != result.SUCCSE {
		log.Println("新增菜单写入数据库失败")
		return result.ERROR
	}

	//写入索引文件
	indexName := utils2.ConstructCategoryNameId(data.Name, dataId)

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
	code = model.WriteToCategoryContent(dirName)
	if code != result.SUCCSE {
		log.Println("新增分类写入磁盘失败")
		return result.ERROR
	}

	return result.SUCCSE
}

func preTreatment(categoryInfo vo.ConditionVO) int {
	if categoryInfo.ParentMenuId != model.TOPLEVELCATEGORY && !HasCategoryInDBTable(categoryInfo.ParentMenuId) {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	children := model.FindDirectCateGoryChildren(categoryInfo.ParentMenuId)
	for _, curr := range children {
		if curr.Name == categoryInfo.Name {
			return result.ERROR_CATEGORY_EXIST
		}
	}
	return result.SUCCSE
}

func DeleteCategory(condition vo.ConditionVO) int {
	var category = vo2Category(condition)
	children := model.FindDirectCateGoryChildren(category.ID)
	var newData model.Menu
	dbBase := utils.OpenDB()

	tx := dbBase.Begin()
	//1.1更新每个孩子节点的父节点Id和名称
	for _, curr := range children {
		//1.1若为待删节点为根节点，其直接子节点将成为顶级父节点
		if category.ParentMenuId == model.TOPLEVELCATEGORY {
			newData.ParentId, newData.ParentName = model.TOPLEVELCATEGORY, curr.Name
		} else {
			newData.ParentId, newData.ParentName = category.ParentMenuId, category.ParentName
		}
		if err := tx.Model(&model.Menu{}).Where("id=?", curr.ID).Updates(newData).Error; err != nil {
			tx.Rollback()
			return result.ERROR
		}
	} //1.2.删除节点本身
	if err := tx.Delete(&model.Menu{}, "id=?", category.ID).Error; err != nil {
		tx.Rollback()
		return result.ERROR
	}
	tx.Commit()

	//2.删除txt文件
	dirName := utils2.ConstructCategoryNameId(category.Name, category.ID)
	if !utils.FoldIsEmptyInContent(dirName) {
		log.Printf("content文件夹下指定文件夹{%v}不为空", dirName)
		return result.ERROR_CATEGORY_NOT_EMPTY
	}
	code := model.DeleteCategoryInContent(dirName)
	if code != result.SUCCSE {
		log.Println("content文件夹下指定文件夹删除失败")
		return result.ERROR
	}

	//3.删除索引文件
	indexName := dirName
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

func SetCategory(id uint, data model.Category) int {
	if !HasCategoryInDBTable(id) {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	return model.UpdateCategoryById(id, data)
}

func HasCategoryInDBTable(id uint) bool {
	return model.GetCategoryById(id) != model.Menu{}
}

func HasCategoryInContent(categoryName string) bool {
	check, _ := utils.HasCategoryInContentDir(categoryName)
	return check
}

func HasCategoryInIndex(indexName string) bool {
	check, _ := model.HasCategoryInIndexDir(indexName)
	return check
}

func vo2Category(vo vo.ConditionVO) model.Category {
	var Do model.Category
	Do.ID = vo.MenuId
	Do.Name = vo.Name
	Do.ParentMenuId = vo.ParentMenuId
	Do.ParentName = vo.ParentName
	return Do
}

func vo2MenuType(vo vo.ConditionVO) model.Menu {
	var Do model.Menu
	Do.Name = vo.Name
	Do.ParentId = vo.ParentMenuId
	Do.ParentName = vo.ParentName
	Do.Type = common.CATEGORY_TYPE
	return Do
}
