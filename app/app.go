// Package app 处理路由的逻辑
package app

import (
	"errors"
	"net/http/cookiejar"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"

	"eduData/database"
	hget "eduData/htmlgetter"
	"eduData/middleware"
	pf "eduData/parse_form"
)

const (
	// 控制研究生获取页面用的
	// LEFTCORUSE 某一周的
	//LEFTCORUSE = "Course/StuCourseWeekQuery.aspx?EID=vB5Ke2TxFzG4yVM8zgJqaQowdgBb6XLK0loEdeh1pyPrNQM0n6oBLQ==&UID="
	// LEFTCORUSEALL 学期的
	LEFTCORUSEALL = "Course/StuCourseQuery.aspx?EID=pLiWBm!3y8J!emOuKhzHa3uED3OEJzAvyCpKfhbkdg9RKe9VDAjrUw==&UID="
)

// 根据学校和研究生本科生判断获取html并解析
func judgeUgOrPgGetInfo(c *gin.Context, cookieJar *cookiejar.Jar) ([]database.Course, error) {
	var table []database.Course
	switch c.PostForm("school") {
	// 哈理工
	case "hrbust":
		switch c.PostForm("studentType") {
		case "1":
			ugHTML, errUg := hget.RevLeftChidUg(cookieJar, "2000")
			if errUg != nil {
				return nil, errUg
			}
			table, errUg = pf.ParseTableUgAll(ugHTML)
			if errUg != nil {
				return nil, errUg
			}
		case "2":
			pgHTML, errPg := hget.RevLeftChidPg(cookieJar, c.PostForm("username"), LEFTCORUSEALL)
			if errPg != nil {
				return nil, errPg
			}
			table, errPg = pf.ParseTablePgAll(pgHTML)
			if errPg != nil {
				return nil, errPg
			}
		}
	// 东北农业大学
	case "neau":
		switch c.PostForm("studentType") {
		case "1":
			GetJSONneau, errNeau := hget.GetJSONneau(cookieJar, "2023-2024-2-1") // todo 设计一下获取学期的函数
			if errNeau != nil {
				return nil, errNeau
			}
			table, errNeau = pf.Parse_json_ug_nd(GetJSONneau)
			if errNeau != nil {
				return nil, errNeau
			}
		case "2":
			return nil, errors.New(c.PostForm("school") + "研究生登陆功能还未开发")
		}
	// 其他没有适配的学校
	default:
		return nil, errors.New("不支持的学校")
	}
	return table, nil
}

// Signin 下发jwt的cookie
func Signin(c *gin.Context) {
	// 获取表单中的用户名
	username := c.PostForm("username")

	//创建jwt
	j := middleware.NewJWT()
	tokenString, err := j.CreateToken(jwt.MapClaims{
		"username": username,
	})
	if err != nil {
		err = c.Error(errors.New("app.Signin()函数中CreateToken的错误: " + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "请重试",
		})
		return
	}

	//返回给前端结果
	c.SetCookie("authentication", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "登陆成功",
	})
}

// UpdataDB 更新数据库中的课程表
func UpdataDB(c *gin.Context) {
	cookieAny, ok := c.Get("cookie")
	if !ok {
		_ = c.Error(errors.New("app : signin中间件没有正常设置username, cookie, studentType")).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "请重试",
		})
		return
	}
	cookie := cookieAny.(*cookiejar.Jar)

	//删除数据库中的所有这个用户名的课程
	database.DeleteUserAllCourse(c.PostForm("username"), c.PostForm("school"))

	table, err := judgeUgOrPgGetInfo(c, cookie)
	if err != nil {
		_ = c.Error(errors.New("app.UpdataDB()函数中judgeUgOrPgGetInfo的错误: " + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "获取课表错误, 请删除小程序后重试",
		})
		return
	}

	//把解析到的课程与学号都一起添加到数据库中
	database.AddCourse(table, c.PostForm("username"))
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "课程更新成功",
	})

}

