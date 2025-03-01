# nbim

### 亮点:

- 自实现tcp网络库，支持海量并发，可根据业务场景定制化
- 拆分ip config server,动静结合策略,动态负载均衡,智能调度
- 自定义消息协议，固定消息头有效处理tcp粘包拆包
- 雪花算法生成分布式id
- 使用时间轮算法优化state server中的大量定时任务(消息，心跳，重连)的资源占用

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

