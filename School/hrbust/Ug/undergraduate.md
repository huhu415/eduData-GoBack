### 哈理工本科生爬虫方法

我们通过最后一次的登陆请求, 可以看到请求过去的只有用户名, 密码, 验证码和cookie. 所以目的就是集齐所有.

## 第一次请求
目的 : 得到cookie
#### 发送
- 方法: GET
- url: http://jwzx.hrbust.edu.cn/academic/index.jsp
- headers:
  - ```User-Agent``` : 自定义
- body: 空

#### 接收
- headers:
  - ```set-cookie``` : 可以被找到
- body:不重要

## 第二次请求
目的 : 得到验证码图片
#### 发送
- 方法: GET
- url: http://jwzx.hrbust.edu.cn/academic/getCaptcha.do
- headers:
  - ```cookie``` : 第一次请求得到的cookie
  - ```User-Agent``` : 自定义
- body: 空

#### 接收
- headers: 不重要
- body:
  - ```body里面就是image, 是base64编码的```

## 第三次请求
目的 : 检验我们用码云平台识别到的数字是否正确
#### 发送
- 方法: POST
- url: http://jwzx.hrbust.edu.cn/academic/checkCaptcha.do?captchaCode=我们识别到的验证码数字
- headers:
  - ```cookie``` : 第一次请求得到的cookie
  - ```User-Agent``` : 自定义
- body:
  - ```captchaCode``` : 我们[识别](#识别验证码)到的验证码数字

#### 接收
- headers: 不重要
- body:
  - ```可能是true或false, 可以通过判断长度来识别, 一个4字母一个5字母, true就是验证码正确了```

## 第四次请求
目的 : 正式的登陆, 获取到可以访问任何页面的cookie
#### 发送
- 方法: POST
- url: http://jwzx.hrbust.edu.cn/academic/j_acegi_security_check
- headers:
  - ```cookie``` : 第一次请求得到的cookie
  - ```User-Agent``` : 自定义
- body:
  - ```j_username``` : 学号
  - ```j_password``` : 密码
  - ```j_captcha```  : 正确的验证码

#### 接收
- headers:
  - ```Set-cookie``` : 这个就是我们要的cookie
- body : 不重要

> **以上登陆完成, 可以访问正式页面.**
>
> 登陆所需代码在[sign_in](sign_in/sign_in_ug.go)文件夹中
> 
> 附:cookie我是使用了cookieJar来控制的, 这个是自动处理cookie的, 就比如收到了请求头里面有set-cookie, 那么你下次请求时, 就会自带这个cookie.

## 识别验证码
我们在登陆的时候, 获取到了验证码图片, 那么我们要识别验证码图片, 以便登陆.

本科生的识别, 我是使用了云码的, 这个平台是机器学习的, 可以识别本科生系统这种复杂奇怪的验证码图片
[云码数字识别](https://www.jfbym.com/test/52.html)
> **识别完验证码, 就可以获取正式页面了**
>
> 云码识别验证码代码在[identimage](../../identimage/jfbym.go)文件夹中

## 访问正式页面
1. 访问首页会得到访问左侧选择拦的地址
2. 访问左侧选择栏会得到课表的地址的随机串
3. 用课表地址加上随机串就能访问到课表页面了
因为左侧本学期课表的链接是带有随机串的, 所以通过需要提取出首页的信息, 然后再把访问课表页面和随机串拼接, 就能访问到课表页面了.

与研究生不同, 本科是获取一个页面, 然后解析出本学期的所有课表.
> **获取到html后, 如果我们需要将它保存到数据库中, 以便利用这些数据, 那么下面要开始解析获取到的html.**
>
> 获取html所需代码在[htmlgetter](htmlgetter/htmlgetter_ug.go)文件夹中


## 解析html
我写了两种解析方式
1. 给定周数, 可以得到那一周的课程表
2. 不给定周数, 可以得到本学期的每一周的所有课程表
> **解析完html, 只需要连接数据库, 将数据发送给数据库就可以了.**
>
> 研究生解析所需代码在[parse_form](parse_form/parse_form_ug.go)文件夹中
