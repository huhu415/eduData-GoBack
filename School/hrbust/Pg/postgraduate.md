### 哈理工研究生爬虫方法

我们登陆的目标是获得激活后的cookie, 因为请求中带着激活后的cookie可以访问任何我们想要访问的页面

所以要完成 2 个事情, 找到cookie, 然后激活cookie.


## 第一次请求
目的 : 找到__VIEWSTATE, __EVENTVALIDATION, ValidateCodeSrc
#### 发送
- 方法: GET
- url: http://yjsjw.hrbeu.edu.cn/
- headers:
    - ```User-Agent``` : 自定义
- body: 空

#### 接收
- headers: 不重要
- body:
    - ```__VIEWSTATE``` : 可以被找到
    - ```__EVENTVALIDATION``` : 可以被找到
    - ```ValidateCode``` : 可以被找到

## 第二次请求
目的 : 通过ValidateCodeSrc获取验证码图片和cookie
#### 发送
- 方法: GET
- url: http://yjsjw.hrbeu.edu.cn/ + ValidateCodeSrc
- headers:
    - ```User-Agent``` : 自定义
- body: 空

#### 接收
- headers:
    - ```Set-Cookie``` : 可以被找到
- body:
    - 验证码图片, 被base64编码的, 并且带有前缀, 要识别的话要把前缀去掉, [识别](#识别验证码)后得到ValidateCode(4位数字)


## 第二次请求
目的 : 把收集到的__VIEWSTATE, __EVENTVALIDATION, ValidateCode, cookie发送给服务器, 来激活cookie
#### 发送
- 方法: POST
- url: http://yjsjw.hrbeu.edu.cn/
- headers:
    - ```User-Agent``` : 自定义
    - ```Cookie``` : 从第二次请求中获取
    - ```Content-Type``` : ```application/x-www-form-urlencoded; charset=UTF-8```
    - ```Referer``` : http://yjsjw.hrbeu.edu.cn/
    - ```Origin``` : http://yjsjw.hrbeu.edu.cn
    - ```X-MicrosoftAjax``` : ```Delta=true```
    - ```X-Requested-With``` : ```XMLHttpRequest```
- body:
    - ```__VIEWSTATE``` : 第一次请求中获取
    - ```____EVENTVALIDATION``` : 第一次请求中获取
    - ```__UserName``` : 教务处账号
    - ```__PassWord``` : 教务处密码
    - ```__ValidateCode``` : 从第二次请求验证码图片中获取
    - ```__ScriptManager1``` : ```UpdatePanel2|btLogin```
    - ```____EVENTTARGET``` : ```btLogin```
    - ```__drpLoginType``` : ```1```
    - ```____ASYNCPOST``` : ```true```
#### 接收
- headers: 根据长短判断成功与否
- body: 提示信息, 如验证码错误, 密码不正确, 用户名不存在等
> **以上登陆完成还差识别验证码没有讲.**
>
> 研究生登陆所需代码在[sign_in](sign_in/sign_in_pg.go)文件夹中


## 识别验证码
我们在登陆的时候, 获取到了验证码图片, 那么我们要识别验证码图片, 以便登陆.

研究生的识别, 我是使用了百度的, 因为有个活动, 1个月1000次免费识别, 所以我就用了百度的.
使用方法如下, 还有凭证要一个月更新一次.[百度手写数字识别](https://ai.baidu.com/tech/ocr_others/numbers)
> **识别完验证码, 就可以获取正式页面了**
>
> baidu识别验证码代码在[identimage](../../identimage/baidu.go)文件夹中


## 访问正式页面
得到激活后的cookie后, 就可以访问正式页面了, 但是要注意,
每次访问正式页面都要带上激活后的cookie:`cookie + LoginType=LoginType=1;`, 否则会被重定向到登陆页面.

下面的地址前面加上(http://yjs.hrbust.edu.cn/Gstudent), 就可以访问到正式页面了
- 教务处左侧栏的全部网页地址有如下:
    - [学期注册信息管理](TrainManage/StudentOnLineReg.aspx?EID=5ENwdLASepv-gHHvpTmBHUbQUAiQCaDTY!ZyJxLCIzloJf59w4ZXla3h3!aa7U5T&UID=)
    - [培养计划信息管理](TrainManage/TeachPlanAddEdit.aspx?EID=k6uTK6gTGwIV3oPm699!TY0kQ-T8QRvEkRi-unAab!Jecic9rJBzlvbqECRTCxIH&UID=)
    - [培养计划信息查询](TrainManage/TeachPlanQuery.aspx?EID=4DR1McYZMWyLdAy1UFkgaKVYM!2f7FLn3YLcGf!Fpf4QqPmxydWoow==&UID=)
    - [师生互选信息管理](TrainManage/StudentSelTutor.aspx?EID=mIB-G8aUXQQOj5s29wiQm3gGWEk3FTuCzxIjS6iCvR1cCQsdy6sTrA==&UID=)
    - [开课目录信息查询](Course/CourseOpenDirQuery.aspx?EID=Mq1AYFLvieWK6jW2XnIkKG-NUDrEkLm3PEZovY6HqgJzloHHEYhZdQ==&UID=)
    - [课程网上选课管理](Course/PlanCourseOnlineSel.aspx?EID=GNyhGO4r6Rgg2dI0j7mSb09KRa2ge11OpShhpbxeiPOKTJSrLTtdRA==&UID=)
    - [选课结果信息查询](Course/CourseSelQuery.aspx?EID=dRSltRtpla1Km67tvn5StQCMFK3Bi-Y7IEJqewCEZqeZJfWSyfLIkA==&UID=)
    - [学期课表信息查询](Course/StuCourseQuery.aspx?EID=pLiWBm!3y8J!emOuKhzHa3uED3OEJzAvyCpKfhbkdg9RKe9VDAjrUw==&UID=)
    - [本周课表信息查询](Course/StuCourseWeekQuery.aspx?EID=vB5Ke2TxFzG4yVM8zgJqaQowdgBb6XLK0loEdeh1pyPrNQM0n6oBLQ==&UID=)
    - [课程成绩信息查询](Course/StudentScoreQuery.aspx?EID=VuYUA7YP5gRRL6Z-IeJgBXS10bXlTWXy-qmX4GxB8li4o6gB-9Zv6w==&UID=)
    - [调停补课信息查询](Course/LessonTTBInfo.aspx?EID=JXNe21ncNv5c741DdkRmRpTPjss4Pwm5wjEnZochoom28sIr7KH7aw==&UID=)
    - [学期考试信息查询](Course/CourseTestInfo.aspx?EID=fZsh-JoXwN3qRxUwRF6mrd5D3na6hHMg8C4!fm0GwGz9XdNhgSGgXQ==&UID=)
    - [课程免修申请管理](Course/Exemption_Apply.aspx?EID=6faXodlwt2Nl2pKsVJ-!nwTmVbNW!!mMD1289ofPc5nm6FHqtrr2zQ==&UID=)
    - [课程重修申请管理](Course/Restudy_Apply.aspx?EID=lq0RV!jpMvkIz0Yr3xnnjZeT9kyvVukoupPNzVGfH8jxG6XSU6osTg==&UID=)
    - [入学CET成绩录入](English/EnglishBeScoreInput.aspx?EID=Bbd1cm9i6CE68aVI4jB7GpM3NiMzkc!DjfRsoYZmiPiya7NKSKN7yw==&UID=)
    - [CET外语考级报名](English/English46SignUp.aspx?EID=xyMe-U-eBwjvs8mWQT5RlMdK8nh51YIu2ZuN6vw-25DUmd!7VfRGxw==&UID=)
    - [CET外语成绩查看](English/EnglishScoreView.aspx?EID=o0i5OMBvZ5P4nhV8i2i3sEtdTdcZqcar2ImQNKVV2SrQgfX7v3WBzA==&UID=)
    - [专业实践](TrainManage/SocialPracticeReg.aspx?EID=!ft1qj7BTJ-B3KlUapum9KwVELHy36qh9q4jCSSjYSpZ0CtqJ6xYTNhorBWjv!7E&UID=)
    - [知识产权与信息检索](TrainManage/TeachPracticeReg.aspx?EID=ky!ONtrVsw8F3hSSC93tt5veKNZ245sK55J4CayDIWxt9pJU!dxCkqKLRdUaMwD7&UID=)
    - [论文开题及阶段报告](TrainManage/BooksReadReg.aspx?EID=WCRuInELBfd!ruRnd4syai5VjLzuA8g-jehfUrqEClVyh2jbhETHQA==&UID=)
    - [前沿课题讲座](TrainManage/SciReport.aspx?EID=CDZHCMtrTAY0FdkHHVhBrpj-yFO8qH85OvTgLPR4er7qnLi0h0gqHA==&UID=)
    - [学术活动信息登记](TrainManage/ScienceReportReg.aspx?EID=lAOCL!KzkTa6E0O4ErgvWl41Tg1KnoCCgB4YzeI6bW0y6LvMu1WPaBwgQpVKuQZK&UID=)
- 以上这些链接前面要加[教务处首页](http://yjs.hrbust.edu.cn/), 后面要加学号, 并且请求头要加上cookie值, 就可以获取到html页面了
- 这些地址可能会变, 可以从[教务处首页](http://yjs.hrbust.edu.cn/)+[这个地址](Gstudent/leftmenu.aspx?UID=)+学号就可以获取到最新的侧边栏的所有地址了

[本周课表信息查询](Course/StuCourseWeekQuery.aspx?EID=vB5Ke2TxFzG4yVM8zgJqaQowdgBb6XLK0loEdeh1pyPrNQM0n6oBLQ==&UID=)
在请求的时候如果在请求头中的cookie: ```cookie + LoginType=LoginType=1; + DropDownListWeeks=DropDownListWeeks=5;```

就会查询第5周的课表, 这样可以精细化调整查询

与本科生不同, 研究生的课表是一周一周的获取, 因为这样便于解析, 可能日后会更改, 因为这样浪费服务器资源.
> **获取到html后, 我们需要解析html, 才能把数据放到数据库, 以便简单的增删改查**
>
> 获取html所需代码在[htmlgetter](htmlgetter/htmlgetter.go)文件夹中


## 解析html
解析html, 我写了两种方式
1. 解析成二维数组, 每一个元素里面就是一节课
   - 但我发现这样是不行的, 因为有些课程是跨行的, 如果只是二维数组, 再如果有的课程表是带有圆角的, 那么跨行的课程就没办法整齐的拼接在一起了.
2. 解析成结构体数组, 每一个元素就是一节课
   - 这样就可以解决上面的问题了, 便于读取, 还更方便的写前端的代码了.

```
type Course struct {
	ID                    uint    //主键
	StuID                 int64   //学号多少
	Week                  int     //第几周
	WeekDay               int     //星期几
	NumberOfLessons       int     //上第几节
	NumberOfLessonsLength int     //上课节数(大学一节课一般为2节, 高中为1节)
	CourseContent         string  //上课内容
	CourseLocation        string  //上课地点
	TeacherName           string  //老师
	BeginWeek             int     //开始周
	EndWeek               int     //结束周
}
```

如果有一周需要上5节课, 那么就是5个结构体的数组.
> **解析完html, 只需要连接数据库, 将数据发送给数据库就可以了.**
>
> 研究生解析所需代码在[parse_form](parse_form/parse_form.go)文件夹中