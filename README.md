## Qiwi service

### Requirements

* Golang 1.11+
* MongoDB 4.0+

### Environment variables
```bash
PORT=50051
DEBUG=false
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DB=qiwi
```

### Start project 
```bash
go mod vendor # install dependencies
go run *.go
```

### Protobuf code generation

##### Requirements
* protobuf
* protoc-gen-go

##### Generate
```bash
protoc -I protobuf protobuf/qiwi.proto --go_out=plugins=grpc:protobuf
``` 
