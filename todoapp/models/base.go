package models

import (
	"github.com/astaxie/beego/config"
	_ "github.com/go-sql-driver/mysql"
)

func TableName(str string) string {
	appConf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		panic(err)
	}
	return appConf.String("database::") + str
}
