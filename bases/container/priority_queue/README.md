# PriorityQueue

`PriorityQueue` 是一个基于堆的数据结构，支持优先级队列的基本操作。它允许你根据优先级对元素进行插入、删除和更新操作。

## 安装

确保你已经安装了 Go 语言环境，并且在你的项目中初始化了 `go.mod` 文件。

```sh
go get github.com/orbit-w/meteor/bases/container/priority_queue
```

## 使用方法

### 导入包

```go
import (
"github.com/orbit-w/meteor/bases/container/priority_queue"
"github.com/orbit-w/meteor/bases/misc/common"
)
```

### 创建 PriorityQueue

```go
pq := priority_queue.New[string, int, common.Integer]()
```

### 插入元素

```go
pq.Push("key1", 100, 1)
pq.Push("key2", 200, 2)
```

### 检查元素是否存在

```go
exists := pq.Exist("key1")
```

### 获取元素

```go
item, found := pq.Get("key1")
```

### 更新元素

```go
pq.Update("key1", 150, 3)
```

### 更新元素优先级

```go
pq.UpdatePriority("key1", 4)
```

### 删除元素

```go
pq.Delete("key1")
```

### 弹出元素

```go
key, value, exists := pq.Pop()
```

### 按优先级弹出元素

```go
pq.PopByScore(5, func(k string, v int) bool {
fmt.Println("Popped:", k, v)
return true
})
```

### 检查队列是否为空

```go
isEmpty := pq.Empty()
```

## 贡献

欢迎提交问题和贡献代码。请确保在提交之前运行所有测试。

## 许可证

此项目基于 MIT 许可证开源。详细信息请参阅 LICENSE 文件。