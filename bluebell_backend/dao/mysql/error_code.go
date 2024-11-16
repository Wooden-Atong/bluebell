package mysql

import "errors"

var(//🌟用errors.New()一般是自己加的业务逻辑判断的错误，而直接获得的error类型的err是程序报错
	ErrorUserExist = errors.New("用户已存在！")
	ErrorUserNotExist = errors.New("用户不存在！")
	ErrorInvalidPassword = errors.New("用户名或密码错误！")
	ErrorInvalidID = errors.New("无效的ID")
)