# eduData(大学教务处爬虫)

用户输入学号, 密码后, 通过js逆向来模仿真人, 来登陆到教务系统, 获取并解析课程表到数据库, 就可以通过后端api增删改查后发送给前端.

### Why open source ?
❗❗️❗️现在从教务处查课表的爬虫小程序, 只要涉及到广告等盈利手段, 那么就是违法的.

❗️❗️❗️只不过就是 __法律意识淡薄__ 或 __大学和商业组织__ 没有找到他, **铤而走险**.

❗️❗️❗️本仓库的所有内容仅供**学习**和**参考**之用, 禁止用于**任何商业用途**.

### Tech
- 语言: Go
- 数据库: Postgresql
- ORM: Gorm
- Web框架: Gin
- 前端: [微信小程序](https://github.com/huhu415/eduData-WxFront)

### Supported universities
- [x] [hrbust_哈尔滨理工大学](school/hrbust)
    - [x] 本科生
    - [x] 研究生
- [x] [neau_东北农业大学](school/neau)
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

### TODOs
- [x] 自定义添加课程
- [x] 日志记录的不全面
- [ ] 自定义删除课程
- [ ] 多适配不同的学校
- [ ] 通过grpc来分布式获取课表, 以防学校封ip
- [ ] 增加多租户功能
- [ ] 以后有钱了, 把百度orc下了, 因为百度的需要一个月更新一次凭证, 到期后也没有提示, 就是500错误


# Get Start
## Traditional
#### Compiled
1. 安装golang, 安装方法自行百度, 国内需要配置go代理
   - ```go env -w GOPROXY=https://goproxy.cn,direct```
   - ```go env -w GOSUMDB=goproxy.cn/sumdb/sum.golang.org```
   - ```go env -w GO111MODULE=on```
2. 要在项目中安装依赖, 在项目根目录下执行```go mod tidy```
3. 安装postpresql
   - 创建一个数据库```CREATE DATABASE Courses;```
   - (如果跨域访问)```sudo vim /etc/postgresql/9.3/main/postgresql.conf```中修改```listen_addresses = '*'```
   - (如果跨域访问)```sudo vim /etc/postgresql/9.3/main/pg_hba.conf```中添加```ipv4 0.0.0.0/0```
4. **环境变量**或**配置文件**中配置好数据库信息等各种参数
5. 在项目根目录下执行```go build .```, 就编译出来可执行文件了
6. 执行```./eduData```就可以运行了

#### Deploy
eduData.service
```shell
# Put the eduData.service file into /etc/systemd/system/
[Service]
Type=simple
ExecStart=/root/eduData-GoBack/eduData
WorkingDirectory=/root/eduData-GoBack
Restart=always

[Install]
WantedBy=multi-user.target
```
Command
```shell
systemctl start eduData
systemctl stop eduData
systemctl status eduData
systemctl enable eduDat
```

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
services:
  edudata:
    image: registry.cn-wulanchabu.aliyuncs.com/zzyan/back-go
    volumes:
      - /home/edudata/config.yaml:/config.yaml
    environment:
      - TZ=Asia/Shanghai
```


## Develop Tips
- 以下任何一个满足, 都代表着课程信息不完整, 都在课表里显示不出来, 只能显示在下面
  - course.NumberOfLessons == 0
  - course.NumberOfLessonsLength == 0
  - course.WeekDay == 0
  - course.week == 0
