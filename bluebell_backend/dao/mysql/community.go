package mysql

import (
	"bluebell/models"
	"database/sql"
	
	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := `select community_id, community_name from community`
	communityList = make([]*models.Community, 0,10)
	if err := db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows { //查询不到不算错误返回，只是写进warning日志
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (communityDetail *models.CommunityDetail, err error) {
	//❓（挖坑待填）注意结构体这里需要new，上面函数切片需要make
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select community_id,community_name,introduction,create_time from community where community_id = ?`
	if err = db.Get(communityDetail, sqlStr, id); err != nil { //🌟db.Get和db.Select有什么区别吗？？？
		if err == sql.ErrNoRows { //查询不到不算错误返回，只是写进warning日志
			zap.L().Warn("there is no this community_id in db")
			err = ErrorInvalidID
		}
	}
	return
}
