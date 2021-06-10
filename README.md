# go-utils
go 的工具库、辅助函数等

## goals
* [ ] 日志
* [ ] Redis
* [ ] MySQL
* [ ] helper - func 辅助函数
* [ ] json
* [ ] error
* [ ] snowflake - 雪花算法，生成 id

## 极客时间-毛剑 go 课程

### 日志&指标&链路跟踪

* https://microservices.io/

MTTR

MTBF

产品上线后，定期做“服务巡检”

服务降级时，日志打印可以使用 warning 

// var _ log.Logger = (*Logger)(nil)  写法的作用是什么？

