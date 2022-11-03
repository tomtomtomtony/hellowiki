package model

import (
	"gorm.io/gorm"
	"hellowiki/common/result"
)

// 文章分类
type Category struct {
	gorm.Model
	Name       string `gorm:"type:varchar(40);not null" json:"name"`
	ParentId   uint   `gorm:"type:int;not null" json:"parentId"`
	ParentName string `gorm:"type:varchar(40);not null" json:"parentName"`
}

var (
	//顶级父类的parentId
	TOPLEVELCATEGORY uint = 0
)

// 顶级父类id为0
func FindCategoryChildren(id uint) []Category {
	var categories []Category
	err := Db.Limit(500).Where("parent_id=?", id).Find(&categories).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []Category{}
	}
	return categories
}

func HasCategoryById(id uint) (code int) {
	var category Category
	Db.Take(&category, "id=?", id)
	if category.ID > 0 {
		//分类已存在
		return result.ERROR_CATEGORY_EXIST
	}
	//分类不存在
	return result.SUCCSE
}

// 新增分类数据
func CreateCategory(data Category) (code int) {
	err := Db.Create(&data).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 查询分类列表
func FindAllCategory(pageSize int, pageNum int) []Category {
	var categories []Category
	err := Db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&categories).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []Category{}
	}
	return categories
}

// 根据id，软删除分类信息
func DeleteCategoryById(id uint) int {
	err := Db.Delete(&Category{}, "id=?", id).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 根据id，更新分类信息
func UpdateCategoryById(id uint, category Category) int {
	err := Db.Model(&category).Where("id=?", id).Updates(category).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 根据parentId更新
func UpdateCategoryByParentId(parentId uint, category Category) (int, error) {
	err := Db.Model(&category).Where("parent_id=?", parentId).Updates(category).Error
	if err != nil {
		return result.ERROR, err
	}
	return result.SUCCSE, err
}
