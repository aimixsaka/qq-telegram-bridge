package config

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

type Config struct {
	QQBot  QQBotConfig
	TGBot  TGBotConfig
	Groups map[string]Group
}

type Group struct {
	QQ uint32 `toml:"qq"`
	TG int64  `toml:"tg"`
}

// BotConfig 代表TOML文件中的bot部分
type QQBotConfig struct {
	Account    uint32 `toml:"account"`
	Password   string `toml:"password"`
	SignServer string `toml:"signServer"`
}

type TGBotConfig struct {
	Token string `toml:"token"`
}

// GlobalConfig 默认全局配置
var GlobalConfig *Config
var GlobalQQTGMap map[uint32]int64
var GlobalTGQQMap map[int64]uint32

func constructQQTGMap() {
	logrus.Infof("Groups: %v\n", GlobalConfig.Groups)
	GlobalQQTGMap = make(map[uint32]int64, len(GlobalConfig.Groups))
	GlobalTGQQMap = make(map[int64]uint32, len(GlobalConfig.Groups))
	for _, pair := range GlobalConfig.Groups {
		GlobalQQTGMap[pair.QQ] = pair.TG
		GlobalTGQQMap[pair.TG] = pair.QQ
	}
}

// Init 使用 ./application.toml 初始化全局配置
func Init() {
	GlobalConfig = &Config{}
	_, err := toml.DecodeFile("application.toml", GlobalConfig)
	if err != nil {
		logrus.WithField("config", "GlobalConfig").
			WithError(err).
			Panicf("unable to read global config")
	}
	constructQQTGMap()
	logrus.Infof("GlobalQQTGMap: %v\n", GlobalQQTGMap)
	logrus.Infof("GlobalTGQQMap: %v\n", GlobalTGQQMap)

}

// InitWithContent 从字节数组中读取配置内容
func InitWithContent(configTOMLContent []byte) {
	_, err := toml.Decode(string(configTOMLContent), GlobalConfig)
	if err != nil {
		panic(err)
	}
}
