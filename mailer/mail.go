package mailer

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xhit/go-simple-mail/v2"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

//var smtpServer *mail.SMTPServer
var channels map[string]*mail.SMTPServer

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b) + "/template"
)

func SetTemplatePath(path string) {
	basepath = path
}

func InitMailChannel(name string, config ChannelConfig) {
	if channels == nil {
		channels = make(map[string]*mail.SMTPServer)
	}

	if _, ok := channels[name]; ok {
		log.Errorf("Channel with name %s already initialized")
		return
	}

	smtpServer := mail.NewSMTPClient()

	smtpServer.Host = config.Host
	smtpServer.Port = config.Port
	smtpServer.Username = config.Username
	smtpServer.Password = config.Password

	switch config.Encrypt {
	case "ssl":
		smtpServer.Encryption = mail.EncryptionSSL
	case "tls":
		smtpServer.Encryption = mail.EncryptionTLS
	default:
		smtpServer.Encryption = mail.EncryptionNone
	}

	smtpServer.KeepAlive = false
	smtpServer.ConnectTimeout = 10 * time.Second
	smtpServer.SendTimeout = 10 * time.Second
	channels[name] = smtpServer
}

func SendMail(channelName string, email *mail.Email) error {
	channel, ok := channels[channelName]
	if !ok {
		return fmt.Errorf("Channel with name %s not exists ", channelName)
	}

	smtpClient, err := channel.Connect()
	if err != nil {
		log.WithError(err).Errorf("Fail connect to SMTP server")
		return err
	}

	defer func() {
		smtpClient.Close()
	}()

	email.SetFrom(channel.Username)

	return email.Send(smtpClient)
}

func CreateEmail(templateName string, data interface{}) (*mail.Email, error) {
	templateFileName := fmt.Sprintf("%s/%s.html", basepath, templateName)
	if _, err := os.Stat(templateFileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("Template %s not found in %s ", templateName, templateFileName)
	}

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return nil, err
	}

	email := mail.NewMSG()
	email.SetBody(mail.TextHTML, buffer.String())

	return email, nil
}
