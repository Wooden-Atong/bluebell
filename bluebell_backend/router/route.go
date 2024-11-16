package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine{
	if mode == gin.ReleaseMode{
		gin.SetMode(gin.ReleaseMode)//🌟gin设置成发布模式，否则gin框架默认是debug模式，运行后输出GIN-debug信息
	}
	//新建一个没有任何默认中间件的路由
	r := gin.New()
	//注册全局中间件,logger为自定义中间件;令牌桶限速为每2s添加一个
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middleware.RateLimitMiddleware(2*time.Second,1))

	v1 := r.Group("api/v1")//🌟嵌套路由，起一个api/v1，方便后续拓展

	//注册 路由
	v1.POST("/signup", controller.SignUpHandler)

	//登陆 路由
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middleware.JWTAuthMiddleware())

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)//🌟/:id 就是加了一个路径参数，但实际请求来的时候不要写/:1，而是直接/1
		
		v1.POST("/post",controller.CreatePostHandler)
		v1.GET("/post/:id",controller.GetPostDetailHandler)
		v1.GET("/posts/",controller.GetPostListHandler)

		//根据时间或者分数获取帖子列表（升级版接口）
		v1.GET("/posts2/",controller.GetPostListHandler2)

		//投票
		v1.POST("/vote",controller.PostVoteController)
	}

	//🌟这里是为“/ping”路由单独注册中间件JWTAuthMiddleware()，如果在这个中间件执行过程中abort了，那么这个路由的处理函数将不会执行
	// r.GET("/ping",middleware.JWTAuthMiddleware(),func (c *gin.Context)  {
	// 	c.String(http.StatusOK,"pong")
	// })

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"msg":"404",
		})
	})

	return r

}

