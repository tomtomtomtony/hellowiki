package service

import (
	"hellowiki/common/result"
	"hellowiki/model"
)

func CreateCategory(category *model.Category) (code int) {
	codeInsert := model.CreateCategory(category)
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
	//若被删除节点为根节点
	if category.ID == category.ParentId {
		for _, curr := range children {
			curr.ParentId, curr.ParentName = curr.ID, curr.Name
		}
		//若为非根节点
	} else {
		for _, curr := range children {
			curr.ParentId, curr.ParentName = category.ParentId, category.ParentName
		}
	}

	return model.DeleteCategoryById(category.ID)

}

func SetCategory(id uint, data model.Category) int {
	if model.HasCategoryById(id) == 200 {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	return model.UpdateCategoryById(id, data)
}
