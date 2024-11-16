package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)


func getIDsFormKey(key string, page, size int64)([]string,error){
	//确定查询的索引
	start := (page-1)*size
	end := start + size - 1
	return client.ZRevRange(key,start,end).Result()//ZRevRange是按照从大到小排序
}

func GetPostIDsInOrder(p *models.ParamPostList)([]string, error){
	//从redis获取id
	 // 根据url带的参数查询
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore{
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFormKey(key,p.Page,p.Size)

}

// 按社区根据ids查询每篇帖子的赞成票数据
func GetCommunityPostIDsInOrder(p *models.ParamPostList)([]string, error){

	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore{
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 社区的key（是一个set的键）
	cKey := getRedisKey(KeyCommunitySetPF+strconv.Itoa(int(p.CommunityID)))
	// 🌟利用缓存key减少zinterstore执行的次数，不需要一上来全部zinterstore，而是涉及用community_id去看按order的帖子的时候，才去创建。
	//🌟使用zinterstore把分区的帖子set与帖子分数的zset生成一个新的zset，针对新的zset再按之前的逻辑取数据
	// 可以就是实际就是post:time:1或post:score:1这样子
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val()<1{
		pipeline := client.Pipeline()
		//如果还不存在，则需要创建
		pipeline.ZInterStore(key,redis.ZStore{
			Aggregate: "MAX",
		},cKey,orderKey)//set 和 zset取交集，返回的是zset
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_,err := pipeline.Exec()
		if err!=nil{
			return nil,err
		}
	}

	return getIDsFormKey(key,p.Page,p.Size)

}

//根据ids查询每篇帖子点赞成票的数据
func GetPostVoteData(ids []string)(data []int64, err error){
	// data = make([]int64, 0,len(ids))
	// for _,id :=range ids{
	// 	key:=getRedisKey(KeyPostVotedZSetPre+id)
	// 	//获取赞成票的数量(也就是分数为1的数量)
	// 	v1 := client.ZCount(key,"1","1").Val()
	// 	data = append(data, v1)
	// }
	//🌟注意，上述写法向redis发的请求太多太频繁了，实际上可以全部一起发过去，然后redis查询到之后再统一返回,减少RTT
	//🌟可以用pipeline封装到事务中去
	pipeline := client.Pipeline()
	for _,id := range ids{
		key:=getRedisKey(KeyPostVotedZSetPre+id)
		pipeline.ZCount(key,"1","1")
	} 
	cmders, err :=pipeline.Exec()
	if err!=nil{
		return nil,err
	}

	data = make([]int64, 0,len(cmders))
	for _,cmder := range cmders{
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	} 
	return
}

