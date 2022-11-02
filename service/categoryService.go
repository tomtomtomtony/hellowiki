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
	var newData model.Category
	//1.若被删除节点为根节点
	if category.ID == category.ParentId {
		tx := model.Db.Begin()
		//1.1更新每个孩子节点的父节点Id和名称
		for _, curr := range children {
			newData.ParentId, newData.ParentName = curr.ID, curr.Name
			tx.Model(&model.Category{}).Where("id=?", curr.ID).Updates(newData)
		} //1.2删除节点本身
		tx.Delete(&model.Category{}, "id=?", category.ID)
		err := tx.Rollback()
		if err != nil {
			return result.ERROR
		}
		tx.Commit()
		return result.SUCCSE
		//2.若为非根节点
	} else {
		tx := model.Db.Begin()
		//2.1更新每个孩子节点的父节点Id和名称为待删除节点的父节点Id和名称
		for _, curr := range children {
			newData.ParentId, newData.ParentName = category.ParentId, category.ParentName
			tx.Model(&model.Category{}).Where("id=?", curr.ID).Updates(newData)
		} //2.2删除节点本身
		tx.Delete(&model.Category{}, "id=?", category.ID)
		err := tx.Rollback()
		if err != nil {
			return result.ERROR
		}
		tx.Commit()
	}
	return result.SUCCSE
}

func SetCategory(id uint, data model.Category) int {
	if model.HasCategoryById(id) == 200 {
		return result.ERROR_CATEGORY_NOT_FOUND
	}
	return model.UpdateCategoryById(id, data)
}
