package test

import (
	"buding-job/common/constant"
	"buding-job/common/utils"
	"buding-job/orm"
	"buding-job/orm/bo"
	"buding-job/orm/do"
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"testing"
	"time"
)

func TestCobra(t *testing.T) {
	//_ = do.JobInfoDo{}
	//_ = do.JobManagementDo{}
}

func TestTime(t *testing.T) {
	now := time.Now().Add(-time.Second * 90)
	fmt.Println(now)
	fmt.Println(now.Before(time.Now()))
}
func TestCron(t *testing.T) {
	//parse, _ := cronexpr.Parse("0/5 * * * * ? *")
	fmt.Println(time.Now())
	parse, _ := cronexpr.Parse("0 0/5 * * * ? *")
	fmt.Println(parse.Next(time.Now()))
}
func TestArr(t *testing.T) {
	fmt.Println(len(getSince()))

}
func getSince() []int {
	return nil
}
func TestRand(t *testing.T) {
	println(utils.RandI64(1))
}
func TestDelete(t *testing.T) {
	orm.DB.Model(&do.JobLockDo{}).
		Where("expire_time < ?", time.Now()).
		Delete(&do.JobLockDo{})
}

func TestSelect(t *testing.T) {
	var jobLogs []bo.JobTimeoutBo
	orm.DB.Raw(constant.TimeoutJob).Scan(&jobLogs)
	for _, value := range jobLogs {
		fmt.Println(value.Email)
		fmt.Println(value.JobLogDo)
	}
}
func TestEmaiil(t *testing.T) {
	e := email.NewEmail()
	e.From = "dj <xxx@126.com>"
	e.To = []string{"935653229@qq.com"}
	e.Cc = []string{"test1@126.com", "test2@126.com"}
	e.Bcc = []string{"secret@126.com"}
	e.Subject = "Awesome web"
	e.Text = []byte("Text Body is, of course, supported!")
	err := e.Send("smtp.126.com:25", smtp.PlainAuth("", "xxx@126.com", "yyy", "smtp.126.com"))
	if err != nil {
		log.Fatal(err)
	}
}
