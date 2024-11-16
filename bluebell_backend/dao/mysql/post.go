package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.Post) error {
	sqlStr := `insert into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)`
	_, err := db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	if err != nil {
		return err
	}
	return err
}

func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}

func GetPostList(pageNum int64, pageSize int64) (posts []*models.Post, err error) {
	//🌟limit ？,？ 意思是从第一个?开始取，取第二个?条数据
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post limit ?,?`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (pageNum-1)*pageSize, pageSize) //🌟只获取一条数据用Get，获取多条用Select，❓（挖坑待太难）这个后面再系统整理一下
	return
}

// 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string)(postsList []*models.Post,err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
	from post  
	where post_id in (?) 
	order by FIND_IN_SET(post_id,?)` //🌟用了mysql内置的FIND_IN_SET进行排序
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))//sqlx是sql的扩展库，这里就是将ids一组值插入在sqlStr中
	if err!=nil{
		return nil,err
	}
	//🌟这里还需要再仔细看看咋个回事
	query = db.Rebind(query)
	db.Select(&postsList,query,args...)//args这里一定要加...
	return
}
