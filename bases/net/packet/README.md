# Packet

`packet` 包提供了用于处理大端序数据包的读写功能。它包含了 `BigEndianPacket` 结构体及其相关方法，用于读取和写入各种数据类型。

## 安装

确保你已经安装了 Go 语言环境，并且在你的项目中初始化了 `go.mod` 文件。

```sh
go get github.com/orbit-w/meteor/bases/net/packet
```

## 使用方法

### 导入包

```go
import "github.com/orbit-w/meteor/bases/net/packet"
```

### 创建 Packet

```go
// 创建一个新的空包
p := packet.New()

// 创建一个指定初始大小的包
p := packet.NewWithInitialSize(1024)

// 从现有数据创建一个包
data := []byte{0x01, 0x02, 0x03}
p := packet.Reader(data)
```

### 写入数据

```go
p.WriteBool(true)
p.WriteByte(0x01)
p.WriteInt16(123)
p.WriteString("hello")
```

### 读取数据

```go
b, err := p.ReadBool()
if err != nil {
// 处理错误
}

i16, err := p.ReadInt16()
if err != nil {
// 处理错误
}

str, err := p.ReadString()
if err != nil {
// 处理错误
}
```

### 重置和返回 Packet

```go
p.Reset()
packet.Return(p)
```

## 方法

### BigEndianPacket 结构体

- `ReadBool() (bool, error)`
- `ReadByte() (byte, error)`
- `ReadInt8() (int8, error)`
- `ReadInt16() (int16, error)`
- `ReadInt32() (int32, error)`
- `ReadInt64() (int64, error)`
- `ReadUint16() (uint16, error)`
- `ReadUint32() (uint32, error)`
- `ReadUint64() (uint64, error)`
- `ReadBytes() ([]byte, error)`
- `ReadBytes32() ([]byte, error)`
- `WriteBool(bool)`
- `WriteByte(byte)`
- `WriteInt8(int8)`
- `WriteInt16(int16)`
- `WriteInt32(int32)`
- `WriteInt64(int64)`
- `WriteUint8(uint8)`
- `WriteUint16(uint16)`
- `WriteUint32(uint32)`
- `WriteUint64(uint64)`
- `WriteBytes([]byte)`
- `WriteBytes32([]byte)`
- `WriteString(string)`

## 贡献

欢迎提交问题和贡献代码。请确保在提交之前运行所有测试。

## 许可证

此项目基于 MIT 许可证开源。详细信息请参阅 LICENSE 文件。