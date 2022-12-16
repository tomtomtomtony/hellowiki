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

// 获取所有文章信息，并以列表形式返回
func GetAllArticle(pageSize int, pageNum int) ([]Menu, int64) {
	dbBase := utils2.OpenDB()
	var res []Menu
	var total int64
	err := dbBase.Limit(pageSize).Offset((pageNum - 1) * pageSize).Where("type=2").Find(&res).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("获取全部文章信息时，出现错误:{%v}\n", err)
		return res, total
	}
	return res, total
}
