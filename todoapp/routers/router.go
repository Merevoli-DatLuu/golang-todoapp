package routers

import (
	"todoapp/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.UserController{}, "post:Login")
	beego.Router("/api/user", &controllers.UserController{}, "get:GetAll;post:Post")
	beego.Router("/api/user/:id", &controllers.UserController{}, "get:GetOne;put:Put;delete:Delete")
	beego.Router("/api/user/checkauth", &controllers.UserController{}, "*:Auth")

	beego.Router("/api/todo", &controllers.TodoController{}, "get:GetAll;post:Post")
	beego.Router("/api/todo/:id", &controllers.TodoController{}, "get:GetOne;put:Put;delete:Delete")
}
