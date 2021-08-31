<div align="center">

# golang-todoapp
  
simple to-do app API using beego

</div>

---

## API Diagrams

```
MODULES

User
------
/login                                                       # POST:   Login
/api/user/                                                   # POST:   Register
/api/user/                                                   # GET:    Get All Users
/api/user/<id>                                               # GET:    Get User by id
/api/user/<id>                                               # PUT:    Update User by id
/api/user/<id>                                               # DELETE: Delete User by id
/api/user/checkauth                                          # Any:    Check Authentication

Todo
------
/api/todo                                                    # POST:   Create A New Todo
/api/todo                                                    # GET:    Get All Todos
/api/todo/<id>                                               # GET:    Get Todo by id
/api/todo/<id>                                               # PUT:    Update Todo by id
/api/todo/<id>                                               # DELETE: Delete Todo by id
```

