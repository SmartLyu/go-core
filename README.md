# go-core

本库中封装了常用功能，方便快速开发。推荐go项目引用本项目，快速实现安全的crud，可以根据需求保存审计信息等。

## 功能

### 数据库
  - 目前支持mysql，未来会支持其他数据库
  - 支持表不存在时自动创建表。
  - 支持数据记录，支持数据加密，支持数据校验。
  - 支持数据权限控制，支持数据审计。
  - 支持数据导出，支持数据导入。
  - 支持数据分页，支持数据排序，支持数据搜索。
  - 支持操作预期和事务控制，保障数据安全。
  - 支持数据缓存，支持数据缓存清理。
    - 如果redis初始化没有执行，将不能进行缓存。
    - 缓存数据会自动过期。
    - 缓存支持压缩存储，压缩可以节约空间和，提高数据传输效率，但会降低响应速度。

### api接口
  - 支持jwt认证，并可以刷新认证信息。
  - 支持数据传输标准化，包括输入输出的格式和数据校验。
  - 支持微服务架构的安全和鉴权。
  - 支持快速和数据库建立crud接口群。

### 日志
  - 统一所有日志输出方式，支持输出到文件，支持输出到控制台。
  - 支持日志文件自动切割，滚动迭代，按日期归档。
  - 数据库、api日志都可以定义和配置，包括debug级别的日志。

### 任务
  - 为保证任务安全，单次任务不要超过30分钟，如果有超长时间任务，建议切割实现
  - 支持任务自定义，支持任务串行、并行，支持定时任务，支持消息队列，支持分布式任务。
  - 支持任务日志，有全链路任务监控，支持任务执行耗时统计。
  - 考虑了分布式安全，分布式锁，避免任务重复执行。
  - 对缓存数据制定了基本的回收规则，保障资源利用最大化。
  - 支持任务重试，任务失败可以重试，支持任务限流，避免任务过载。
  - 支持任务告警，可以实时监控任务执行情况。
  - 支持分布式下的变量传递，并完成资源回收。
  - 为保证安全，任务会优雅退出，如果有幂等任务，不需要等待，则可以在声明中忽略。

### 配置
  - 支持yaml和json两种格式，支持环境变量覆盖。
  - 在使用中引用 `BasicConfig` 类的时候需要添加yaml参数，可仿照 `pkg/config` 进行添加。

## 使用方法

### 直接运行Demo

本地运行该项目代码，需要先初始化，启动后，即可访问

```shell
make init   # 初始化本地环境，下载依赖包
make run    # 运行任务
```

### 作为库引用

可以仿照库中pkg目录使用该包，下载方法 `go get github.com/yuanJewel/go-core`

一个简单的web项目可以仿照如下：

```
├── asset               (go-bindata 自动生成)
├── docs                (swag init 自动生成)
├── Dockerfile          (构建镜像，可以参照本项目)
├── Makefile            (编译文件，可以参照本项目)
├── README.md
├── application.yml     (默认配置文件，不建议提交到git中)
├── go.mod              (go mod init 自动生成)
├── main.go             (入口文件，可以参照本项目)
├── pkg                 (集成项目，自定义自己web项目逻辑的位置，结构可以参照本项目)
│    ├── api            (接口类的实现)
│    ├── config         (配置类的实现)
│    └── db             (数据库类的实现)
└── views               (自定义需要打包进过项目的前端页面，目前只能应用于'/api/v1'下的路由使用)
    ├── 404.html
    └── index.html
```

### 环境变量

- `RECORD_DATA`: 是否开启数据记录
  - `true`: 开启数据记录(默认)
  - `false`: 关闭数据记录
- `CHECK_TABLE_EXISTS`: 是否检查表是否存在
  - `true`: 检查表是否存在，该功能需要数据库高级权限(默认)
  - `false`: 不检查表是否存在
- `LOGGER_RETAIN_NUMBER`: 日志文件保留个数，默认为3
- `LOGGER_FILE_SIZE`: 日志文件最大大小，单位为Mb，默认为50
- `LOGGER_ROOT_PATH`: 日志文件存放位置，默认为`./logs`
- `LOGGER_OUT_STYLE`: 日志输出方式
  - `file`: 输出到文件(默认)
  - `stdout`: 输出到控制台
- `LOGGER_ACCESS_OUT_STYLE`: api访问日志输出格式
  - `file`: 输出到文件(默认)
  - `stdout`: 输出到控制台
- `LOGGER_OUT_LEVEL`: 日志输出级别
  - `info`: 输出信息级别的日志(默认)
  - `warn`: 输出警告及以上级别日志
  - `error`: 输出错误及以上级别日志
  - `debug`: 输出调试及以上级别日志

### 配置基础格式

可以查看 `config/object.go` 代码查看必须的配置信息，部分配置有默认值可以不填写

```yaml
apiVersion: v1

server:
  port: 8080

auth:
  key:                  # jwt的密钥
  timeout: 600
  refresh: 300
  cryptoKey:            # 数据加密密钥
  cryptoPrefix:         # 数据加密前缀

redis:
  host: 
  db: 
  password:
  pool_size: 10
  timeout: 5
  expiration: 600
  retry_delay_ms: 100
  is_zip: true

db:
  # 目前只支持mysql
  driver: mysql
  host:
  port: 3306
  db:
  user:
  password:
  charset: utf8
  idle_connections: 2
  max_connections: 10

# 如果需要使用任务模块功能
task:
  tag:
  concurrency: 10
  worker: true
  redis:
    host:
    db:
    password: 
  mq:
    host: 
    username: 
    password: 
```

## 目录结构

### api

目录 `api` 中声明，引用iris框架，实现api接口，封装了认证基本接口，安全检查接口

### db

目录 `db` 中声明包括数据库目录

- `object` 数据库的数据结构
- `service` 数据库接口
- `mysql` 数据库具体方法
- `redis` 缓存方法和接口

### 日志

目录 `logger` 中声明，引用logrus库实现日志的基本功能

### 任务

目录 `task` 中声明，引用machinery库实现任务自定义，可以串行、并行，支持定时任务，支持消息队列，支持分布式任务。

任务中需要用到redis和rabbitmq，为防止产生冲突缓存key出现冲突，自定义相关参数配置，主服务使用的redis和任务系统的redis区分配置。

## 声明

引用或商用请注明出处，详细授权使用声明请查看: [LICENSE](https://github.com/yuanJewel/go-core/blob/main/LICENSE)

大版本间隔可能会存在不兼容情况，更新迭代，请自行评估。

如果有扩展需求或使用问题，请提交issue，或者邮件联系作者：`luyu151111@163.com`