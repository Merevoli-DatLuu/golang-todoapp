package controllers

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"todoapp/models"
	"todoapp/utils"
)

type UserController struct {
	BaseController
}

func (c *UserController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

func (c *UserController) Post() {
	var v models.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if errorMessage := utils.CheckNewUserPost(v.Username, v.Password,
			v.Age, v.Gender, v.Email); errorMessage != "ok" {
			c.RetResponse(&Response{FailStatus, 403, errorMessage, ""}, 403)
			return
		}

		if models.CheckUserName(v.Username) {
			c.RetResponse(&Response{FailStatus, 403, "username is already registered", ""}, 403)
			return
		}

		if models.CheckEmail(v.Email) {
			c.RetResponse(&Response{FailStatus, 403, "email is already registered", ""}, 403)
			return
		}

		if user, err := models.AddUser(&v); err == nil {
			var returnData = &UserSuccessLoginData{user.Token, user.Username}
			c.RetResponse(&Response{SuccessStatus, 0, "register successfully", returnData}, 201)
			return
		} else {
			c.RetResponse(&Response{FailStatus, 400, err.Error(), ""}, 400)
			return
		}
	} else {
		c.RetResponse(&Response{FailStatus, 400, err.Error(), ""}, 400)
		return
	}
}

func (c *UserController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int = 10
	var offset int

	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	} else {
		fields = strings.Split("Id,Username,Gender,Age,Email", ",")
	}
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	}
	if v, err := c.GetInt("offset"); err == nil {
		offset = v
	}
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.RetResponse(&Response{FailStatus, 400, "invalid query key/value pair", ""}, 400)
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllUser(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.RetResponse(&Response{FailStatus, 400, err.Error(), ""}, 400)
		return
	} else {
		c.RetResponse(&Response{SuccessStatus, 0, "Get all users successfully", l}, 200)
		return
	}
}

func (c *UserController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")

	id, _ := strconv.Atoi(idStr)
	v, err := models.GetUserById(id)
	if v == nil {
		c.RetResponse(&Response{FailStatus, 400, err.Error(), ""}, 400)
		return
	} else {
		c.RetResponse(&Response{SuccessStatus, 0, "Get user successfully", v}, 200)
		return
	}
}

func (c *UserController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.User{Id: id}

	u, err := models.GetUserById(id)
	if u == nil {
		c.RetResponse(&Response{FailStatus, 404, err.Error(), ""}, 404)
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateUserById(&v, u); err == nil {
			c.RetResponse(&Response{SuccessStatus, 0, "Update user successfully", ""}, 200)
			return
		} else {
			c.RetResponse(&Response{FailStatus, 400, err.Error(), ""}, 400)
			return
		}
	} else {
		c.RetResponse(&Response{FailStatus, 400, err.Error(), ""}, 400)
		return
	}
}

func (c *UserController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteUser(id); err == nil {
		c.RetResponse(&Response{SuccessStatus, 0, "Delete user successfully", ""}, 200)
		return
	} else {
		c.RetResponse(&Response{FailStatus, 400, err.Error(), ""}, 400)
		return
	}
}

func (c *UserController) Login() {
	var reqData struct {
		Username string `valid:"Required"`
		Password string `valid:"Required"`
	}
	var token string

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqData); err == nil {
		if errorMessage := utils.CheckUsernamePassword(reqData.Username, reqData.Password); errorMessage != "ok" {
			c.RetResponse(&Response{FailStatus, 403, errorMessage, ""}, 403)
			return
		}
		if ok, user := models.Login(reqData.Username, reqData.Password); ok {
			et := utils.EasyToken{}
			validation, _ := et.ValidateToken(user.Token)
			if !validation {
				et = utils.EasyToken{
					Username: user.Username,
					Uid:      int64(user.Id),
					Expires:  time.Now().Unix() + 2*3600,
				}
				token, err = et.GetToken()
				if token == "" || err != nil {
					c.RetResponse(errUserToken, 400)
					return
				} else {
					models.UpdateUserToken(user, token)
				}
			} else {
				token = user.Token
			}

			var returnData = &UserSuccessLoginData{token, user.Username}
			c.RetResponse(&Response{SuccessStatus, 0, "login successfully", returnData}, 200)
			return
		} else {
			c.RetResponse(errNoUserOrPass, 200)
			return
		}
	} else {
		c.RetResponse(errNoUserOrPass, 200)
		return
	}
}

func (c *UserController) Auth() {
	c.RetResponse(&Response{SuccessStatus, 0, "token valid", ""}, 200)
}

func (u *UserController) Logout() {
	u.Data["json"] = successReturn
	u.ServeJSON()
}
