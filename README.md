一个集成通用包的gin框架小例子

## Uses

| 组件名                  | 介绍       | 链接                                                    |
|----------------------|----------|-------------------------------------------------------|
| gin                  | Web框架    | https://github.com/gin-gonic/gin                      |
| gorm                 | ORM库     | https://github.com/go-gorm/gorm                       |
| dtm                  | 分布式事务服务  | https://github.com/dtm-labs/dtm                       |
| jwt                  | JWT库     | https://github.com/golang-jwt/jwt                     |
| zap                  | 日志库      | https://github.com/uber-go/zap                        |
| viper                | 配置管理     | https://github.com/spf13/viper                        |
| go-redis             | Redis库   | https://github.com/redis/go-redi                      |
| kafka-go             | Kafka库   | https://github.com/segmentio/kafka-go                 |
| sonic                | JSON库    | https://github.com/bytedance/sonic                    |
| gin-swagger          | Swagger库 | https://github.com/swaggo/gin-swagger                 |
| universal-translator | i18n库    | https://github.com/go-playground/universal-translator |
| file-rotatelogs      | 日志分割库    | https://github.com/lestrrat-go/file-rotatelogs        |

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
├── queue             (异步队列层)
│   └── drive         (队列驱动层)
├── route             (路由层)
├── service           (service层)
└── utils             (工具包)
```

## 参与者

[![Contributors](https://contributors-img.web.app/image?repo=Gekkoou/gin-example)](https://github.com/Gekkoou/gin-example/graphs/contributors)
