# go-core

## 使用方法

### 直接运行Demon

本地运行该项目代码，需要先初始化，启动后，即可访问

```shell
make init   # 初始化本地环境，下载依赖包
make run    # 运行任务
```

### 作为库引用

可以仿照库中pkg目录使用该包，下载方法 `go get github.com/SmartLyu/go-core`

一个简单的web项目可以仿照如下：

```
├── asset               (go-bindata 自动生成)
├── docs                (swag init 自动生成)
├── Dockerfile          (构建镜像，可以参照本项目)
├── Makefile            (编译文件，可以参照本项目)
├── README.md
├── VERSION             (记录版本号，配合编译文件构建制品)
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

## 目录结构

### api

目录 `api` 中声明，引用iris框架，实现api接口，封装了认证基本接口，安全检查接口

### db

包括数据库目录
- `db` 数据库的数据结构
- `cmdb` 数据库接口
- `mysql` 数据库具体方法

### 日志

目录 `logger` 中声明，引用logrus库实现日志的基本功能