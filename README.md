# infrastructure

![](https://img.shields.io/badge/LICENSE-MIT-blue.svg) ![](https://goreportcard.com/badge/github.com/yeqown/infrastructure) [![](https://godoc.org/github.com/yeqown/infrastructure?status.svg)](http://godoc.org/github.com/yeqown/infrastructure)



Collecting some utilities those will be useful when coding a web application with Go.

## Features

* Framework based `gorm`, `logrus`, `go-redis/redis`, `gin`, `etcd`, `mgo.v2`, `jwt-go`.
* Health Checking
* Lang utils, like: `MultiSorter`, `ConvertStruct2Map`, `WalkFolder`...
* Some Types including: `Envrion`, `Database Cfg`, `Codes to response`...

## Todos: 

- [x] fill gormic package
- [x] finish logger test cases
- [x] finish code test cases
- [x] finish utils test cases
- [ ] finish framework test cases
- [x] finish gormic test cases
- [x] support `validator.v8`.`ResourceCheck`

## Golang Model Struct to Service Struct

*moved to [jademperor/go-tools](github.com/jademperor/go-tools)*

## Examples

* [gin resource checker](examples/gin-resource-checker)
> `gin` do not support validator.v9 for now ...

    resource-validator to check if the resource exists.

* [health checker](examples/health-checker)

    checking healthy of `Mongo`, `Redis`, `SQL-DB`, `Service over TCP`.

* [amqp wrapper](examples/amqp-wrapper)

    wrap `amqp.Connection` with reconnection ability.