package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"todoapp/models"
	"todoapp/utils"

	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
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
			c.Ctx.ResponseWriter.WriteHeader(403)
			c.Data["json"] = Response{403, 403, errorMessage, ""}
			c.ServeJSON()
			return
		}
		if models.CheckUserName(v.Username) {
			c.Ctx.ResponseWriter.WriteHeader(403)
			c.Data["json"] = Response{403, 403, "username is already registered", ""}
			c.ServeJSON()
			return
		}
		if models.CheckEmail(v.Email) {
			c.Ctx.ResponseWriter.WriteHeader(403)
			c.Data["json"] = Response{403, 403, "email is already registered", ""}
			c.ServeJSON()
			return
		}

		if user, err := models.AddUser(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			var returnData = &UserSuccessLoginData{user.Token, user.Username}
			c.Data["json"] = &Response{0, 0, "ok", returnData}
		} else {
			c.Data["json"] = &Response{1, 1, "user registration failed", err.Error()}
		}
	} else {
		c.Data["json"] = &Response{1, 1, "user registration failed", err.Error()}
	}
	c.ServeJSON()
}

func (c *UserController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int = 10
	var offset int

	token := c.Ctx.Input.Header("token")
	et := utils.EasyToken{}
	validation, err := et.ValidateToken(token)
	if !validation {
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Data["json"] = Response{401, 401, fmt.Sprintf("%s", err), ""}
		c.ServeJSON()
		return
	}

	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	} else {
		fields = strings.Split("Username,Gender,Age,Email,Token", ",")
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
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllUser(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

func (c *UserController) GetOne() {
	token := c.Ctx.Input.Header("token")
	idStr := c.Ctx.Input.Param(":id")
	et := utils.EasyToken{}
	valido, err := et.ValidateToken(token)
	if !valido {
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Data["json"] = Response{401, 401, fmt.Sprintf("%s", err), ""}
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(idStr)
	v, err := models.GetUserById(id)
	if v == nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()

}

func (c *UserController) Put() {
	token := c.Ctx.Input.Header("token")
	et := utils.EasyToken{}
	valido, err := et.ValidateToken(token)
	if !valido {
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Data["json"] = Response{401, 401, fmt.Sprintf("%s", err), ""}
		c.ServeJSON()
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.User{Id: id}

	u, err := models.GetUserById(id)
	if u == nil {
		c.Data["json"] = err.Error()
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateUserById(&v, u); err == nil {
			c.Data["json"] = successReturn
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *UserController) Delete() {
	token := c.Ctx.Input.Header("token")
	et := utils.EasyToken{}
	valido, err := et.ValidateToken(token)
	if !valido {
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Data["json"] = Response{401, 401, fmt.Sprintf("%s", err), ""}
		c.ServeJSON()
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteUser(id); err == nil {
		c.Data["json"] = successReturn
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *UserController) Login() {
	var reqData struct {
		Username string `valid:"Required"`
		Password string `valid:"Required"`
	}
	var token string

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqData); err == nil {
		if errorMessage := utils.CheckUsernamePassword(reqData.Username, reqData.Password); errorMessage != "ok" {
			c.Ctx.ResponseWriter.WriteHeader(403)
			c.Data["json"] = Response{403, 403, errorMessage, ""}
			c.ServeJSON()
			return
		}
		if ok, user := models.Login(reqData.Username, reqData.Password); ok {
			et := utils.EasyToken{}
			validation, err := et.ValidateToken(user.Token)
			if !validation {
				et = utils.EasyToken{
					Username: user.Username,
					Uid:      int64(user.Id),
					Expires:  time.Now().Unix() + 2*3600,
				}
				token, err = et.GetToken()
				if token == "" || err != nil {
					c.Data["json"] = errUserToken
					c.ServeJSON()
					return
				} else {
					models.UpdateUserToken(user, token)
				}
			} else {
				token = user.Token
			}

			var returnData = &UserSuccessLoginData{token, user.Username}
			c.Data["json"] = &Response{0, 0, "ok", returnData}
		} else {
			c.Data["json"] = &errNoUserOrPass
		}
	} else {
		c.Data["json"] = &errNoUserOrPass
	}
	c.ServeJSON()
}

func (c *UserController) Auth() {
	et := utils.EasyToken{}
	token := strings.TrimSpace(c.Ctx.Request.Header.Get("Authorization"))
	validation, err := et.ValidateToken(token)
	if !validation {
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Data["json"] = Response{401, 401, fmt.Sprintf("%s", err), ""}
		c.ServeJSON()
		return
	}

	c.Data["json"] = Response{0, 0, "is login", ""}
	c.ServeJSON()
}

func (u *UserController) Logout() {
	u.Data["json"] = successReturn
	u.ServeJSON()
}
