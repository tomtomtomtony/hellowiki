package model

import "gorm.io/gorm"

type CatGory struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}
