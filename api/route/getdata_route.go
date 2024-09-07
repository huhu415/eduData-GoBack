package route

import (
	"eduData/api/controller"
	"eduData/repository"
	"eduData/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewGetDataRouter(db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewRepository(db)
	lc := &controller.GetDataController{
		LoginUsecase: usecase.NewUsecase(ur),
	}

	group.POST("/getweekcoure/:week", lc.GetWeekCoure)
	group.POST("/getgrade", lc.GetGrade)
	group.POST("/getTimeTable", lc.GetTimeTable)
	group.POST("/addcoures", lc.AddCoures)

}
