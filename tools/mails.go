package tools

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/go-gomail/gomail"
	"go.uber.org/zap"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type EmailParam struct {
	// ServerHost 邮箱服务器地址，如腾讯企业邮箱为smtp.exmail.qq.com
	serverHost string
	// ServerPort 邮箱服务器端口，如腾讯企业邮箱为465
	serverPort int
	// FromEmail　发件人邮箱地址
	fromEmail string
	// FromPasswd 发件人邮箱密码（注意，这里是明文形式），
	fromPassword string
	// Toers 接收者邮件，如有多个，则以英文逗号(“,”)隔开，不能为空
	Toers []string
	// CCers 抄送者邮件，如有多个，则以英文逗号(“,”)隔开，可以为空
	CCers []string
	//主题
	Subject string
	//正文
	Text string
}

// SendEmail body支持html格式字符串
func SendEmail(ep *EmailParam) bool {
	m := gomail.NewMessage()
	if len(ep.CCers)+len(ep.Toers) == 0 {
		return false
	}
	m.SetHeader("To", ep.Toers...)
	m.SetHeader("Cc", ep.CCers...)

	// 第三个参数为发件人别名，可以为空（此时则为邮箱名称）
	m.SetAddressHeader("From", ep.fromEmail, "Mxbox verification code(Do not reply)")
	m.SetHeader("Subject", ep.Subject)
	m.SetBody("text/html", ep.Text)
	d := gomail.NewDialer(ep.serverHost, ep.serverPort, ep.fromEmail, ep.fromPassword)
	err := d.DialAndSend(m)
	if err != nil {
		Log.Error(RunFuncName(), zap.Any("emailParam", ep), zap.Any("err", err.Error()))
		return false
	}
	Log.Info(RunFuncName(), zap.Any("ok", ep))
	return true
}

// isEmailValid checks if the email provided passes the required structure// and length test. It also checks the domain has a valid MX record.
func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	return err == nil && len(mx) != 0
}

func SendEmailVerificationCode(vCode int32, mail string) bool {
	text := fmt.Sprintf("<h3>Verification code:%d, valid in 2 minutes, Do not give this code to anyone!!!</h3>", vCode)
	myEmail := &EmailParam{
		serverHost:   "smtp.*.com",
		serverPort:   465,
		fromEmail:    "user@***.com",
		fromPassword: "xxxx",
		Toers:        []string{mail},
		CCers:        []string{},
		Subject:      "verification code",
		Text:         text,
	}
	return SendEmail(myEmail)
}
