package config

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"hellowiki/common/utils"
	"hellowiki/model"
	"hellowiki/service/search"
	utils2 "hellowiki/service/utils"
	"log"
	"os"
	"time"
)

var DbBase *gorm.DB

func InitData() {
	//初始化数据目录
	check, err := utils.HasDirectory(Cfg.DataDir.Location)
	if err != nil {
		log.Fatalf("未能检查data文件夹:{%v}", err)
	}
	if !check {
		if err := os.Mkdir(Cfg.DataDir.Location, os.ModePerm); err != nil {
			log.Fatalf("创建文件夹{%v}失败:{%v}", Cfg.DataDir.Location, err)
		}
	}

	initDataBase()
	initIndexDir()
	initContentDir()

}

func initDataBase() {
	//初始化数据库目录
	Cfg.DataBase.AbsPath = Cfg.DataDir.Location + string(os.PathSeparator) + Cfg.DataBase.Location
	check, err := utils.HasDirectory(Cfg.DataBase.AbsPath)
	if err != nil {
		log.Fatalf("未能检查data文件夹:{%v}", err)
	}
	if !check {
		if err := os.Mkdir(Cfg.DataBase.AbsPath, os.ModePerm); err != nil {
			log.Fatalf("创建db文件夹失败:{%v}", err)
		}
	}
	//初始化数据库
	DbBase, err = gorm.Open(sqlite.Open(Cfg.DataBase.AbsPath+string(os.PathSeparator)+Cfg.DataBase.DefaultDBName), &gorm.Config{
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
	DbBase.Set("gorm:table_options", "AUTO_INCREMENT=1")
	err = DbBase.AutoMigrate(&model.RegUser{}, &model.Category{})
	if err != nil {
		log.Fatalf("建表失败")
	}
}

func initIndexDir() {
	//初始化索引目录
	Cfg.SearchDB.AbsPath = Cfg.DataDir.Location + string(os.PathSeparator) + Cfg.SearchDB.Location
	check, err := utils.HasDirectory(Cfg.SearchDB.AbsPath)
	if err != nil {
		log.Fatalf("未能检查index文件夹:{%v}", err)
	}
	if !check {
		if err := os.Mkdir(Cfg.SearchDB.AbsPath, os.ModePerm); err != nil {
			log.Fatalf("创建文件夹{%v}失败:{%v}", Cfg.SearchDB.AbsPath, err)
		}
	}

	//初始化索引
	_, err = bleve.Open(Cfg.SearchDB.AbsPath + string(os.PathSeparator) + Cfg.SearchDB.DefaultIndex)
	if err != nil {
		if err != bleve.ErrorIndexPathDoesNotExist {
			log.Fatalf("搜索索引失败")
		}

		tokenOpt := map[string]interface{}{
			"dicts":     Cfg.Analyze.Dict,
			"stop":      "",
			"opt":       "search-hmm",
			"trim":      "trim",
			"alpha":     false,
			"type":      search.TokenName,
			"tokenizer": search.TokenName,
		}

		articlesMapping := utils2.BuildArticleMapping(tokenOpt)
		_, err = bleve.New(Cfg.SearchDB.AbsPath+string(os.PathSeparator)+Cfg.SearchDB.DefaultIndex, articlesMapping)
		if err != nil {
			log.Fatal("创建索引失败", err)
		}
	}
}

func initContentDir() {
	//初始化文章存储
	Cfg.DirDB.AbsPath = Cfg.DataDir.Location + string(os.PathSeparator) + Cfg.DirDB.Location
	check, err := utils.HasDirectory(Cfg.DirDB.AbsPath)
	if err != nil {
		log.Fatalf("未能检查index文件夹:{%v}", err)
	}
	if !check {
		if err := os.Mkdir(Cfg.DirDB.AbsPath, os.ModePerm); err != nil {
			log.Fatalf("创建文件夹{%v}失败:{%v}", Cfg.DirDB.AbsPath, err)
		}
	}

	//初始化默认文章目录
	_, err = os.OpenFile(Cfg.DirDB.AbsPath+string(os.PathSeparator)+Cfg.DirDB.DefaultDir, os.O_RDONLY, os.ModeDir)
	if err != nil {
		if err != os.ErrNotExist {
			log.Printf("目录{%v}不存在\n", Cfg.DirDB.AbsPath+string(os.PathSeparator)+Cfg.DirDB.DefaultDir)
		}
		err = os.Mkdir(Cfg.DirDB.AbsPath+string(os.PathSeparator)+Cfg.DirDB.DefaultDir, os.ModePerm)
		if err != nil {
			log.Fatal("创建txt默认目录失败", err)
		}
	}
}
