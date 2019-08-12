# Reserve gRPC

Other language: 
### **[中文](README.zh.md)**

This is the project for me to learn writing Microservice in GRPC and Go, and it is not intended for production use. 
I originally took the code from [Alan Shreve's gRPC cache service](https://about.sourcegraph.com/go/grpc-in-production-alan-shreve)， then made some changes on retry and timeout features.

## Getting Started

### Installing

```
go get github.com/jfeng45/reservegrpc
```

Run Server
```
cd reserveserver
go run cacheJinServer.go
```
Run Client
```
cd reserveclient
go run cacheJinClient.go
```
## License

[MIT](LICENSE.txt) License



