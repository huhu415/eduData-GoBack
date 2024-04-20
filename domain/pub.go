package domain

// LoginForm POST中body的内容, 这个结构体只用在判断表单完整性来使用
type LoginForm struct {
	Username    string `form:"username" binding:"required"`
	Password    string `form:"password" binding:"required"`
	School      string `form:"school" binding:"required"`
	StudentType int    `form:"studentType" binding:"required,min=1,max=2"` // 1 本科生 2 研究生, 不可以是其他的
}
