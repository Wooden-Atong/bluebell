package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600 //一周时间的秒数
	scorePerVote = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeat = errors.New("不允许重复投票")
)

func CreatePost(postID,communityID int64)error{
	//🌟事务操作，必须同时成功
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet),redis.Z{
		Score:float64(time.Now().Unix()),
		Member: postID,
	})//🌟由于不需要返回结果，所以后面不用加上.Result()
	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet),redis.Z{
		Score:float64(time.Now().Unix()),//帖子分数由发布时间初始化
		Member: postID,
	})//.Result()

	//把帖子id加到community得set
	cKey := getRedisKey(KeyCommunitySetPF+strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey,postID)
	_,err := pipeline.Exec()
	return err
}

//用户投一票就加432分 86400s/200 需要200张赞成票才可以给帖子续一天（时间戳增大一天，否则排后面了） 《redis实战》

// 投票的几种情况：
// direction=1
//   之前没有投过票，现在投赞成票 --> 更新分数和投票记录 abs 1 +432
//   之前投反对票，现在改投赞成票 --> 更新分数和投票记录 abs 2 +432*2
// direction=0
//   之前投过赞成票，现在要取消投票 --> 更新分数和投票记录 1 -432
//   之前投过反对票，现在要取消投票 --> 更新分数和投票记录 1 +432
// direction=-1
//   之前没有投过票，现在投反对票 --> 更新分数和投票记录 1 -432
//   之前投赞成票，现在改投反对票 --> 更新分数和投票记录 2 -432*2

// 投票的限制：
//每个帖子自发表之日后的一个星期允许用户投票，超过一个星期不允许再投票。
//到期之后讲redis中保存的赞成票数及反对票数存储到mysql
//到期之后删除 KeyPostVotedZSetPre

func VoteForPost(userID string, postID string, direction float64) error {
	//1.判断是否已经过期
	//🌟这里的key是查找zset的key，而这里的postID则是键值对中的键
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	
	//2.更新帖子分数
	// 先查该用户之前给当前帖子的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPre+postID), userID).Val()
	if ov == direction{
		return ErrVoteRepeat
	}

	var dir float64
	if direction > ov{
		dir = 1
	}else{
		dir = -1
	}
	diff := math.Abs(ov-direction) 
	//2和3需要放在一个事务中去
	//🌟就是修改指定zset中指定key对应的value
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet),dir*diff*scorePerVote,postID)
	
	//3.记录该用户给帖子投票的数据
	if direction==0{//取消投票
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPre+postID),postID).Result() //🌟这里是直接移除了这一个投票项，也就是不管之前投的啥都没了
	}else{
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPre+postID),redis.Z{
		Score: direction,//🌟这里不是前面的分数，这里只是记录票数,赞成票还是反对票
		Member: userID,
	}).Result()}
	_,err:=pipeline.Exec()
	return err
}
