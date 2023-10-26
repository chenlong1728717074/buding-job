package orm

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB *gorm.DB
)

type dbConfig struct {
	name    string
	setting string
}

func init() {

	var err error
	userName := "root"
	passWorld := "990927"
	addr := "127.0.0.1"
	port := "3306"
	table := "xll-job"
	url := "charset=utf8mb4&parseTime=True&loc=Local"
	setting := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", userName, passWorld, addr, port, table, url)
	fmt.Println(setting)
	source, err := dataSource(dbConfig{
		name:    "mysql",
		setting: setting,
	})
	if err != nil {
		panic(source)
	}
	DB, err = gorm.Open(source,
		&gorm.Config{Logger: getLogger(),
			SkipDefaultTransaction: true,
		})
	if err != nil {
		panic("连接失败:" + err.Error())
	}
	//DB = DB.Scopes(WithDeletedFalse())
}
func dataSource(config dbConfig) (gorm.Dialector, error) {
	var dataBase gorm.Dialector
	switch config.name {
	case "mysql":
		dataBase = mysql.Open(config.setting)
	default:
		return nil, errors.New("no suitable data source matching")
	}
	return dataBase, nil
}
func getLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
}
