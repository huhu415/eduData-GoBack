package route

import (
	"eduData/api/controller"
	"eduData/api/middleware"
	"eduData/repository"
	"eduData/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewUpDataRouter(db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewRepository(db)
	lc := &controller.SigninController{
		LoginUsecase: usecase.NewUsecase(ur),
	}

	// 登录路由
	group.POST("/signin", lc.LogIn)

	// 需要认证的路由组
	authGroup := group.Use(middleware.RequireAuthJwt())
	{
		// 课程相关路由
		courseGroup := group.Group("/courses")
		{
			courseGroup.POST("/renew", lc.UpdateCourse)  // 更新课程
			courseGroup.POST("/add", lc.AddCourse)       // 添加课程
			courseGroup.POST("/:week", lc.GetWeekCourse) // 获取某一周的课程
		}

		// 成绩相关路由
		gradeGroup := group.Group("/grades")
		{
			gradeGroup.POST("/renew", lc.UpdateGrade) // 更新成绩
			gradeGroup.POST("/all", lc.GetGrade)      // 获取成绩
		}

		// 时间表相关路由
		authGroup.POST("/timetable", lc.GetTimeTable) // 获取时间表
	}
}
