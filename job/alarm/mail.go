package alarm

import (
	"buding-job/common/log"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"time"
)

func init() {
	Mail = NewMailAlarm()
}

const subject = "[来自BuDing-Job的的预警信息]"
const textTemp = "尊敬的%s先生/女士,布丁向您发来温馨的预警,您所负责的任务管理器:%s下的任务:%s于%s执行失败,请您登录BuDing-Job管理平台查看"

var Mail *MailAlarm

type MailAlarm struct {
	MailPool *email.Pool
}

func NewMailAlarm() *MailAlarm {
	pool, err := email.NewPool(
		"smtp.qq.com:587",
		5,
		smtp.PlainAuth("", "1728717074@qq.com", "xalbdtmrqtlzccei", "smtp.qq.com"),
	)

	if err != nil {
		log.GetLog().Fatal("failed to create pool:", err)
	}

	return &MailAlarm{
		pool,
	}
}
func NewAlarmEmail(author string, mailAddr string, manageName string, jobName string) *email.Email {
	e := email.NewEmail()
	e.From = "1728717074@qq.com"
	e.To = []string{mailAddr}
	e.Subject = subject
	text := fmt.Sprintf(textTemp, author, manageName, jobName, time.Now().Format("2006-01-02 15:04:05"))
	e.Text = []byte(text)
	return e
}

func (mail *MailAlarm) SendTxt(information *email.Email) bool {
	err := mail.MailPool.Send(information, 2*time.Second)
	if err != nil {
		return false
	}
	return true
}

func (mail *MailAlarm) CommonAlarm(author string, mailAddr string, manageName string, jobName string) bool {
	err := mail.MailPool.Send(NewAlarmEmail(author, mailAddr, manageName, jobName), 2*time.Second)
	if err != nil {
		return false
	}
	return true
}
