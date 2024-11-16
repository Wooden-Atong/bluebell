package models


//🌟对于请求Post的参数大多放在这里,tag里是json和binding，Get则少一些，是和binding；
//🌟post参数校验利用gin内置validator，所以都需要在tag中加入binding字段

const (
	OrderTime = "time"
	OrderScore = "score"
)

//ParamSignUp 注册请求参数
type ParamSignUp struct{
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`//eqfield=Password表明这个字段要和RePassword相等
}

//ParamLogin 登陆请求参数
type ParamLogin struct{
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票请求参数 
type ParamVoteData struct{
	//🌟UserID 从请求中获取当前用户，不需要显示定义
	PostID string `json:"post_id" binding:"required"`
	Direction int8 `josn:"direction,string" binding:"oneof=0 1 -1"` //赞成票（1）、反对票（-1）、取消投票（0）
	//🌟多个binding检测中间用逗号隔开，不能有多余空格. 同时注意如果binding中出现了required，这个值是false或者0或者""，他会自动忽略掉，认为你没有填值
}

//ParamPostList 获取帖子列表query string参数
type ParamPostList struct{
	CommunityID int `form:"community_id"` //可以为空
	Page int64 `form:"page"`
	Size int64 `form:"size"`
	Order string `form:"order"`
}

//
type ParamCommunityPostList struct{
	*ParamPostList
}