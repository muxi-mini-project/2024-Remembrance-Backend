package email

import (
	"fmt"
	"net/smtp"
	"remembrance/app/common"
	"remembrance/app/common/tool"
	"time"

	"github.com/go-redis/redis"
)

type Message struct {
	Email string `json:"email"`   //邮箱
	Code  string `json:"code"`    //验证码
	Type  string `json:"gettype"` //请求类型 注册"register" 改密码"change"
}

//  common.CONFIG.Ema.Qqsmtp
//  common.CONFIG.Ema.Qqport
//  common.CONFIG.Ema.From      //发件人地址
//  common.CONFIG.Ema.Form_code //授权码

// 发送验证码
func SendCode(to string, way string) {
	//生成验证码
	code := tool.Randnum(6)
	fmt.Println(code)
	// 邮件主题和内容
	subject := "Code"
	body := code
	// 邮件内容
	message := "Subject: " + subject + "\r\n" +
		"To: " + to + "\r\n" +
		"From: " + common.CONFIG.Ema.From + "\r\n" +
		"\r\n" +
		body

		// 使用 SMTP 连接到 QQ 邮箱服务器
	auth := smtp.PlainAuth("", common.CONFIG.Ema.From, common.CONFIG.Ema.Form_code, common.CONFIG.Ema.Qqsmtp)
	err := smtp.SendMail(common.CONFIG.Ema.Qqsmtp+":"+common.CONFIG.Ema.Qqport, auth, common.CONFIG.Ema.From, []string{to}, []byte(message))
	if err != nil {
		// 处理发送邮件时的错误
		panic(err)
	}

	// 邮件发送成功
	println("Email sent successfully")
	//检查方法
	switch way {
	case "register":
		err = common.RDB.Set("verification"+to+way, code, 5*time.Minute).Err() // 验证码有效期为5分钟
	case "change":
		err = common.RDB.Set("verification"+to+way, code, 10*time.Minute).Err() // 验证码有效期为10分钟
	}
	if err != nil {
		panic(err)
	}
	return
}

// 检查验证码
func (mes Message) CheckCode() (int, string) {
	val, err := common.RDB.Get("verification" + mes.Email + mes.Type).Result()
	if err == redis.Nil {
		fmt.Println("验证码不存在或已过期")
		return 400, "验证码不存在或已过期"
	} else if err != nil {
		panic(err)
		return 400, "错误"
	} else if val == mes.Code {
		fmt.Println("验证码正确")
		return 200, "验证码正确"
	} else {
		fmt.Println("验证码错误")
		return 400, "验证码错误"
	}
}
