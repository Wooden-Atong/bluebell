package redis

//redis key

const (
	//🌟对于redis的key，用:分割命名空间，方便业务查询和拆分。冒号不是唯一，只要能分割都可以，逗号斜线当然也可以
	KeyPrefix = "bluebell:" // 在redis集群中能快速找到这个项目相关的key
	KeyPostTimeZSet = "post:time" // 帖子 及 发帖时间
	KeyPostScoreZSet = "post:score" //帖子 及 投票分数 基准值是帖子的发帖时间的时间戳
	KeyPostVotedZSetPre = "post:voted:" //记录用户及投票的类型；参数是post-id
	KeyCommunitySetPF = "community:" // set;保存每个分区下帖子的id
)

// getRedisKey 给key加上前缀
func getRedisKey(key string)string{
	return KeyPrefix+key
}