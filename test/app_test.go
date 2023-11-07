package test

import (
	"buding-job/common/utils"
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/jordan-wright/email"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net/smtp"
	"runtime"
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
	n := parse.NextN(time.Now(), 10)
	for _, v := range n {
		fmt.Println(v)
	}
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
	//orm.DB.Model(&do.JobLockDo{}).
	//	Where("expire_time < ?", time.Now()).
	//	Delete(&do.JobLockDo{})
}

func TestSelect(t *testing.T) {
	//var jobLogs []bo.JobTimeoutBo
	//orm.DB.Raw(constant.TimeoutJob).Scan(&jobLogs)
	//for _, value := range jobLogs {
	//	fmt.Println(value.Email)
	//	fmt.Println(value.JobLogDo)
	//}
}
func TestEmaiil(t *testing.T) {
	e := email.NewEmail()
	e.From = "1728717074@qq.com"
	e.To = []string{"319087181@qq.com"}
	e.Subject = "Awesome web"
	e.Text = []byte("Text Body is, of course, supported!")

	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "1728717074@qq.com",
		"xalbdtmrqtlzccei", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}
func TestOs(t *testing.T) {
	fmt.Println("os:%S\n", runtime.GOOS)
	fmt.Printf("CPU数量: %d\n", runtime.NumCPU())
	memory, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("无法获取内存信息: %v\n", err)
		return
	}

	fmt.Printf("总内存: %v GB\n", memory.Total/1024/1024/1024)
	fmt.Printf("可用内存: %v GB\n", memory.Available/1024/1024/1024)
	fmt.Printf("已用内存: %v GB\n", memory.Used/1024/1024/1024)
	fmt.Printf("内存使用率: %f%%\n", memory.UsedPercent)
	println(runtime.Version())
	println(runtime.GOROOT())
	println(runtime.NumGoroutine())
	var arr []byte
	runtime.Stack(arr, true)
	for _, v := range arr {
		fmt.Println(v)

	}
}
