package model

import (
	"gorm.io/gorm"
	utils2 "hellowiki/common/utils"
	"log"
)

/*
*
type字段用于标识菜单类型，其中1表示分类，2表示文章
*/
type Menu struct {
	gorm.Model
	Name       string `gorm:"type:varchar(40);not null" json:"name"`
	ParentId   uint   `gorm:"type:int;not null" json:"parentId"`
	ParentName string `gorm:"type:varchar(40);not null" json:"parentName"`
	Type       uint8  `gorm:"type:varchar(10);not null" json:"type"`
}

func GetAllDirectMenuChildren(id uint) []Menu {
	var menus []Menu
	dbBase := utils2.OpenDB()
	err := dbBase.Limit(500).Where("parent_id=?", id).Find(&menus).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("查找id为{%v}的直接子类时，出现错误:{%v}\n", id, err)
		return []Menu{}
	}
	return menus
}

func GetTopLevelMenu() []Menu {
	var menus []Menu
	dbBase := utils2.OpenDB()
	err := dbBase.Limit(500).Where("parent_id=0").Find(&menus).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("查询顶级分类时出现错误:{%v}\n", err)
		return []Menu{}
	}
	return menus
}
