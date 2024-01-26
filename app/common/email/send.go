package email

import (
	"fmt"
	"net/smtp"
	"remembrance/app/common/tool"
)

type Message struct {
	Email string `json:"email"`   //邮箱
	Code  string `json:"code"`    //验证码
	Type  string `json:"gettype"` //请求类型 注册"register" 改密码"change"
}

const (
	qqsmtp    = "smtp.qq.com"
	qqport    = "587"
	from      = "2804366305@qq.com" //发件人地址
	form_code = "srqjyiexqmdvdchh"  //授权码
)

func SendCode(to string) string {
	//生成验证码
	code := tool.Randnum(6)
	fmt.Println(code)
	// 邮件主题和内容
	subject := "Code"
	body := code
	// 邮件内容
	message := "Subject: " + subject + "\r\n" +
		"To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"\r\n" +
		body

		// 使用 SMTP 连接到 QQ 邮箱服务器
	auth := smtp.PlainAuth("", from, form_code, qqsmtp)
	err := smtp.SendMail(qqsmtp+":"+qqport, auth, from, []string{to}, []byte(message))
	if err != nil {
		// 处理发送邮件时的错误
		panic(err)
	}

	// 邮件发送成功
	println("Email sent successfully")
	return code
}
