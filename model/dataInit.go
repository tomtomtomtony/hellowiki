package model

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"hellowiki/common/utils"
	"hellowiki/config"
	"hellowiki/model/searchConfig"
	utils2 "hellowiki/model/utils"
	"log"
	"os"
	"time"
)

func init() {
	//初始化数据目录
	check, err := utils.HasDirectoryOrFile(config.Cfg.DataDir.Location)
	if err != nil {
		log.Printf("未能检测到data文件夹:{%v}\n", err)
	}
	if !check {
		if err := os.Mkdir(config.Cfg.DataDir.Location, os.ModePerm); err != nil {
			log.Fatalf("创建文件夹{%v}失败:{%v}", config.Cfg.DataDir.Location, err)
		}
	}

	initDataBase()
	initIndexDir()
	initContentDir()

}

func initDataBase() {
	//初始化数据库目录
	config.Cfg.DataBase.AbsPath = config.Cfg.DataDir.Location + string(os.PathSeparator) + config.Cfg.DataBase.Location
	check, err := utils.HasDirectoryOrFile(config.Cfg.DataBase.AbsPath)
	if err != nil {
		log.Printf("未能检测到data文件夹:{%v}", err)
	}
	if !check {
		if err := os.Mkdir(config.Cfg.DataBase.AbsPath, os.ModePerm); err != nil {
			log.Fatalf("创建db文件夹失败:{%v}", err)
		}
	}
	//初始化数据库
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
	dbBase.Set("gorm:table_options", "AUTO_INCREMENT=1")
	err = dbBase.AutoMigrate(RegUser{})
	if err != nil {
		log.Fatalf("建表失败")
	}
}

func initIndexDir() {
	//初始化索引目录
	config.Cfg.SearchDB.AbsPath = config.Cfg.DataDir.Location + string(os.PathSeparator) + config.Cfg.SearchDB.Location
	check, err := utils.HasDirectoryOrFile(config.Cfg.SearchDB.AbsPath)
	if err != nil {
		log.Printf("未能检测到index文件夹:{%v}", err)
	}
	if !check {
		if err := os.Mkdir(config.Cfg.SearchDB.AbsPath, os.ModePerm); err != nil {
			log.Fatalf("创建文件夹{%v}失败:{%v}", config.Cfg.SearchDB.AbsPath, err)
		}
	}

	//初始化索引
	_, err = bleve.Open(config.Cfg.SearchDB.AbsPath + string(os.PathSeparator) + config.Cfg.SearchDB.DefaultIndex)
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
			"type":      searchConfig.TokenName,
			"tokenizer": searchConfig.TokenName,
		}

		articlesMapping := utils2.BuildArticleMapping(tokenOpt)
		_, err = bleve.New(config.Cfg.SearchDB.AbsPath+string(os.PathSeparator)+config.Cfg.SearchDB.DefaultIndex, articlesMapping)
		if err != nil {
			log.Fatal("创建索引失败", err)
		}
	}
}

func initContentDir() {
	//初始化文章存储
	config.Cfg.DirDB.AbsPath = config.Cfg.DataDir.Location + string(os.PathSeparator) + config.Cfg.DirDB.Location
	check, err := utils.HasDirectoryOrFile(config.Cfg.DirDB.AbsPath)
	if err != nil {
		log.Printf("未能检测到index文件夹:{%v}", err)
	}
	if !check {
		if err := os.Mkdir(config.Cfg.DirDB.AbsPath, os.ModePerm); err != nil {
			log.Fatalf("创建文件夹{%v}失败:{%v}", config.Cfg.DirDB.AbsPath, err)
		}
	}

	//初始化默认文章目录
	_, err = os.OpenFile(config.Cfg.DirDB.AbsPath+string(os.PathSeparator)+config.Cfg.DirDB.DefaultDir, os.O_RDONLY, os.ModeDir)
	if err != nil {
		if err != os.ErrNotExist {
			log.Printf("目录{%v}不存在\n", config.Cfg.DirDB.AbsPath+string(os.PathSeparator)+config.Cfg.DirDB.DefaultDir)
		}
		err = os.Mkdir(config.Cfg.DirDB.AbsPath+string(os.PathSeparator)+config.Cfg.DirDB.DefaultDir, os.ModePerm)
		if err != nil {
			log.Fatal("创建txt默认目录失败", err)
		}
	}
}
