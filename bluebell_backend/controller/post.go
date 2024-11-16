package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)




func CreatePostHandler(c *gin.Context){
	//获取参数及参数校验
	p := new(models.Post)
	//🌟c.ShouldBindJSON() 做校验实际上是利用gin内嵌的validator，根据binding tag指定的去校验
	if err := c.ShouldBindJSON(p);err!=nil{
		zap.L().Debug("c.ShouldBindJSON(p) error",zap.Any("err",err))//🌟学习一下用这个debug打印err
		ResponseError(c,CodeInvalidParam)
		return
	} 

	//获取当前登陆用户的ID
	userID,err := getCurrentUserID(c)
	if err!=nil{
		ResponseError(c,CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	//创建帖子
	if err:=logic.CreatePost(p);err!=nil{
		zap.L().Error("logic.CreatePost() failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c,nil)
}

//GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context){
	//获取参数（路径参数）
	pidStr:=c.Param("id")
	pid,err:=strconv.ParseInt(pidStr,10,64)
	if err!=nil{
		zap.L().Error("get post detail with invalid param",zap.Error(err))
		ResponseError(c,CodeInvalidParam)
		return
	}
	//根据id去除帖子数据
	data, err := logic.GetPostByID(pid)
	if err!=nil{
		zap.L().Error("logic.GetPostByID failed")
		ResponseError(c,CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c,data)
}

func GetPostListHandler(c *gin.Context){
	//获取分页参数
	pageNum,pageSize:=getPageInfo(c)


	//获取数据
	data,err:=logic.GetPostList(pageNum,pageSize)
	if err!=nil{
		zap.L().Error("logic.GetPostList failed!",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c,data)
}

//GetPostListHandler2 升级版帖子列表接口
//根据前端传来的参数（分数 或 创建时间）动态的获取帖子列表
func GetPostListHandler2(c *gin.Context){
	p := &models.ParamPostList{
		Page: 1,
		Size: 10,
		Order: models.OrderTime,
	}
	//1.获取参数（按时间 还是 按分数）
	if err:=c.ShouldBindQuery(p);err!=nil{
		zap.L().Error("GetPostListHandler2 with invalid param.", zap.Error(err),zap.String("order:",p.Order))
		ResponseError(c,CodeInvalidParam)
		return
	}
	

	//获取数据
	data,err:=logic.GetPostListNew(p)
	if err!=nil{
		zap.L().Error("logic.GetPostListNew failed!",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c,data)
}


//根据社区查询帖子列表
// func GetCommunityPostListHandler2(c *gin.Context){
// 	p := &models.ParamCommunityPostList{
// 		ParamPostList: &models.ParamPostList{
// 			Page: 1,
// 			Size: 10,
// 			Order: models.OrderTime,
// 		},
		
// 	}
// 	//1.获取参数（按时间 还是 按分数）
// 	if err:=c.ShouldBindQuery(p);err!=nil{
// 		zap.L().Error("GetPostListHandler2 with invalid param.", zap.Error(err),zap.String("order:",p.Order))
// 		ResponseError(c,CodeInvalidParam)
// 		return
// 	}
	

// 	//获取数据
// 	data,err:=logic.GetCommunityPostList2(p)
// 	if err!=nil{
// 		zap.L().Error("logic.GetPostList failed!",zap.Error(err))
// 		ResponseError(c,CodeServerBusy)
// 		return
// 	}
// 	//返回响应
// 	ResponseSuccess(c,data)
// }
