package utils

import (
	"fmt"

	"github.com/astaxie/beego/validation"
)

func CheckUsernamePassword(username string, password string) (errorMessage string) {
	valid := validation.Validation{}
	valid.Required(username, "Username").Message("username is required")
	valid.Required(password, "Password").Message("password is required")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return fmt.Sprintf("%s", err.Message)
		}
	}
	return fmt.Sprintf("%s", "ok")
}

func CheckNewUserPost(Username string, Password string, Age int,
	Gender int, Email string) (errorMessage string) {
	valid := validation.Validation{}
	valid.Required(Username, "Username").Message("username is required")
	valid.AlphaNumeric(Username, "Username").Message("username must have digits or characters")
	valid.Required(Password, "Password").Message("password is required")
	valid.MinSize(Password, 6, "Password").Message("password cannot be less than 6 digits")
	valid.MaxSize(Password, 20, "Password").Message("password cannot be more than 20 digits")
	valid.Required(Age, "Age").Message("age is required")
	valid.Range(Age, 1, 100, "Age").Message("age must between 1 and 100 years old")
	valid.Range(Gender, 0, 2, "Gender").Message("incorrect gender")
	valid.Required(Email, "Email").Message("email is required")
	valid.Email(Email, "Email").Message("email format is incorrect")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return fmt.Sprintf("%s", err.Message)
		}
	}
	return fmt.Sprintf("%s", "ok")
}

func CheckNewTodo(name string, description string, userid int, status string) (errorMessage string) {
	valid := validation.Validation{}
	valid.Required(name, "Name").Message("用户ID必填")
	valid.Required(description, "Description").Message("设备名必填")
	valid.Required(userid, "User Id").Message("地址必填")
	valid.Required(status, "Status").Message("用户ID必填")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return fmt.Sprintf("%s", err.Message)
		}
	}
	return fmt.Sprintf("%s", "ok")
}
