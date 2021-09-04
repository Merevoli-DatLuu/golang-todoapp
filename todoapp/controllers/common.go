package controllers

import (
	"fmt"

	"github.com/astaxie/beego/context"
)

const (
	SuccessStatus   = "success"
	FailStatus      = "fail"
	ErrorStatus     = "error"
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
	successReturn   = &Response{SuccessStatus, 0, "success", "success"}
	err404          = &Response{FailStatus, 404, "Webpage not found", "Webpage not found"}
	errInputData    = &Response{FailStatus, 10001, "Data input error", "Client parameter error"}
	errDupUser      = &Response{FailStatus, 10003, "User information already exists", "Duplicate database records"}
	errNoUser       = &Response{FailStatus, 10004, "User information does not exist", "Database record does not exist"}
	errPass         = &Response{FailStatus, 10005, "User information does not exist or the password is incorrect", "Incorrect password"}
	errNoUserOrPass = &Response{FailStatus, 10006, "The user does not exist or the password is incorrect", "The database record does not exist or the password is incorrect"}
	errNoUserChange = &Response{FailStatus, 10007, "The user does not exist or the data has not changed", "The database record does not exist or the data has not changed"}
	errInvalidUser  = &Response{FailStatus, 10008, "User information is incorrect", "Session information is incorrect"}
	errExpired      = &Response{FailStatus, 10012, "Login has expired", "Verify that the token expires"}
	errPermission   = &Response{FailStatus, 10013, "Permission denied", "No operation authority"}
	errUserToken    = &Response{ErrorStatus, 10002, "Server Error", "Token operation error"}
	errDatabase     = &Response{ErrorStatus, 10002, "Server Error", "Database operation error"}
	errOpenFile     = &Response{ErrorStatus, 10009, "Server Error", "Error opening file"}
	errWriteFile    = &Response{ErrorStatus, 10010, "Server Error", "Error writing file"}
	errSystem       = &Response{ErrorStatus, 10011, "Server Error", "Operating system error"}
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

func (base *BaseController) RetResponse(e *Response, status int) {
	base.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	base.Ctx.ResponseWriter.WriteHeader(status)
	base.Data["json"] = e
	base.ServeJSON()
	// base.StopRun()
}

func RetUnauthorizedResponse(ctx *context.Context, err string) {
	response_data := &Response{FailStatus, 401, fmt.Sprintf("%s", err), ""}
	ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	ctx.ResponseWriter.WriteHeader(401)
	ctx.Output.JSON(response_data, true, true)
}

var sqlOp = map[string]string{
	"eq": "=",
	"ne": "<>",
	"gt": ">",
	"ge": ">=",
	"lt": "<",
	"le": "<=",
}
