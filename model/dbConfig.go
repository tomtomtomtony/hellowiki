package model

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"hellowiki/config"
	"time"
)

var Db *gorm.DB
var err error

func InitDB() {
	Db, err = gorm.Open(sqlite.Open(config.Cfg.DataBase.Location), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	Db.Set("gorm:table_options", "AUTO_INCREMENT=1")
	if err != nil {
		panic("fail to connect database")
	}
	err = Db.AutoMigrate(&RegUser{}, &Category{}, &Article{})
	if err != nil {
		panic("建表失败")
	}

}
