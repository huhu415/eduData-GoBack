# eduData(大学教务处爬虫)
这个爬虫是通过学生学号和密码, 模仿真人来登陆到教务处系统, 获取并解析课程表到数据库, 然后增删改查到前端

### 为什么开源
因为现在的爬虫类小程序, 只要涉及到广告等盈利手段, 那么就是违法的.

只不过就是 法律意识淡泊 或 大学或商业组织 没有找到他, 铤而走险罢了.

我开源出来知识供大家**学习**和**使用**, 不允许用于**任何商业用途**.

### 技术栈
- 使用了golang作为后台语言
- postgresql作为数据库
- gorm来操作数据库
- gin作为web框架
- 微信小程序为前端显示

### 目前支持的大学

- [x] [哈尔滨理工大学](School/hrbust)
    - [x] 本科生
    - [x] 研究生
- [x] [东北农业大学](School/neau)
    - [x] 本科生
    - [ ] 研究生
- [ ] 黑龙江大学
    - [ ] 本科生
    - [ ] 研究生
- [ ] 东北林业大学
    - [ ] 本科生
    - [ ] 研究生
- [ ] 哈尔滨师范大学
    - [ ] 本科生
    - [ ] 研究生


# 公共部分

## 初始化数据库
要使用数据库, 要先初始化数据库, 定义好增删改查这些方法, 和新建数据库的方法.

我使用的是orm框架, 所以只需要定义好结构体, 然后调用orm的方法就可以了.

我使用的是PostpreSql, 据说这个现在比Mysql强
> **初始化数据库后, 可以通过web框架来处理逻辑, 处理过程中使用数据库.**
>
> 数据库代码在[database](database/database.go)文件夹中


## 发送数据到自己的数据库
我们要根据用户的请求来决定是发送数据到数据库, 还是从数据库中获取数据.

所以需要一个web框架, 来处理路由, 然后根据路由来处理逻辑.

我使用的是gin框架, 这个是golang中最强的web框架了, 用起来很方便.

> **不同的请求, 代表着不同的处理方式, 也就是api已经完成, 下一步就是前端调用与显示**
>
> 路由代码在[router](router/router.go)文件夹中;
> 处理路由后的逻辑在[app](app/app.go)文件夹中


## 前端
我用的是微信小程序.
> **用来展示**
>
> 这里是[前端代码](https://github.com/huhu415/eduData-WxFront)的项目地址

## 不同学校的差异文件夹
不同学校的差异, 是登陆的方式不同, 获取课表的不同, 解析课表的不同.

为了查找, 一个学校一个文件夹, 其中里面有两个个文件夹(本科研究生), 一个文件夹里有signin,GetData,ParseData,Test,md说明文档.

## 其余文件夹

- [setting](setting/setting.go) : 读取ini配置文件所需代码
- [middleware](middleware/jwt.go) : 中间件所需代码(jwt技术和日志记录)

-------------------

# 使用方法
1. 首先要安装golang, 安装方法自行百度, 国内需要配置go代理
   - ```go env -w GOPROXY=https://goproxy.cn,direct```
   - ```go env -w GOSUMDB=goproxy.cn/sumdb/sum.golang.org```
   - ```go env -w GO111MODULE=on```
2. 要在项目中安装依赖, 在项目根目录下执行```go mod tidy```
3. 安装postpresql
   - 创建一个数据库```CREATE DATABASE Courses;```
   - [config.ini](config/config.ini)中配置好数据库的连接信息
   - ```sudo vim /etc/postgresql/9.3/main/postgresql.conf```中修改```listen_addresses = '*'```
   - ```sudo vim /etc/postgresql/9.3/main/pg_hba.conf```中添加```ipv4 0.0.0.0/0```
4. 初始化表, 在[database](database/database.go)中有初始化表的方法, 是被注释掉的, 可以取消注释, 创建表就创建一次就行了, 用完取消注释
    - 如果要用代码创建表的话, 是要build出来后, 执行一次的
5. 查看一下[config.ini](config/config.ini)配置文件, 看看有没有需要修改的地方
6. 在项目根目录下执行```go build .```, 就编译出来可执行文件了
7. 执行```./eduData```就可以运行了
8. (可选)如果要部署的话, 用[systemctl命令](eduData.service)这个文件来管理
    - 把这个文件修改好参数, 然后放到```/etc/systemd/system/```目录下
    - 执行```systemctl start eduData```就可以启动了
    - 执行```systemctl stop eduData```就可以停止了
    - 执行```systemctl status eduData```就可以查看状态和日志
    - 执行```systemctl enable eduData```就可以开机自动启动了

### 部署注意事项
1. 检查一下[config.ini](config/config.ini)中的配置
2. 检查[database](database/database.go)中表的初始化那行注释与否
3. 检查[router](router/router.go)中的路由与中间件是否逻辑正确


## TODOs

- [ ] 设配不同的学校
- [ ] 通过grpc来分布式获取课表, 以防学校封ip
- [ ] 可以用k8s来部署分布式
- [ ] 日志方面, 要搞懂, 记得清楚
- [ ] 增加多租户功能
- [ ] 以后有钱了, 把百度orc下了, 因为百度的需要一个月更新一次凭证, 凭证坏了, 也不给什么提示, 就是500错误


## 开发提示
- 以下任何一个满足, 都代表着课程信息不完整, 都在课表里显示不出来, 只能显示在下面
  - course.NumberOfLessons == 0 
  - course.NumberOfLessonsLength == 0 
  - course.WeekDay == 0
  - course.week == 0