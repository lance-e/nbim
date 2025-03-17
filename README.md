# nbim

### 亮点

- 自实现tcp网络库，支持海量并发，可根据业务场景定制化
- 拆分ip config server,动静结合策略,动态负载均衡,智能调度客户端连接的目标网关机
- 自定义消息协议，固定消息头有效处理tcp粘包拆包，解析消息信令执行不同业务逻辑
- 雪花算法生成分布式id
- 使用时间轮算法优化state server中的大量定时任务(消息，心跳，重连)的资源占用
- 单聊写扩散提高消息及时性，群聊读扩散减少数据库写入压力，保证系统稳定性

### 基本架构：

消息协议：

![](./image/protocal)

接入层：gateway和state的拆分：

![](./image/gateway_state)



### 目录结构

项目目录结构遵循 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

~~~
/cmd     		启动入口
/internal		服务私有文件
/configs		配置
/pkg			复用包
/examples		示例
/docs			文档

~~~

### 技术栈

- mysql
- redis
- grpc
- protobuf
- zap
- gorm

### 基本架构

- Connection:接入层，分为tcp/ws长连接和http短链接两板块。
- Logic:逻辑层，负责设备信息，好友信息，群组消息管理，消息转发等逻辑。

### 快速开始
##### 建库建表:

~~~shell
mysql -u lance -p < ./sql/tables.sql
~~~

##### 构建：

~~~shell
make all
~~~

