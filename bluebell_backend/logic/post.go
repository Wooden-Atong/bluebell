package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//生成post_id
	p.ID = snowflake.GenID()
	//保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
}

func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {

	//获取post信息
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	//根据作者id查询作者名称
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}
	//根据社区id查询社区信息
	community_detail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysel.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}
	//🌟这里可以不要用data = new(models.ApiPostDetail)这个在前面初始化，因为反正要复制，不如直接最后初始化并赋值，当然先new一个也可以
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		CommunityDetail: community_detail,
		Post:            post,
	}
	return
}

func GetPostList(pageNum int64, pageSize int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, errs := mysql.GetUserByID(post.AuthorID)
		if errs != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			err = errs
			return
		}
		//根据社区id查询社区信息
		community_detail, errs := mysql.GetCommunityDetailByID(post.CommunityID)
		if errs != nil {
			zap.L().Error("mysel.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			err = errs
			return
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			CommunityDetail: community_detail,
			Post:            post,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	//2.去redis中查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}

	//3.根据id去数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	for idx, post := range posts {
		user, errs := mysql.GetUserByID(post.AuthorID)
		if errs != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			err = errs
			return
		}
		//根据社区id查询社区信息
		community_detail, errs := mysql.GetCommunityDetailByID(post.CommunityID)
		if errs != nil {
			zap.L().Error("mysel.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			err = errs
			return
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			CommunityDetail: community_detail,
			Post:            post,
		}
		data = append(data, postDetail)
	}
	return data, err

}

func GetCommunityPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//2.去redis中查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}

	// 3.根据id去数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	for idx, post := range posts {
		user, errs := mysql.GetUserByID(post.AuthorID)
		if errs != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			err = errs
			return
		}
		//根据社区id查询社区信息
		community_detail, errs := mysql.GetCommunityDetailByID(post.CommunityID)
		if errs != nil {
			zap.L().Error("mysel.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			err = errs
			return
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			CommunityDetail: community_detail,
			Post:            post,
		}
		data = append(data, postDetail)
	}
	return data, err

}



func GetPostListNew(p *models.ParamPostList)(data []*models.ApiPostDetail, err error){
	if p.CommunityID==0{ //如果没有传communityID参数
		data, err = GetPostList2(p)
	}else{
		data, err = GetCommunityPostList2(p)
	}
	if err!=nil{
		zap.L().Error("GetPostListNew failed",zap.Error(err))
		return nil,err
	}
	return
}