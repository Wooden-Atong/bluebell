package models

import "time"

type Community struct{
	ID int64 `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`

}

type CommunityDetail struct{
	ID int64 `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
	Introduction string `json:"introduction,omitempty" db:"introduction"`
	CreateTime time.Time `json:"create_time" db:"create_time"` //🌟注意这里是time.Time的时间类型，连接数据库的时候要加上parseTime=true，让数据库自己做一个类型转换
}//❓（挖坑待填）如果db这里不指定tag，就会报错，create_time和CreateTime就没法绑定，这中间到底是一个什么界限，还需要仔细思考弄清楚。

