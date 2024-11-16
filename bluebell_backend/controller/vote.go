package controller

import (
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)



func PostVoteController(c *gin.Context){
	//参数校验
	p := new(models.ParamVoteData)
	if err:=c.ShouldBindJSON(p);err!=nil{
		errs,ok := err.(validator.ValidationErrors)
		if !ok{
			zap.L().Error("c.ShouldBindJSON出错,但不是validator错误",zap.Error(err))
			ResponseError(c,CodeInvalidParam)
			return
		}
		zap.L().Error("c.ShouldBindJSON出错",zap.Error(err))

		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c,CodeInvalidParam,errData)
		return
	}
	userID,err:= getCurrentUserID(c)
	if err!=nil{
		ResponseError(c,CodeNeedLogin)
	}
	if err:=logic.VoteForPost(userID,p);err!=nil{
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	ResponseSuccess(c,nil)
}