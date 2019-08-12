## Reserve gRPC

其他语言：
### **[English](README.md)**

这是我在学习GRPC和Go中编写的微服务的试验项目，它可能不适合直接移植到生产环境使用。 
我最初从[Alan Shreve的gRPC缓存服务](https://about.sourcegraph.com/go/grpc-in-production-alan-shreve)中获取了代码, 然后对重试和超时功能进行了一些更改。

### 安装和运行

#### 安装

```
go get github.com/jfeng45/reservegrpc
```

运行服务器
```
cd reserveserver
go run cacheJinServer.go
```
运行客户端
```
cd reserveclient
go run cacheJinClient.go
```
### 授权

[MIT](LICENSE.txt) 授权



