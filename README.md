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