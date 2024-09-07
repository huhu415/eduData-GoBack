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

	group.POST("/signin", lc.LogIn) // 登陆
	r := group.Use(middleware.RequireAuthJwt())
	{
		r.POST("/updata", lc.Updata)           // 更新课程
		r.POST("/updataGrade", lc.UpdataGrade) // 更新成绩
	}
}
