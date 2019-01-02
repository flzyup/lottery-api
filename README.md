# 年会抽奖程序后端API服务
本API服务基于`golang`编写，原用于荣威俱乐部RX8车友会2018第一届年会抽奖程序

# Require
注意：由于使用了`golang`的`module`功能，所以需要

```go
go verson >= 1.11
```

很简单的API抽奖程序，`db-*.sql`一个是数据库表结构，另外一个是测试数据样本

跟[RX8抽奖前端](https://github.com/flzyup/rx8-lottery)配合一起使用

# Usage

- Import db structure and data

使用*.sql文件导入数据结构和测试数据

- Go compile
```go
go build i4o.xyz/rx8lottery
```

# License
Apache License 2.0

