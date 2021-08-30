package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"todoapp/models"
	"todoapp/utils"
)

type TodoController struct {
	BaseController
}

func (c *TodoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

func (c *TodoController) Post() {
	var v models.Todo

	token := c.Ctx.Input.Header("token")
	et := utils.EasyToken{}
	valido, err := et.ValidateToken(token)
	if !valido {
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Data["json"] = Response{401, 401, fmt.Sprintf("%s", err), ""}
		c.ServeJSON()
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if errorMessage := utils.CheckNewTodo(v.Name, v.Description, v.UserId, v.Status); errorMessage != "ok" {
			c.Ctx.ResponseWriter.WriteHeader(403)
			c.Data["json"] = Response{403, 403, errorMessage, ""}
			c.ServeJSON()
			return
		}
		if !models.CheckUserId(v.UserId) {
			c.Ctx.ResponseWriter.WriteHeader(403)
			c.Data["json"] = Response{403, 403, "UserId is not exist", ""}
			c.ServeJSON()
			return
		}

		if todoId, err := models.AddTodo(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			var returnData = &CreateObjectData{int(todoId)}
			c.Data["json"] = &Response{0, 0, "ok", returnData}
		} else {
			c.Ctx.ResponseWriter.WriteHeader(403)
			c.Data["json"] = err.Error()
			c.ServeJSON()
			return
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *TodoController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetTodoById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

func (c *TodoController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int = 20
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

	found, _ := models.GetUserByToken(token)
	if !found {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = &Response{401, 401, "user is not exist", ""}
		c.ServeJSON()
		return
	}

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
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

	l, err := models.GetAllTodos(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		var returnData = &GetTodoData{l}
		c.Data["json"] = &Response{0, 0, "ok", returnData}
	}
	c.ServeJSON()
}

func (c *TodoController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Todo{Id: id}

	u, err := models.GetTodoById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateTodoById(&v, u); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *TodoController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteTodo(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
