package utils

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"hellowiki/config"
	"log"
	"os"
	"time"
)

func OpenDB() *gorm.DB {
	dbBase, err := gorm.Open(sqlite.Open(config.Cfg.DataBase.AbsPath+string(os.PathSeparator)+config.Cfg.DataBase.DefaultDBName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	if err != nil {
		log.Fatalf("fail to connect database:{%v}", err)
	}
	return dbBase
}
