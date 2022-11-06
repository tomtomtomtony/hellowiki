package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type tomlConfig struct {
	DataBase dataBase
	SearchDB searchDB
	DirDB    dirDB
	WikiInfo wikiInfo
	Server   server
	Analyze  analyze
}

type dataBase struct {
	userName string
	password string
	Location string
}

type searchDB struct {
	Location      string
	DefaultIndex  string
	IndexPathName string
}

type dirDB struct {
	Location    string
	DefaultDir  string
	TxtPathName string
}

type analyze struct {
	Dict string
}

type wikiInfo struct {
	WikiName string
}

type server struct {
	AppMode string
	Port    string
}

var Cfg *tomlConfig

func init() {
	Cfg = new(tomlConfig)
	//读取配置文件
	viper.SetConfigName("config") //设置文件名时不要带后缀
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config") //搜索路径可以设置多个，viper 会根据设置顺序依次查找
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}
	//输出读取到的值
	//2.获取所有值
	fmt.Println("all settings: ", viper.AllSettings())
	//3.映射到结构体
	err := viper.Unmarshal(&Cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

}
