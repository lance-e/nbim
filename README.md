# nbim

### 亮点:

to be continue...

### 目录结构：

项目目录结构遵循 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

~~~
/cmd
/internal
/configs
/pkg
/examples
/docs

~~~

### 技术栈

- mysql
- redis
- grpc
- protobuf
- zap
- gorm

### 基本架构

- Business:业务服务器，可根据业务需求拓展
- Connection:长连接接入层，拆分出ip config ,gateway,state,用于维护与客户端的长连接，消息分发等
- Logic:逻辑层，负责设备信息，好友信息，群组消息管理，消息转发等逻辑

### 快速开始
##### 建表:

~~~shell
mysql -u lance -p nbim < ./sql/tables.sql
~~~

##### 构建：

~~~shell
make all
~~~

