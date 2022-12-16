package model

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
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

	initIndexDir()
	initContentDir()
	initAuthenticationDB()
	initDataBase()

}

func initDataBase() {
	//初始化数据库目录
	config.Cfg.DataBase.AbsPath = config.Cfg.DataDir.Location + string(os.PathSeparator) + config.Cfg.DataBase.Location
	check, err := utils.HasDirectoryOrFile(config.Cfg.DataBase.AbsPath)
	if err != nil {
		log.Printf("未能检测到db文件夹:{%v}", err)
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
		log.Fatalf("连接数据库失败:{%v}", err)
	}
	dbBase.Set("gorm:table_options", "AUTO_INCREMENT=1")

	err = dbBase.AutoMigrate(RegUser{})
	if err != nil {
		log.Fatalf("建表失败")
	}
	var root RegUser
	//创建超级管理员,仅在程序第一次运行时起作用
	var userNumber int64
	dbBase.Model(&RegUser{}).Count(&userNumber)
	if userNumber < 1 {
		root.UserName = config.Cfg.SuperAdmin.UserName
		hashWord, err := bcrypt.GenerateFromPassword([]byte(config.Cfg.SuperAdmin.PassWord), bcrypt.DefaultCost)
		root.PassWord = string(hashWord)
		root.IsEnable = true
		err = dbBase.Create(&root).Error
		if err != nil {
			log.Fatalf("创建超级管理员失败:{%v}", err)
		}
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
		log.Printf("未能检测到content文件夹:{%v}", err)
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

func initAuthenticationDB() {
	//初始化权限管理数据
	config.Cfg.AuthenticationDB.AbsPath = config.Cfg.DataDir.Location + string(os.PathSeparator) + config.Cfg.AuthenticationDB.Location
	check, err := utils.HasDirectoryOrFile(config.Cfg.AuthenticationDB.AbsPath)
	if err != nil {
		log.Printf("未能检测到authenDB文件夹:{%v}", err)
	}
	if !check {
		if err := os.Mkdir(config.Cfg.AuthenticationDB.AbsPath, os.ModePerm); err != nil {
			log.Fatalf("创建文件夹{%v}失败:{%v}", config.Cfg.AuthenticationDB.AbsPath, err)
		}
	}

	//初始化默认文章目录
	//创建鉴权模型文件
	modelFile := config.Cfg.AuthenticationDB.AbsPath + string(os.PathSeparator) + config.Cfg.AuthenticationDB.ModelFile
	file, err := os.Open(modelFile)
	defer func() { file.Close() }()
	if err != nil && os.IsNotExist(err) {
		newFile, err := os.Create(modelFile)
		if err != nil {
			log.Fatalf("创建鉴权模型文件失败:{%v}", err)
		}
		defer func() { newFile.Close() }()
		newFile.WriteString("[request_definition]" + "\n")
		newFile.WriteString(config.Cfg.AuthenticationDB.RequestDefinition + "\n")
		newFile.WriteString("[policy_definition]" + "\n")
		newFile.WriteString(config.Cfg.AuthenticationDB.PolicyDefinition + "\n")
		newFile.WriteString("[matchers]" + "\n")
		newFile.WriteString(config.Cfg.AuthenticationDB.Matchers + "\n")
		newFile.WriteString("[policy_effect]" + "\n")
		newFile.WriteString(config.Cfg.AuthenticationDB.PolicyEffect + "\n")
		newFile.WriteString("[role_definition]" + "\n")
		newFile.WriteString(config.Cfg.AuthenticationDB.RoleDefinition + "\n")
	}
	//创建鉴权策略文件
	policyFile := config.Cfg.AuthenticationDB.AbsPath + string(os.PathSeparator) + config.Cfg.AuthenticationDB.PolicyFile
	file, err = os.Open(policyFile)
	defer func() { file.Close() }()
	if err != nil && os.IsNotExist(err) {
		_, err = os.Create(policyFile)
		if err != nil {
			log.Fatalf("创建鉴权策略文件失败:{%v}", err)
		}
		csvAdapter := fileadapter.NewAdapter(policyFile)
		enforcer, err := casbin.NewEnforcer(modelFile, csvAdapter)
		if err != nil {
			log.Fatalf("创建鉴权器失败:{%v}", err)
		}
		//将用户数据库第一位设为超级管理员
		//“1”表示超级管理员在用户数据库中的ID
		enforcer.AddRoleForUser("1", config.Cfg.SuperAdmin.Role)
		enforcer.SavePolicy()

	}

}
