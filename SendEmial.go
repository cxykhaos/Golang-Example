package utils

import (
	"log"
	"net/smtp"

	"github.com/spf13/viper"
)

var addr string
var username string
var host string
var password string

func InitEmail() {
	addr = viper.GetString("website.addr")
	username = viper.GetString("website.username")
	host = viper.GetString("website.host")
	password = viper.GetString("website.password")
}

func Sentemail(To string, Subject string, content string) {
	// 设置接收方的邮箱
	To1 := []string{To}
	auth := smtp.PlainAuth("", username, password, host)
	AllC := []byte("To: " + To + "\r\nFrom: " + username + "\r\nSubject: " + Subject + "\r\n\r\n" + content)
	err := smtp.SendMail(addr, auth, username, To1, AllC)
	if err != nil {
		log.Fatal(err)
	}
}
