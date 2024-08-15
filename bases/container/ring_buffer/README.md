# RingBuffer

`ring_buffer` 包提供了一个通用的环形缓冲区实现。它支持基本的插入、删除和查看操作，并且可以根据需要动态调整大小。

## 安装

确保你已经安装了 Go 语言环境，并且在你的项目中初始化了 `go.mod` 文件。

```sh
go get github.com/orbit-w/meteor/bases/container/ring_buffer
```

## 使用方法

### 导入包

```go
import "github.com/orbit-w/meteor/bases/container/ring_buffer"
```

### 创建 RingBuffer

```go
// 创建一个初始大小为 8 的 RingBuffer
rb := ring_buffer.New[int](8)
```

### 插入元素

```go
rb.Push(10)
rb.Push(20)
```

### 弹出元素

```go
item, ok := rb.Pop()
if ok {
    fmt.Println("Popped item:", item)
} else {
    fmt.Println("RingBuffer is empty")
}
```

### 查看队头元素

```go
item := rb.Peek()
fmt.Println("Head item:", item)
```

### 检查是否为空

```go
isEmpty := rb.IsEmpty()
fmt.Println("Is empty:", isEmpty)
```

### 获取长度

```go
length := rb.Length()
fmt.Println("Length:", length)
```

### 重置 RingBuffer

```go
rb.Reset()
```

### 收缩 RingBuffer

```go
contracted := rb.Contract()
fmt.Println("Contracted:", contracted)
```

## 方法

### RingBuffer 结构体

- `Push(item V)`
- `Pop() (V, bool)`
- `Peek() V`
- `Length() int`
- `IsEmpty() bool`
- `Reset()`
- `Contract() bool`

## 贡献

欢迎提交问题和贡献代码。请确保在提交之前运行所有测试。

## 许可证

此项目基于 MIT 许可证开源。详细信息请参阅 LICENSE 文件。