package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	// "go.uber.org/zap"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// 🌟jwt包自带的jwt.StandardClaims只包含了官方字段，我们这里需要额外记录字段，所以要自定义结构体，保存更多信息也是可以的，只要是不敏感的有助于识别用户的数据都可以
type MyClaims struct {
	UserID int64 `json:"userid"`
	Username string `json:"username"`
	jwt.StandardClaims
}


var MySecret = []byte("压心底压心底不能告诉你")


// GenToken 根据token type 获取对应token
func GenClaim(userID int64,username string) (*MyClaims){
	c := &MyClaims{
		userID,
		username, // 自定义字段
		jwt.StandardClaims{//🌟这里有七个数据，这里只指定了两个
			//🌟viper.GetInt()从配置信息中获取整数。
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_access_expire"))*time.Second).Unix(), // 10分钟过期时间
			Issuer:    "bluebell",                               // 签发人
			Subject: "aToken",
		},
	}

	return c
}


// GenTokenHandler 获取token处理器
func GenARToken(userID int64,username string) (aToken string, rToken string, err error) {
	
	
	//返回aToken和rToken
	// 获取aToken
	claim := GenClaim(userID,username)
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(MySecret)
	if err!=nil{
		return "","",err
	}
	// 获取rToken
	claim.ExpiresAt = time.Now().Add(time.Duration(viper.GetInt("auth.jwt_refresh_expire"))*time.Hour).Unix() //30天过期时间
	claim.Subject = "rToken"
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(MySecret)// 使用指定的签名方法创建签名对象,使用指定的secret签名并获得完整的编码后的字符串token
	if err!=nil{
		return "","",err
	}
	
	return aToken,rToken,nil

	
}


// GenTokenHandler 获取token处理器
func GenAToken(userID int64,username string) (aToken string, err error) {
	
	
	//返回aToken和rToken
	// 获取aToken
	claim := GenClaim(userID,username)
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(MySecret)
	if err!=nil{
		return "",err
	}
	
	return aToken,nil
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	mc := new(MyClaims)
	// 解析token
	token, _ := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	
	// zap.L().Debug("判断token_type类型，",zap.Any("token_type=",mc.Subject),zap.Any("解析错误=",err))

	//如果是aToken
	if mc.Subject == "aToken"{
		if token.Valid { // 校验token
			return mc, nil
		}
		return nil, errors.New("access_token expired")//过期返回
	}
	//如果是rToken
	if token.Valid { // 校验token
		return mc, errors.New("refresh_token not expired") //没过期但也按照错误返回
	}
	return nil, errors.New("refresh_token expired")//过期返回
}