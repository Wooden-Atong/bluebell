package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code":100, //状态码
	"msg":,xx,//提示信息
	"data":{},//数据
}
*/

//🌟这里msg和data都定义为空接口，方便传输数据.其实也可以不用定义，返回响应直接用gin.H{map[string]any}是一样的
type ResponseData struct{
	Code ResCode	`json:"code"`
	Msg interface{} `json:"msg"`
	Data interface{}`json:"data,omitempty"`//🌟omitempty就是当这个字段为空的时候，就不会展示它
}

func ResponseError(c *gin.Context, code ResCode){
	rd := &ResponseData{
		Code:code,
		Msg:code.Msg(),
		Data:nil,
	}
	c.JSON(http.StatusOK,rd)
}


func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}){
	rd := &ResponseData{
		Code:code,
		Msg:msg,//这样可以传进来的时候自行制定msg
		Data:nil,
	}
	c.JSON(http.StatusOK,rd)
}

//🌟这里本来我的想法是还是先传进来code，然后把它丢在code里面去判断一下，但是视频这里直接传data（因为已经成功了）。我的思路是有问题的，我以为我们获取一个code然后去找问题，实际上是我们知道问题用code方便表述罢了
func ResponseSuccess(c *gin.Context, data interface{}){
	rd := &ResponseData{
		Code:CodeSuccess,
		Msg:CodeSuccess.Msg(),
		Data:data,
	}
	c.JSON(http.StatusOK,rd)
}