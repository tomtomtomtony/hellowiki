package model

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"hellowiki/config"
	"hellowiki/service/search"
	utils2 "hellowiki/service/utils"
	"log"
	"os"
	"time"
)

var DbBase *gorm.DB
var err error

func InitDB() {
	DbBase, err = gorm.Open(sqlite.Open(config.Cfg.DataBase.Location), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	DbBase.Set("gorm:table_options", "AUTO_INCREMENT=1")
	if err != nil {
		panic("fail to connect database")
	}
	err = DbBase.AutoMigrate(&RegUser{}, &Category{})
	if err != nil {
		panic("建表失败")
	}

	//初始化索引目录
	_, err = os.OpenFile(config.Cfg.SearchDB.Location, os.O_RDONLY, os.ModeDir)
	if err != nil {
		log.Printf("打开{%v}文件夹失败\n", config.Cfg.SearchDB.Location)
		if err == os.ErrNotExist {
			log.Printf("{%v}不存在\n", config.Cfg.SearchDB.Location)
		}
		err = os.Mkdir(config.Cfg.SearchDB.Location, os.ModePerm)
		if err != nil {
			panic("创建文章目录失败")
		}
	}
	//初始化索引
	config.Cfg.SearchDB.IndexPathName = config.Cfg.SearchDB.Location + config.Cfg.SearchDB.DefaultIndex
	_, err := bleve.Open(config.Cfg.SearchDB.IndexPathName)
	if err != nil {
		if err != bleve.ErrorIndexPathDoesNotExist {
			log.Fatalf("搜索索引失败")
		}

		tokenOpt := map[string]interface{}{
			"dicts":     config.Cfg.Analyze.Dict,
			"stop":      "",
			"opt":       "search-hmm",
			"trim":      "trim",
			"alpha":     false,
			"type":      search.TokenName,
			"tokenizer": search.TokenName,
		}

		articlesMapping := utils2.BuildArticleMapping(tokenOpt)
		_, err = bleve.New(config.Cfg.SearchDB.IndexPathName, articlesMapping)
		if err != nil {
			log.Fatalf("创建索引失败", err)
		}
	}

	//初始化txt存储目录
	_, err = os.OpenFile(config.Cfg.DirDB.Location, os.O_RDONLY, os.ModeDir)
	if err != nil {
		log.Printf("打开{%v}文件夹失败\n", config.Cfg.DirDB.Location)
		if err == os.ErrNotExist {
			log.Printf("{%v}不存在\n", config.Cfg.DirDB.Location)
		}
		err = os.Mkdir(config.Cfg.DirDB.Location, os.ModePerm)
		if err != nil {
			panic("创建文章目录失败")
		}
	}

	//初始化默认txt目录
	config.Cfg.DirDB.TxtPathName = config.Cfg.DirDB.Location + config.Cfg.DirDB.DefaultDir
	_, err = os.OpenFile(config.Cfg.DirDB.TxtPathName, os.O_RDONLY, os.ModeDir)
	if err != nil {
		if err != os.ErrNotExist {
			log.Printf("目录{%v}不存在\n", config.Cfg.DirDB.TxtPathName)
		}
		err = os.Mkdir(config.Cfg.DirDB.TxtPathName, os.ModePerm)
		if err != nil {
			log.Fatalf("创建txt默认目录{%v}失败", config.Cfg.DirDB.TxtPathName)
		}
	}
}
