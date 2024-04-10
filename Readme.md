# eduData(大学教务处爬虫)
这个爬虫是通过学生学号和密码, 模仿真人来登陆到教务处系统, 获取并解析课程表到数据库, 然后增删改查到前端

### 为什么开源
❗❗️❗️因为现在的爬虫类小程序, 只要涉及到广告等盈利手段, 那么就是违法的.

❗️❗️❗️只不过就是 __法律意识淡薄__ 或 __大学和商业组织__ 没有找到他, **铤而走险**.

❗️❗️❗️本仓库的所有内容仅供**学习**和**参考**之用, 禁止用于**任何商业用途**.

### 技术栈
- 使用了golang作为后台语言
- postgresql作为数据库
- gorm来操作数据库
- gin作为web框架
- [微信小程序](https://github.com/huhu415/eduData-WxFront)为前端显示

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

## TODOs

- [ ] 设配不同的学校
- [ ] 通过grpc来分布式获取课表, 以防学校封ip
- [ ] 可以用k8s来部署分布式
- [ ] 日志方面, 要搞懂, 记得清楚
- [ ] 增加多租户功能
- [ ] 以后有钱了, 把百度orc下了, 因为百度的需要一个月更新一次凭证, 凭证坏了, 也不给什么提示, 就是500错误

----------------


# Get Start
## Self Compiled
1. 安装golang, 安装方法自行百度, 国内需要配置go代理
   - ```go env -w GOPROXY=https://goproxy.cn,direct```
   - ```go env -w GOSUMDB=goproxy.cn/sumdb/sum.golang.org```
   - ```go env -w GO111MODULE=on```
2. 要在项目中安装依赖, 在项目根目录下执行```go mod tidy```
3. 安装postpresql
   - 创建一个数据库```CREATE DATABASE Courses;```
   - **环境变量**或**配置文件**中配置好数据库的连接信息
   - ```sudo vim /etc/postgresql/9.3/main/postgresql.conf```中修改```listen_addresses = '*'```
   - ```sudo vim /etc/postgresql/9.3/main/pg_hba.conf```中添加```ipv4 0.0.0.0/0```
4. 在项目根目录下执行```go build .```, 就编译出来可执行文件了
5. 执行```./eduData```就可以运行了
6. (可选)如果要部署的话, 用[systemctl命令](eduData.service)这个文件来管理
    - 把这个文件修改好参数, 然后放到```/etc/systemd/system/```目录下
    - 执行```systemctl start eduData```就可以启动了
    - 执行```systemctl stop eduData```就可以停止了
    - 执行```systemctl status eduData```就可以查看状态和日志
    - 执行```systemctl enable eduData```就可以开机自动启动了


## docker
### Use Env
docker-compose.yml
```yaml
version: '3'
services:
  edudata:
    image: registry.cn-wulanchabu.aliyuncs.com/zzyan/back-go
    environment:
      - TZ=Asia/Shanghai
      - EDU_PG_CONFIG= xxxx
      - EDU_JFYM_TOKEN= xxxx
      - EDU_ACCESSTOKEN= xxx
```

### Use Config File
config.yaml
``` yaml
pg_config : host=localhost user=postgre password=123 dbname=123 port=5432 sslmode=disable TimeZone=Asia/Shanghai
listen_port : 8080
jwt_key : 9385g0x98n347tx980y****s****
UserAgent : Mozilla/5.0 (Macintosh Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36

# https://www.jfbym.com/test/52.html 云码,地址
jfym_request_url : http://www.jfbym.com/api/YmServer/customApi
jfym_token : XNcI2JjmD882dHRqDiFMYibe*******

# BaiduAccessToken要30天更新一次，否则会报错.
# https://console.bce.baidu.com/tools/#/index 更新token地址
baidu_request_url : https://aip.baidubce.com/rest/2.0/ocr/v1/numbers
baidu_accessToken : 24.aa2194d8df42835d0763*****852ba5.*****0******
```

docker-compose.yml
```yaml
version: '3'
services:
  edudata:
    image: registry.cn-wulanchabu.aliyuncs.com/zzyan/back-go
    volumes:
      - /home/edudata/config.yaml:/config.yaml
    environment:
      - TZ=Asia/Shanghai
```

----------------

## 开发提示
- 以下任何一个满足, 都代表着课程信息不完整, 都在课表里显示不出来, 只能显示在下面
  - course.NumberOfLessons == 0 
  - course.NumberOfLessonsLength == 0 
  - course.WeekDay == 0
  - course.week == 0