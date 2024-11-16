package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

//🌟在controller层出错了，是日志记录错误，并直接返回响应（c.JSON(http.StatusOK,gin.H{"msg":"有误"}），
//🌟而在下面的层（比如logic层、dao层、models层），报错了一般是return err，将错误返回上一层。

func SignUpHandler(c *gin.Context) {
	//1.获取参数 和 参数校验
	// json格式校验
	p := new(models.ParamSignUp) //🌟var p models.ParamSignUp，这倒也是一种写法，但这样写1）只是分配内存空间，没有初始化；2）是值类型，而new一个返回的是指针类型
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors) //🌟这里用的是类型断言，判断错误是不是ValidationErrors
		if ok {
			ResponseErrorWithMsg(c,CodeInvalidParam,removeTopStruct(errs.Translate(trans)))//把错误翻译了再返回
			return
		}
		ResponseError(c,CodeInvalidParam)
		return
	}

	// 手动业务规则校验
	// if len(p.Username) == 0 || len(p.Password)==0||len(p.RePassword)==0||p.Password!=p.RePassword{
	// 	zap.L().Error("SignUp with invalid param")
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"msg":"请求参数有误",
	// 	})
	// }

	//2.业务处理 （放在logic层）
	if err := logic.SignUp(p); err != nil {
		if errors.Is(err, mysql.ErrorUserExist){
			ResponseError(c,CodeUserExist)
			return
		}
		ResponseError(c,CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c,nil)
}

func LoginHandler(c *gin.Context) {
	p := new(models.ParamLogin)
	//获取请求参数及参数校验
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			ResponseErrorWithMsg(c,CodeInvalidParam,removeTopStruct(errs.Translate(trans)))//把错误翻译了再返回
			return
		}
		ResponseError(c,CodeInvalidParam)
		return
	}
	//业务逻辑校验
	user,err := logic.Login(p)//🌟这一句必须单独拿出来写，如果在下面if后面跟着的话，user就属于if分支局部变量，后面返回响应的时候无法使用
	if err!=nil{
		zap.L().Error("Login.login() failed", zap.String("username",p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist){//可以用errors.Is()判断err是否是xx类型
			ResponseError(c,CodeUserNotExist)
			return
		}
		ResponseError(c,CodeServerBusy)
		return
	}

	//返回响应
	ResponseSuccess(c,gin.H{
		//🌟还是传给前端id失真的问题，但这里没在结构体定义的时候在tag-json中加string，
		//🌟这是因为，我们user结构体压根没有用到json序列化，所以干脆直接在这里改
		"user_id":fmt.Sprintf("%d",user.UserID),
		"user_name":user.Username,
		"a_token":user.AToken,
		"r_token":user.RToken,
	})
}
