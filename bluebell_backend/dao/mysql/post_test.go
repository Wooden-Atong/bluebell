package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

//🌟这里需要init函数因为要测试的createPost这个函数当中调用了db，而db如果需要提前初始化
func init(){
	dbCfg :=settings.MySQLConfig{
		Host        :"127.0.0.1",
		User        :"root",
		Password    :"133248563s",
		DB          :"bluebell",
		Port        :3306,
		MaxOpenConns :10,    
		MaxIdleConns :10,  
	}
	err := Init(&dbCfg)
	if err!=nil{
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:10,
		AuthorID: 123,
		CommunityID: 1,
		Title: "test",
		Content: "just a test",
	}
	err := CreatePost(&post)
	if err!=nil{
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n",err)
	}
	t.Logf("CreatePost insert record into mysql success")
}