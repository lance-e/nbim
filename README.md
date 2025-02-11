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

### 快速开始
##### 建表:

~~~shell
mysql -u lance -p nbim < ./sql/tables.sql
~~~

##### 构建：

~~~shell
make all
~~~

