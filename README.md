一个集成通用包的gin框架小例子

## Uses
 - [gin](https://github.com/gin-gonic/gin)
 - [gorm](https://github.com/go-gorm/gorm)
 - [jwt](https://github.com/golang-jwt/jwt)
 - [zap](https://github.com/uber-go/zap)
 - [viper](https://github.com/spf13/viper)
 - [redis](https://github.com/redis/go-redis)
 - [sonic](https://github.com/bytedance/sonic)
 - [swagger](https://github.com/swaggo/gin-swagger)
 - [universal-translator](https://github.com/go-playground/universal-translator)
 - [file-rotatelogs](https://github.com/lestrrat-go/file-rotatelogs)


## 目录结构
```
├── api               (api层)
├── config            (配置包)
├── dao               (dao层)
├── docs              (swagger文档目录)
├── global            (全局对象)
├── initialize        (初始化)
├── log               (日志)
├── middleware        (中间件层)
├── model             (模型层)
│   ├── request       (入参结构体)
│   └── response      (出参结构体)
├── route             (路由层)
├── service           (service层)
└── utils             (工具包)
```

## 参与者
[![Contributors](https://contributors-img.web.app/image?repo=Gekkoou/gin-example)](https://github.com/Gekkoou/gin-example/graphs/contributors)
