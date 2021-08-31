package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Todo struct {
	Id          int    `json:"id" orm:"column(id);pk;unique;auto;int(11)"`
	Name        string `json:"name" orm:"column(name);size(100)"`
	Description string `json:"description, omitempty" orm:"column(description);size(255)"`
	UserId      int    `json:"user_id" orm:"column(user_id);size(11)"`
	Status      string `json:"status" orm:"column(status);size(27)"`
	CreatedDate int64  `json:"created_date, omitempty" orm:"column(created_at);size(11)"`
	UpdateDate  int64  `json:"updated_date, omitempty" orm:"column(updated_date);size(11)"`
}

func init() {
	orm.RegisterModel(new(Todo))
}

func (t *Todo) TableName() string {
	return TableName("todo")
}

func Todos() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Todo))
}

func AddTodo(m *Todo) (id int64, err error) {
	o := orm.NewOrm()

	CreatedAt := time.Now().UTC().Unix()
	UpdatedAt := CreatedAt

	todo := Todo{
		Name:        m.Name,
		Description: m.Description,
		UserId:      m.UserId,
		Status:      m.Status,
		CreatedDate: CreatedAt,
		UpdateDate:  UpdatedAt,
	}

	fmt.Println(todo.Id)

	id, err = o.Insert(&todo)
	if err == nil {
		return id, err
	}

	return 0, err
}

func GetTodoById(id int) (v *Todo, err error) {
	o := orm.NewOrm()
	v = &Todo{Id: id}
	if err = o.QueryTable(new(Todo)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAllTodos(query map[string]string, fields []string, sortby []string, order []string,
	offset int, limit int) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Todo))
	for k, v := range query {
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Todo
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

func UpdateTodoById(m *Todo, u *Todo) (err error) {
	o := orm.NewOrm()
	v := Todo{Id: m.Id}

	if m.Name == "" {
		m.Name = u.Name
	}
	if m.UserId == 0 {
		m.UserId = u.UserId
	}
	if m.Description == "" {
		m.Description = u.Description
	}
	if m.Status == "" {
		m.Status = u.Status
	}
	if m.CreatedDate == 0 {
		m.CreatedDate = u.CreatedDate
	}

	m.UpdateDate = time.Now().UTC().Unix()

	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteTodo(id int) (err error) {
	o := orm.NewOrm()
	v := Todo{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Todo{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
