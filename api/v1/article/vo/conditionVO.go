package vo

import "gorm.io/gorm"

type ConditionVO struct {
	gorm.Model
	Name       string `gorm:"-" json:"name"`
	EngName    string `gorm:"-" json:"engName"`
	ParentId   uint   `gorm:"-" json:"parentId"`
	ParentName string `gorm:"-" json:"parentName"`
	PageSize   int    `gorm:"-" json:"pageSize"`
	PageNum    int    `gorm:"-" json:"pageNum"`
}
