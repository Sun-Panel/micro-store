package mail

import "gopkg.in/gomail.v2"

type EmailInfo struct {
	Username string // 账号
	Password string // 密码
	Host     string // 服务器地址
	Port     int    // 端口 默认465
}

type Emailer struct {
	EmailInfo EmailInfo
	Dialer    *gomail.Dialer
}

type ContentInfo struct {
	MailTo     []string
	SendName   string
	Title      string
	Body       string
	RespaceArg map[string]string
}
