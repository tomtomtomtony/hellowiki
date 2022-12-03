package model

import (
	"github.com/blevesearch/bleve/v2/mapping"
	"gorm.io/gorm"
	"hellowiki/api/result"
	utils2 "hellowiki/common/utils"
	"hellowiki/config"
	"log"
	"os"
)

// 文章分类
type Category struct {
	gorm.Model
	Name         string `gorm:"type:varchar(40);not null" json:"name"`
	ParentMenuId uint   `gorm:"type:int;not null" json:"parentMenuId"`
	ParentName   string `gorm:"type:varchar(40);not null" json:"parentName"`
}

var (
	//顶级父类的parentId
	TOPLEVELCATEGORY uint = 0
)

// 顶级父类id为0
func FindDirectCateGoryChildren(id uint) []Menu {
	var categories []Menu
	dbBase := utils2.OpenDB()
	err := dbBase.Limit(500).Where("parent_id=? and type=1", id).Find(&categories).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("查找id为{%v}的直接子类时，出现错误:{%v}\n", id, err)
		return []Menu{}
	}
	return categories
}

// 从数据库菜单表中查询分类
func GetCategoryById(id uint) Menu {
	var menu Menu
	dbBase := utils2.OpenDB()
	err := dbBase.Take(&menu, "id=? and type=1", id).Error
	if err != nil {
		return Menu{}
	}
	return menu
}

func HasCategoryTable(tableName string) bool {
	dbBase := utils2.OpenDB()
	return dbBase.Migrator().HasTable(tableName)
}

func HasCategoryInIndexDir(categoryName string) (bool, error) {
	return utils2.HasDirectoryOrFile(config.Cfg.SearchDB.AbsPath + string(os.PathSeparator) + categoryName)
}

//

// 索引写入
func WriteToCategoryIndex(indexName string, mapping mapping.IndexMapping) int {
	return utils2.WriteToIndexDir(indexName, mapping)
}

// data/content下分类文件夹写入
func WriteToCategoryContent(categoryName string) int {
	return utils2.CreateFoldContent(categoryName)
}

// 索引文件夹删除
func DeleteCategoryInIndex(indexName string) int {
	return utils2.DeleteFold(config.Cfg.SearchDB.AbsPath + string(os.PathSeparator) + indexName)
}

func DeleteCategoryInContent(categoryName string) int {
	return utils2.DeleteFold(config.Cfg.DirDB.AbsPath + string(os.PathSeparator) + categoryName)
}

// 根据id，更新分类信息
func UpdateCategoryById(id uint, category Category) int {
	dbBase := utils2.OpenDB()
	err := dbBase.Model(&category).Where("id=?", id).Updates(category).Error
	if err != nil {
		return result.ERROR
	}
	return result.SUCCSE
}

// 根据parentId更新
func UpdateCategoryByParentId(parentId uint, category Category) (int, error) {
	dbBase := utils2.OpenDB()
	err := dbBase.Model(&category).Where("parent_id=?", parentId).Updates(category).Error
	if err != nil {
		return result.ERROR, err
	}
	return result.SUCCSE, err
}

// 数据库写入菜单表
func CategoryWriteToDBMenuTable(data Menu) (code int, id uint) {
	dbBase := utils2.OpenDB()
	err := dbBase.Create(&data).Error
	if err != nil {
		return result.ERROR, 0
	}
	return result.SUCCSE, data.ID
}
