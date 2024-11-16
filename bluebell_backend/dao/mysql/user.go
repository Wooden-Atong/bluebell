package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)


const secrete = "limutong"


//存放操作数据库代码

//CheckUserExist 检查指定用户名的用户是否已经存在
func CheckUserExist(username string)(err error){
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	//🌟注意这里传入的是count地址，所以不需要思维定式一定要返回count
	if err := db.Get(&count, sqlStr, username); err!=nil{
		return err //报err的情况就不考虑count数了，直接false
	}
	if count >0{
		return ErrorUserExist//🌟一般不要直接在代码中出现奇怪字符串，比如errors.New("用户已存在")，可以统一定义在一组
	}
	return 
}

//🌟go中建议为首字母大写的函数（不仅对本包使用，还供外面包使用的函数）写注释的时候要带上这个函数名
// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User)(err error){
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_,err = db.Exec(sqlStr,user.UserID,user.Username,user.Password)
	return
}

func encryptPassword(oPassword string) string{
	h := md5.New()
	h.Write([]byte(secrete))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User)(err error){
	oPassword := user.Password //用户登录时传进来的密码
	sqlStr := `select user_id, username, password from user where username=?`
	err=db.Get(user, sqlStr, user.Username)//🌟这里查询后的结果会报存在在user结构体当中，所以在前面需要先提前把原始密码保存在oPassword中
	if err==sql.ErrNoRows{
		return ErrorUserNotExist
	}
	if err != nil{//走到这一步err说明是查询数据库失败了
		return err
	}
	password := encryptPassword(oPassword)
	if password != user.Password{
		return ErrorInvalidPassword
	}
	return
}


func GetUserByID(uid int64)(user *models.User,err error){
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?` 
	err = db.Get(user,sqlStr,uid)
	return
}