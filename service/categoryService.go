package service

import (
	"hellowiki/common/result"
	"hellowiki/model"
)

// 创建根分类
func CreateRootCategory(categoryName string) (code int) {
	var data model.Category
	//查询顶级父类
	children := model.FindCategoryChildren(0)
	for _, curr := range children {
		if curr.Name == categoryName {
			return result.ERROR_CATEGORY_EXIST
		}
	}
	data.Name = categoryName
	data.ParentId = 0
	data.ParentName = categoryName
	codeInsert := model.CreateCategory(data)
	if codeInsert != 200 {
		return codeInsert
	}
	return codeInsert
}

// 创建非根分类
func CreateNonRootCategory(categoryInfo model.Category) (code int) {
	var data model.Category
	data.Name = categoryInfo.Name
	if model.HasCategoryById(categoryInfo.ParentId) == 200 {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	data.ParentId = categoryInfo.ParentId
	data.ParentName = categoryInfo.ParentName
	codeInsert := model.CreateCategory(data)
	if codeInsert != 200 {
		return codeInsert
	}
	return codeInsert
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
		//若为待删节点为根节点，其直接子节点将成为顶级父节点
		if category.ParentId == 0 {
			newData.ParentId, newData.ParentName = 0, curr.Name
		} else {
			newData.ParentId, newData.ParentName = category.ParentId, category.ParentName
		}
		if err := tx.Model(&model.Category{}).Where("id=?", curr.ID).Updates(newData).Error; err != nil {
			tx.Rollback()
			return result.ERROR
		}
	} //2删除节点本身
	if err := tx.Delete(&model.Category{}, "id=?", category.ID).Error; err != nil {
		tx.Rollback()
		return result.ERROR
	}
	tx.Commit()
	return result.SUCCSE
}

func SetCategory(id uint, data model.Category) int {
	if model.HasCategoryById(id) == 200 {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	return model.UpdateCategoryById(id, data)
}
