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

### Protobuf documentation

See [Here](protobuf/readme.md)

*Powered by: https://github.com/pseudomuto/protoc-gen-doc*

##### Generate

```bash
docker run --rm \
  -v $(pwd)/protobuf:/out \
  -v $(pwd)/protobuf:/protos \
  pseudomuto/protoc-gen-doc --doc_opt=markdown,readme.md
```

### Protobuf code generation

##### Requirements
* protobuf
* protoc-gen-go

##### Generate
```bash
protoc -I protobuf protobuf/qiwi.proto --go_out=plugins=grpc:protobuf
``` 
