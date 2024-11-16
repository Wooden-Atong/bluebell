package models

import "time"

//🌟内存对齐 节省内存，这里定义的时候相同类型变量放在一起❓ （挖坑待填）
type Post struct{
	//🌟注意，前端的数字表示大小要比int64小，所以可能出现id传递失真的问题。
	//🌟一般的解决方法就是json序列化的时候把user_id转为字符串再传递，反序列化的时候也把前端传来的字符串变为int64，在go语言中直接在json-tag中加一个string就好了
	ID int64 `json:"id,string" db:"post_id"`
	AuthorID int64 `json:"author_id" db:"author_id"`
	CommunityID int64 `json:"community_id" db:"community_id" binding:"required"`
	Status int32 `json:"status" db:"status"`
	Title string `json:"title" db:"title" binding:"required"`
	Content string `json:"content" db:"content" binding:"required"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}


//ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct{
	AuthorName string `json:"author_name"`
	VoteNum int64 `json:"vote_num"`
	//🌟实现继承，将post和community_detail嵌入，实现信息的拓展
	*Post
	*CommunityDetail `json:"community"` //🌟这样子可以再创一层json。相当于CommunityDetail的信息都在community字段下了
}