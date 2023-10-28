package orm

import (
	"errors"
	"fmt"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
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
	setting := buildDbConfig()
	source, err := dataSource(setting)
	if err != nil {
		panic(err)
	}
	DB, err = gorm.Open(source,
		&gorm.Config{Logger: getLogger(),
			SkipDefaultTransaction: true,
		})
	if err != nil {
		panic("连接失败:" + err.Error())
	}
}

func buildDbConfig() dbConfig {
	userName := "root"
	passWorld := "990927"
	addr := "127.0.0.1"
	port := "3306"
	table := "buding-job"
	url := "charset=utf8mb4&parseTime=True&loc=Local"
	setting := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", userName, passWorld, addr, port, table, url)
	return dbConfig{
		name:    "mysql",
		setting: setting,
	}
}

func dataSource(config dbConfig) (gorm.Dialector, error) {
	var dataBase gorm.Dialector
	switch config.name {
	case "mysql":
		dataBase = mysql.Open(config.setting)
	case "postgres":
		dataBase = postgres.Open(config.setting)
	case "sqlite":
		dataBase = sqlite.Open("gorm.db")
	case "sqlserver":
		dataBase = sqlserver.Open("gorm.db")
	case "clickhouse":
		dataBase = clickhouse.Open("gorm.db")
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
