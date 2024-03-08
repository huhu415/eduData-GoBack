// 本科生教务系统左侧菜单地址, 可以通过地址之间访问对应的页面
var nodes0 = [{
id: 1,
pId: 0,
name: '学期注册信息管理',
t: '学期注册信息管理',
url: 'TrainManage/StudentOnLineReg.aspx?EID=5ENwdLASepv-gHHvpTmBHUbQUAiQCaDTY!ZyJxLCIzloJf59w4ZXla3h3!aa7U5T&UID=2320410125',
target: 'PageFrame'
},
{
    id: 2,
    pId: 0,
    name: '培养计划信息管理',
    t: '培养计划信息管理',
    url: 'TrainManage/TeachPlanAddEdit.aspx?EID=k6uTK6gTGwIV3oPm699!TY0kQ-T8QRvEkRi-unAab!Jecic9rJBzlvbqECRTCxIH&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 3,
    pId: 0,
    name: '培养计划信息查询',
    t: '培养计划信息查询',
    url: 'TrainManage/TeachPlanQuery.aspx?EID=4DR1McYZMWyLdAy1UFkgaKVYM!2f7FLn3YLcGf!Fpf4QqPmxydWoow==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 4,
    pId: 0,
    name: '师生互选信息管理',
    t: '',
    url: 'TrainManage/StudentSelTutor.aspx?EID=mIB-G8aUXQQOj5s29wiQm3gGWEk3FTuCzxIjS6iCvR1cCQsdy6sTrA==&UID=2320410125',
    target: 'PageFrame'
}];
var nodes1 = [{
id: 1,
pId: 0,
name: '开课目录信息查询',
t: '开课目录信息查询',
url: 'Course/CourseOpenDirQuery.aspx?EID=Mq1AYFLvieWK6jW2XnIkKG-NUDrEkLm3PEZovY6HqgJzloHHEYhZdQ==&UID=2320410125',
target: 'PageFrame'
},
{
    id: 2,
    pId: 0,
    name: '课程网上选课管理',
    t: '课程网上选课管理',
    url: 'Course/PlanCourseOnlineSel.aspx?EID=GNyhGO4r6Rgg2dI0j7mSb09KRa2ge11OpShhpbxeiPOKTJSrLTtdRA==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 3,
    pId: 0,
    name: '选课结果信息查询',
    t: '选课结果信息查询',
    url: 'Course/CourseSelQuery.aspx?EID=dRSltRtpla1Km67tvn5StQCMFK3Bi-Y7IEJqewCEZqeZJfWSyfLIkA==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 4,
    pId: 0,
    name: '学期课表信息查询',
    t: '学期课表信息查询',
    url: 'Course/StuCourseQuery.aspx?EID=pLiWBm!3y8J!emOuKhzHa3uED3OEJzAvyCpKfhbkdg9RKe9VDAjrUw==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 5,
    pId: 0,
    name: '本周课表信息查询',
    t: '本周课表信息查询',
    url: 'Course/StuCourseWeekQuery.aspx?EID=vB5Ke2TxFzG4yVM8zgJqaQowdgBb6XLK0loEdeh1pyPrNQM0n6oBLQ==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 6,
    pId: 0,
    name: '课程成绩信息查询',
    t: '课程成绩信息查询',
    url: 'Course/StudentScoreQuery.aspx?EID=VuYUA7YP5gRRL6Z-IeJgBXS10bXlTWXy-qmX4GxB8li4o6gB-9Zv6w==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 7,
    pId: 0,
    name: '调停补课信息查询',
    t: '调停补课信息查询',
    url: 'Course/LessonTTBInfo.aspx?EID=JXNe21ncNv5c741DdkRmRpTPjss4Pwm5wjEnZochoom28sIr7KH7aw==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 8,
    pId: 0,
    name: '学期考试信息查询',
    t: '学期考试信息查询',
    url: 'Course/CourseTestInfo.aspx?EID=fZsh-JoXwN3qRxUwRF6mrd5D3na6hHMg8C4!fm0GwGz9XdNhgSGgXQ==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 9,
    pId: 0,
    name: '课程免修申请管理',
    t: '课程免修申请管理',
    url: 'Course/Exemption_Apply.aspx?EID=6faXodlwt2Nl2pKsVJ-!nwTmVbNW!!mMD1289ofPc5nm6FHqtrr2zQ==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 10,
    pId: 0,
    name: '课程重修申请管理',
    t: '课程重修申请管理',
    url: 'Course/Restudy_Apply.aspx?EID=lq0RV!jpMvkIz0Yr3xnnjZeT9kyvVukoupPNzVGfH8jxG6XSU6osTg==&UID=2320410125',
    target: 'PageFrame'
}];
var nodes2 = [{
id: 1,
pId: 0,
name: '入学CET成绩录入',
t: '入学CET成绩录入',
url: 'English/EnglishBeScoreInput.aspx?EID=Bbd1cm9i6CE68aVI4jB7GpM3NiMzkc!DjfRsoYZmiPiya7NKSKN7yw==&UID=2320410125',
target: 'PageFrame'
},
{
    id: 2,
    pId: 0,
    name: 'CET外语考级报名',
    t: 'CET外语考级报名',
    url: 'English/English46SignUp.aspx?EID=xyMe-U-eBwjvs8mWQT5RlMdK8nh51YIu2ZuN6vw-25DUmd!7VfRGxw==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 3,
    pId: 0,
    name: 'CET外语成绩查看',
    t: 'CET外语成绩查看',
    url: 'English/EnglishScoreView.aspx?EID=o0i5OMBvZ5P4nhV8i2i3sEtdTdcZqcar2ImQNKVV2SrQgfX7v3WBzA==&UID=2320410125',
    target: 'PageFrame'
}];
var nodes3 = [{
id: 1,
pId: 0,
name: '专业实践',
t: '专业实践',
url: 'TrainManage/SocialPracticeReg.aspx?EID=!ft1qj7BTJ-B3KlUapum9KwVELHy36qh9q4jCSSjYSpZ0CtqJ6xYTNhorBWjv!7E&UID=2320410125',
target: 'PageFrame'
},
{
    id: 2,
    pId: 0,
    name: '知识产权与信息检索',
    t: '知识产权与信息检索',
    url: 'TrainManage/TeachPracticeReg.aspx?EID=ky!ONtrVsw8F3hSSC93tt5veKNZ245sK55J4CayDIWxt9pJU!dxCkqKLRdUaMwD7&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 3,
    pId: 0,
    name: '论文开题及阶段报告',
    t: '论文开题及阶段报告',
    url: 'TrainManage/BooksReadReg.aspx?EID=WCRuInELBfd!ruRnd4syai5VjLzuA8g-jehfUrqEClVyh2jbhETHQA==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 4,
    pId: 0,
    name: '前沿课题讲座',
    t: '前沿课题讲座',
    url: 'TrainManage/SciReport.aspx?EID=CDZHCMtrTAY0FdkHHVhBrpj-yFO8qH85OvTgLPR4er7qnLi0h0gqHA==&UID=2320410125',
    target: 'PageFrame'
},
{
    id: 5,
    pId: 0,
    name: '学术活动信息登记',
    t: '学术活动信息登记',
    url: 'TrainManage/ScienceReportReg.aspx?EID=lAOCL!KzkTa6E0O4ErgvWl41Tg1KnoCCgB4YzeI6bW0y6LvMu1WPaBwgQpVKuQZK&UID=2320410125',
    target: 'PageFrame'
}];