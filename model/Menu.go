package model

import "gorm.io/gorm"

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
