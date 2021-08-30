package controllers

import (
	"github.com/astaxie/beego"
)

const (
	ErrInputData    = "input data error"
	ErrDatabase     = "database error"
	ErrDupUser      = "duplication user erroe"
	ErrNoUser       = "no user error"
	ErrPass         = "password error"
	ErrNoUserPass   = "no user pass error"
	ErrNoUserChange = "no user change error"
	ErrInvalidUser  = "invalid user error"
	ErrOpenFile     = "open file error"
	ErrWriteFile    = "write file error"
	ErrSystem       = "system error"
)

var (
	successReturn   = &Response{200, 0, "ok", "ok"}
	err404          = &Response{404, 404, "Webpage not found", "Webpage not found"}
	errInputData    = &Response{400, 10001, "Data input error", "Client parameter error"}
	errDatabase     = &Response{500, 10002, "Server Error", "Database operation error"}
	errUserToken    = &Response{500, 10002, "Server Error", "Token operation error"}
	errDupUser      = &Response{400, 10003, "User information already exists", "Duplicate database records"}
	errNoUser       = &Response{400, 10004, "User information does not exist", "Database record does not exist"}
	errPass         = &Response{400, 10005, "User information does not exist or the password is incorrect", "Incorrect password"}
	errNoUserOrPass = &Response{400, 10006, "The user does not exist or the password is incorrect", "The database record does not exist or the password is incorrect"}
	errNoUserChange = &Response{400, 10007, "The user does not exist or the data has not changed", "The database record does not exist or the data has not changed"}
	errInvalidUser  = &Response{400, 10008, "User information is incorrect", "Session information is incorrect"}
	errOpenFile     = &Response{500, 10009, "Server Error", "Error opening file"}
	errWriteFile    = &Response{500, 10010, "Server Error", "Error writing file"}
	errSystem       = &Response{500, 10011, "Server Error", "Operating system error"}
	errExpired      = &Response{400, 10012, "Login has expired", "Verify that the token expires"}
	errPermission   = &Response{400, 10013, "Permission denied", "No operation authority"}
)

type UserSuccessLoginData struct {
	AccessToken string `json:"access_token"`
	UserName    string `json:"user_name"`
}

type CreateObjectData struct {
	Id int `json:"id"`
}

type GetTodoData struct {
	Todo interface{} `json:"todo"`
}

func (base *BaseController) RetError(e *Response) {
	if mode := beego.AppConfig.String("runmode"); mode == "prod" {
		e.Data = ""
	}

	base.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	base.Ctx.ResponseWriter.WriteHeader(e.Status)
	base.Data["json"] = e
	base.ServeJSON()
	base.StopRun()
}

var sqlOp = map[string]string{
	"eq": "=",
	"ne": "<>",
	"gt": ">",
	"ge": ">=",
	"lt": "<",
	"le": "<=",
}
