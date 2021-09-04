package routers

import (
	"strings"
	"todoapp/controllers"
	"todoapp/models"
	"todoapp/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
			beego.NSRouter("/register", &controllers.UserController{}, "post:Post"),
			beego.NSNamespace("/user",
				beego.NSBefore(auth),
				beego.NSRouter("/", &controllers.UserController{}, "get:GetAll"),
				beego.NSRouter("/:id", &controllers.UserController{}, "get:GetOne;put:Put;delete:Delete"),
				beego.NSRouter("/check_auth", &controllers.UserController{}, "*:Auth"),
			),
			beego.NSNamespace("/todo",
				beego.NSBefore(auth),
				beego.NSRouter("/", &controllers.TodoController{}, "get:GetAll;post:Post"),
				beego.NSRouter("/:id", &controllers.TodoController{}, "get:GetOne;put:Put;delete:Delete"),
			),
		),
	)

	beego.AddNamespace(ns)
}

var auth = func(ctx *context.Context) {
	token_raw := ctx.Input.Header("Authorization")
	token_fields := strings.Split(token_raw, " ")

	if len(token_fields) == 2 && token_fields[0] == "Bearer" {
		token := token_fields[1]
		println(token)
		et := utils.EasyToken{}
		validation, err := et.ValidateToken(token)

		if !validation {
			controllers.RetUnauthorizedResponse(ctx, err.Error())
			return
		}

		found, _ := models.GetUserByToken(token)
		if !found {
			controllers.RetUnauthorizedResponse(ctx, "user is not exist")
			return
		}

	} else {
		controllers.RetUnauthorizedResponse(ctx, "Wrong format token")
	}
}
