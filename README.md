# go-web

## 工程搭建
1. 下载模块
```shell
go mod download
```

2. 编译
```shell
go build cmd/goweb.go
```

3. 运行
```shell
./goweb
```

4. 访问
打开:http://localhost:8080/

## 工程结构
`Feature-based`工程结构。

# 任务状态转移图
![](./docs/uml/task-state-transition.png)
```xml

@startuml

[*] --> Created: Task created
Created --> Running: Start task
Running --> Done: Task completed
Running --> Cancelled: Task cancelled
Running --> Failure: Task failed
Created --> Cancelled: Task cancelled before started
Done --> [*]
Cancelled --> [*]
Failure --> [*]

@enduml

```

# worker管理
worker由worker manager管理并存储。Scheduler定时pingworker，worker manager根据ping结果更新worker状态。Worker状态unhealthy时，worker manager会删除worker。

# task 管理
task 由task manager管理，task manager根据worker状态分配task。scheduler从task manager获取task，并分配给worker（Pull model）。

# Task状态存储与管理
Task与Task状态由Task Store存储， Task Store V1版本使用内存存储，Task Store V2版本使用redis存储。Task使用本地存储以后，进行分布式扩展时会出现问题，所以目前只支持单机版本。
V2使用Redis cluster存储后无此问题。V2可以考虑将状态存储到非Redis的存储中，如MySQL。