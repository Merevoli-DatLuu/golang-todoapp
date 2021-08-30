package main

import (
	_ "todoapp/routers"
	"todoapp/utils"

	"github.com/astaxie/beego"
)

func main() {
	utils.InitSql()
	beego.Run()
}
