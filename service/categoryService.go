package service

import (
	"hellowiki/common"
	"hellowiki/common/result"
	"hellowiki/model"
	"strconv"
)

func CreateCategory(category model.Category) (code int) {
	var data model.Category
	if category.ParentId == model.TOPLEVELCATEGORY {
		code = handleCreateRootCategory(category)
	} else {
		code = handleCreateNonRootCategory(category)
	}
	if code != result.SUCCSE {
		return result.ERROR
	}
	data.Name = category.Name
	data.ParentId = model.TOPLEVELCATEGORY
	data.ParentName = category.Name
	data.EngName = category.EngName
	tx := model.Db.Begin()
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return result.ERROR
	}
	tableName := data.EngName + common.UNDER_SCORE + strconv.Itoa(int(data.ID))
	if err := tx.Table(tableName).AutoMigrate(&model.Article{}); err != nil {
		tx.Rollback()
		return result.ERROR
	}
	tx.Commit()
	return result.SUCCSE
}

func handleCreateRootCategory(categoryInfo model.Category) int {
	children := model.FindCategoryChildren(model.TOPLEVELCATEGORY)
	for _, curr := range children {
		if curr.Name == categoryInfo.Name {
			return result.ERROR_CATEGORY_EXIST
		}
	}
	return result.SUCCSE
}

func handleCreateNonRootCategory(categoryInfo model.Category) int {
	if !HasCategory(categoryInfo.ID) {
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
	tx := model.Db.Begin()
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

//func repairTable(categoryId uint) int {
//	var category model.Category
//	category = model.GetCategoryById(categoryId)
//
//}

func ConstructStandardTableName(categoryEngName string, categoryId uint) string {
	return categoryEngName + common.UNDER_SCORE + strconv.Itoa(int(categoryId))
}
