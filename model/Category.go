package model

import (
	"gorm.io/gorm"
	"hellowiki/common/result"
	"hellowiki/config"
	"log"
	"os"
)

// 文章分类
type Category struct {
	gorm.Model
	Name       string `gorm:"type:varchar(40);not null" json:"name"`
	EngName    string `gorm:"type:varchar(40);not null "json:"engName"`
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
	err := config.DbBase.Limit(500).Where("parent_id=?", id).Find(&categories).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("查找id为{%v}的直接子类时，出现错误:{%v}\n", id, err)
		return []Category{}
	}
	return categories
}

func GetCategoryById(id uint) Category {
	var category Category
	err := config.DbBase.Take(&category, "id=?", id).Error
	if err != nil {
		return Category{}
	}
	return category
}

func HasCategoryIndex(indexName string) bool {
	_, err := os.OpenFile(config.Cfg.SearchDB.Location+indexName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return false
	}
	return true
}

func HasCategoryTable(tableName string) bool {
	return config.DbBase.Migrator().HasTable(tableName)
}

// 新增分类数据
func CreateCategory(data Category) (code int, id uint) {
	err := config.DbBase.Create(&data).Error
	if err != nil {
		return result.ERROR, 0
	}
	return result.SUCCSE, data.ID
}

// 查询分类列表
func FindAllCategory(pageSize int, pageNum int) []Category {
	var categories []Category
	err := config.DbBase.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&categories).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []Category{}
	}
	return categories
}

// 根据id，软删除分类信息
func DeleteCategoryById(id uint) int {
	err := config.DbBase.Delete(&Category{}, "id=?", id).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 根据id，更新分类信息
func UpdateCategoryById(id uint, category Category) int {
	err := config.DbBase.Model(&category).Where("id=?", id).Updates(category).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 根据parentId更新
func UpdateCategoryByParentId(parentId uint, category Category) (int, error) {
	err := config.DbBase.Model(&category).Where("parent_id=?", parentId).Updates(category).Error
	if err != nil {
		return result.ERROR, err
	}
	return result.SUCCSE, err
}
