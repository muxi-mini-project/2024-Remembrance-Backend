package common

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	CONFIG Config
)

type Config struct {
	Dsn string
	Oss struct {
		AccessKey  string
		SecretKey  string
		Bucket string
		Domain string
	}
}

func Loadconfig() {
	viper.SetConfigName("config") // 配置文件的名称（没有文件扩展名）
	viper.SetConfigType("yaml")   // 配置文件的类型，例如 "yaml"
	viper.AddConfigPath(".")      // 配置文件的路径

	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&CONFIG) // 将配置文件内容赋值给 Config 结构体
	if err != nil {
		log.Fatalf("Unable to decode into struct, %s", err)
	}
}
