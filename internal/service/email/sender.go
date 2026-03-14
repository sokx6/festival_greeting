package email

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/smtp"
)

type EmailSender struct {
	From     string
	To       string
	Subject  string
	Port     int
	Host     string
	Username string
	Password string
	Content  string
}

// sendMailWithTLS 使用 TLS 加密连接发送邮件（用于端口 465）
func sendMailWithTLS(s *EmailSender, auth smtp.Auth, message string) error {
	// 创建 TLS 配置
	tlsConfig := &tls.Config{
		ServerName: s.Host,
	}

	// 建立 TLS 连接
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port), tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS connection failed: %v", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %v", err)
	}
	defer client.Close()

	// 认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %v", err)
	}

	// 设置发件人
	if err = client.Mail(s.From); err != nil {
		return fmt.Errorf("setting sender failed: %v", err)
	}

	// 设置收件人
	addr := s.To
	if err = client.Rcpt(addr); err != nil {
		return fmt.Errorf("setting recipient failed: %v", err)
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data command failed: %v", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("writing message failed: %v", err)
	} else {
		fmt.Println("邮件内容写入成功")
		fmt.Printf("邮件主题: %s\n", s.Subject)
		fmt.Printf("邮件收件人: %s\n", s.To)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("closing data writer failed: %v", err)
	}

	return client.Quit()
}

func (s *EmailSender) Send() error {

	var err error
	// 设置邮件头部
	header := make(map[string]string)
	header["From"] = s.From
	header["To"] = s.To
	header["Subject"] = mime.QEncoding.Encode("UTF-8", s.Subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + s.Content

	// 认证信息
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	// 根据端口选择发送方式
	if s.Port == 465 {
		// 使用 SSL/TLS 加密连接（端口 465）
		err = sendMailWithTLS(s, auth, message)
	} else {
		// 使用 STARTTLS（端口 587）或明文（端口 25，不推荐）
		err = smtp.SendMail(fmt.Sprintf("%s:%d", s.Host, s.Port), auth, s.From, []string{s.To}, []byte(message))
	}

	if err != nil {
		return err
	}
	return nil
}

func NewEmailSender(from, to, subject, host string, port int, username, password, content string) *EmailSender {
	return &EmailSender{
		From:     from,
		To:       to,
		Subject:  subject,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Content:  content,
	}
}
