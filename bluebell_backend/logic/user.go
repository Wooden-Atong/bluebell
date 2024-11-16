package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

//存放业务逻辑代码

func SignUp(p *models.ParamSignUp) (err error) {

	//1.判断用户是否存在（dao层）
	if err := mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询出错
		return err
	}
	
	//2.生成UID
	userID := snowflake.GenID()
	// 构造一个User实例
	user := &models.User{
		UserID: userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3.写进数据库（dao层）
	return mysql.InsertUser(user)

}

func Login(p *models.ParamLogin)(user *models.User,err error){
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//🌟传递的user是指针，所以这里执行完后user就有userID了
	if err=mysql.Login(user); err!=nil{
		return nil, err
	}
	//生成jwt
	aToken,rToken,err := jwt.GenARToken(user.UserID, user.Username)
	if err!=nil{
		return nil,err
	}
	user.AToken,user.RToken = aToken,rToken
	return 
}