// GetWeekCoure 获取某周课程表
func GetWeekCoure(c *gin.Context) {
	//获取url中的周数的参数
	week := c.Param("week")
	weekInt, err := strconv.Atoi(week)
	if err != nil {
		_ = c.Error(errors.New("app.GetWeekCoure():周数必须是数字")).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "周数必须是数字",
		})
		return
	} else if !(weekInt >= 0 && weekInt <= 30) {
		_ = c.Error(errors.New("app.GetWeekCoure():周数必须在0-30之间")).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "周数必须在0-30之间",
		})
		return
	}
	//fmt.Printf("GetWeekCoure()中的username: %s, week:%d\n", username, weekInt)

	//获取数据库中的课程
	courseByWeekUsername := database.CourseByWeekUsername(weekInt, c.PostForm("username"), c.PostForm("school"))
	c.JSON(http.StatusOK, courseByWeekUsername)
}

/*
// UpdataDB2 理论上高性能版获取课程表, 20个协程
func UpdataDB2(c *gin.Context) {
	//检测表单是否合法
	var loginForm LoginForm
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("UpdataDB()中的username: %s, password: %s\n", loginForm.Username, loginForm.Password)

	//登陆
	cookie, err := signin.SingIn(loginForm.Username, loginForm.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	fmt.Println(cookie, "登陆成功")

	//删除数据库中的所有这个用户名的课程
	database.DeleteUserAllCourse(loginForm.Username)

	//获取全部课程表
	var wg sync.WaitGroup
	wg.Add(20)
	resultChan := make(chan error, 25)
	for i := 1; i <= 20; i++ {
		go func(w int) {
			defer wg.Done()
			//获取课程表
			leftChid, routineerr := hget.RevLeftChid(cookie, loginForm.Username, LEFTCORUSE, "DropDownListWeeks=DropDownListWeeks="+strconv.Itoa(w))
			if routineerr != nil {
				resultChan <- routineerr
				return
			}

			//解析课程表
			table, routineerr := pf.ParseTable1D(leftChid)
			if routineerr != nil {
				resultChan <- routineerr
				return
			}

			//把解析到的课程与学号都一起添加到数据库中
			database.AddCourse(table, loginForm.Username, loginForm.StudentType)
		}(i)
	}
	wg.Wait()
	close(resultChan)

	if err = <-resultChan; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "课程更新成功",
		})
	}
}

// UpdataDB1 理论上低性能版获取课程表, 5个协程来回工作
func UpdataDB1(c *gin.Context) {
	//检测表单是否合法
	var loginForm LoginForm
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("app.UpdataDB1()函数中username: %s, password: %s, studentType: %d\n", loginForm.Username, loginForm.Password, loginForm.StudentType)

	// 根据本科生还是研究生判断登陆
	cookie, err := JudgeUgOrPgSignIn(loginForm)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	fmt.Println(cookie, "数据库登陆成功")

	//删除数据库中的所有这个用户名的课程
	database.DeleteUserAllCourse(loginForm.Username)

	// 如果是研究生
	if loginForm.StudentType == 2 {
		//获取全部课程表
		//创建一个channel, 5个goroutine
		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errs, ctx := errgroup.WithContext(ctx)
		msg := make(chan int, 20)
		for i := 0; i < 5; i++ {
			errs.Go(func() error {
				for {
					select {
					case <-ctx.Done():
						return nil
					case w := <-msg:
						wg.Add(1)
						leftChid, routineerr := hget.RevLeftChid(cookie, loginForm.Username, LEFTCORUSE, "DropDownListWeeks=DropDownListWeeks="+strconv.Itoa(w))
						if routineerr != nil {
							wg.Done()
							return routineerr
						}
						//解析课程表
						table, routineerr := pf.ParseTable1D(leftChid)
						if routineerr != nil {
							wg.Done()
							return routineerr
						}

						//把解析到的课程与学号都一起添加到数据库中
						database.AddCourse(table, loginForm.Username)

						wg.Done()
					default:
						continue
					}
				}
			})
		}
		//从第几周开始到第几周, 一般来说是从1-20
		for i := 1; i <= 20; i++ {
			msg <- i
			fmt.Print("send:", i, "\n")
		}
		wg.Wait()
		cancel()
		if errs.Wait() != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errs.Wait().Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "课程更新成功",
	})

}
*/
