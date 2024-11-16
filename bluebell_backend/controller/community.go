package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//和社区相关


func CommunityHandler(c *gin.Context){
	//返回（查询）所有的社区选项 community_id, community_name
	data, err := logic.GetCommunityList()
	if err != nil{
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c,data)
}

//CommunityDetailHandler 根据分类id查询社区详情
func CommunityDetailHandler(c *gin.Context){
	//获取社区id
	idStr := c.Param("id")//🌟获取path参数，实际上就是从url中获取path参数，所以c.Param("xxx")要是url中:后面的xxx对应上
	id, err:= strconv.ParseInt(idStr,10,64)//🌟相当于参数校验并转化，将str转为10进制int64，如果不能转还会报错说明这个参数有问题
	if err!=nil{
		ResponseError(c,CodeInvalidParam)
		return
	}
	//返回（查询）社区细节
	data, err := logic.GetCommunityDetail(id)
	if err != nil{
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c,data)
}
