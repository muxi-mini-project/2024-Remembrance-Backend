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

type Config struct {
	Dsn string
	Oss *struct {
		AccessKey string
		SecretKey string
		Bucket    string
		Domain    string
	}
}

// type CaptchaConf struct {
// 	Height          int
// 	Width           int
// 	Length          int
// 	Maxskew         float64
// 	Dotcount        int
// 	ExpireTime      int
// 	DebugExpireTime int
// 	TestingKey      string
// }

// type Config struct {
// 	rest.RestConf
// 	AccountCenterConf zrpc.RpcClientConf
// 	JwtAuth           struct {
// 		AccessSecret string
// 		AccessExpire int64
// 	}
// 	EmailConf *struct {
// 		Host     string
// 		Port     string
// 		UserName string
// 		Password string
// 	}
// 	Oss *struct {
// 		AccessKey  string
// 		SecretKey  string
// 		BucketName string
// 		DomainName string
// 	}
// 	CaptchaConf      *CaptchaConf
// 	EmailCodeExpired int
// 	RedisConf        redis.RedisConf
// 	JwtAuthChPass    struct {
// 		AccessSecret string
// 		AccessExpire int64
// 	}
// }
