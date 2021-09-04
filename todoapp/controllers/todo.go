package controllers

import (
	"encoding/json"
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

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if errorMessage := utils.CheckNewTodo(v.Name, v.Description, v.UserId, v.Status); errorMessage != "ok" {
			c.RetResponse(&Response{FailStatus, 403, errorMessage, ""}, 403)
			return
		}

		if !models.CheckUserId(v.UserId) {
			c.RetResponse(&Response{FailStatus, 403, "UserId is not exist", ""}, 403)
			return
		}

		if todoId, err := models.AddTodo(&v); err == nil {
			var returnData = &CreateObjectData{int(todoId)}
			c.RetResponse(&Response{SuccessStatus, 0, "Create todo successfully", returnData}, 201)
			return
		} else {
			c.RetResponse(&Response{FailStatus, 403, err.Error(), ""}, 403)
			return
		}
	} else {
		c.RetResponse(&Response{FailStatus, 403, err.Error(), ""}, 403)
		return
	}
}

func (c *TodoController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetTodoById(id)

	if err != nil {
		c.RetResponse(&Response{FailStatus, 404, err.Error(), ""}, 404)
		return
	} else {
		c.RetResponse(&Response{SuccessStatus, 200, "Get todo successfully", v}, 200)
		return
	}
}

func (c *TodoController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int = 20
	var offset int

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
				c.RetResponse(&Response{FailStatus, 400, "Error: invalid query key/value pair", ""}, 400)
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllTodos(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.RetResponse(&Response{FailStatus, 400, err.Error(), ""}, 400)
		return
	} else {
		c.RetResponse(&Response{SuccessStatus, 0, "Get all todos successfully", l}, 200)
		return
	}
}

func (c *TodoController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Todo{Id: id}

	u, err := models.GetTodoById(id)
	if err != nil {
		c.RetResponse(&Response{FailStatus, 404, err.Error(), ""}, 404)
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateTodoById(&v, u); err == nil {
			c.RetResponse(&Response{SuccessStatus, 0, "Update todo successfully", ""}, 201)
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

func (c *TodoController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteTodo(id); err == nil {
		c.RetResponse(&Response{SuccessStatus, 0, "Delete todo successfully", ""}, 200)
		return
	} else {
		c.RetResponse(&Response{FailStatus, 400, err.Error(), ""}, 400)
		return
	}
}
