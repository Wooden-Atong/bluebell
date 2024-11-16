package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota //🌟只要这里第一行写了ResCode类型，这一组常量都是这个类型；并且iota初始值为0，往下依次➕1
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidToken
	CodeExpiredAToken
	CodeExpiredRToken
	CodeNotExpiredRToken
	CodeNeedLogin
)

// 🌟这里定义的时候建议一组都加上一个Code，这样在别的文件调用的时候打Code，编译器就会提示，这样很方便
var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken:  "无效token",
	CodeExpiredAToken: "access token过期了",
	CodeExpiredRToken: "refresh token过期了",
	CodeNotExpiredRToken: "refresh token没过期, 返回了刷新的access token",
	CodeNeedLogin:     "需要登陆",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok { //相当于如果状态码不在codeMsgMap当中，那么直接定义为CodeServerBusy类型状态码
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
