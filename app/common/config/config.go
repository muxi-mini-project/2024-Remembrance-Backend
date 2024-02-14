package config

import (
	"log"
	"remembrance/app/common"

	"github.com/spf13/viper"
)

func Loadconfig() {
	viper.SetConfigName("config") // 配置文件的名称（没有文件扩展名）
	viper.SetConfigType("yaml")   // 配置文件的类型，例如 "yaml"
	viper.AddConfigPath(".")      // 配置文件的路径

	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&common.CONFIG) // 将配置文件内容赋值给 Config 结构体
	if err != nil {
		log.Fatalf("Unable to decode into struct, %s", err)
	}
}

