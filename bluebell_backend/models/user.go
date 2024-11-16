package models

type User struct{
	UserID int64 `db:"user_id"`//🌟雪花算法生成ID 用int64
	Username string `db:"username"`
	Password string `db:"password"`
	AToken string
	RToken string 

}