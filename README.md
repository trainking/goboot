# goboot

**goboot**是个一个**Golang**开发过程中，使用的一些经验总结的库，涵盖`Web`开发和游戏开发领域。

## 目录

* `bin`: 编译后程序目录
* `configs`: 配置文件存放目录
* `internal`: 具体项目的逻辑存放目录
* `pkg`: 通用包目录
  - `errgroup`：一个携程编组实现
  - `httpapi`: 基于Echo框架实现的http api实现
  - `jwt`： jwt的实现
  - `log`: 日志库的实现
  - `redisx`: 基于redis6.0的工具实现
  - `utils`: 一些工具函数的实现