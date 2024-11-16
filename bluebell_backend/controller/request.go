package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

func getCurrentUserID(c *gin.Context)(userID int64, err error){
	uid, ok := c.Get(ContextUserIDKey)
	if !ok{
		//🌟由于在函数返回参数已经定义了err，所以这里写返回ErrorUserNotLogin，实际上是赋值给err然后返回
		//🌟这里userID没有改动过，但定义的时候就已经初始化，所以返回的是默认值
		return userID,ErrorUserNotLogin
	}
	userID, ok = uid.(int64)
	if !ok{
		return userID,ErrorUserNotLogin
	}
	return
}

func getPageInfo(c *gin.Context)(int64,int64){
	var (
		pageNum int64
		pageSize int64
		err error
	)

	//获取分页参数
	pageNumStr := c.Query("page") //🌟获取querystring参数
	pageSizeStr := c.Query("size")

	pageNum,err = strconv.ParseInt(pageNumStr,10,64)
	if err!=nil{
		pageNum = 1
	}
	pageSize,err = strconv.ParseInt(pageSizeStr,10,64)
	if err!=nil{
		pageSize=10
	}
	return pageNum,pageSize
}