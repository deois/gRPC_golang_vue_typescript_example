# gRPC, gRPC-Gateway, gRPC-Web, Vue3, Vite, Typescript Example

- Development environment based on Windows 11
- Implementation of gRPC server and client (50011)
- Implementation of gRPC-Gateway server (8080)
- Implementation of gRPC-Reflection server (50051)
- Implementation of gRPC-Web client
- Implementation of Vue3, Vite, Typescript frontend
- Integrated execution of gRPC server, gRPC-Gateway server, gRPC-Web client, and Vue3 frontend

## Operation Example

### golang gRPC server
- Implementation of SayHello RPC method in MyService
  - SayHello: When String is requested, responds with "Hello" + String
  - StreamTime: When requested, responds with current time at 1-second intervals

### golang gRPC client
- Calling SayHello RPC method in MyService
  - When requesting "World", responds with "Hello World"
  - When requesting "Gopher", responds with "Hello Gopher"

### vue3 frontend
- Request to gRPC server through gRPC-Web client from typescript frontend
  - When requesting sayHello "World", outputs "Hello World" to debug console
  - When requesting streamTime, outputs current time to debug console at 1-second intervals


## Installing protoc 
### using binary (windows 11) 

- Download win64 binary from https://github.com/protocolbuffers/protobuf/releases
- ex, https://github.com/protocolbuffers/protobuf/releases/download/v28.1/protoc-28.1-win64.zip
- Set PATH for protoc executable
- Verify protoc --version

#### Check version after running protoc
```shell
protoc --version
```

```log
libprotoc 28.1
```

## Writing proto code

```
mkdir -p backend/proto
```

```backend/proto/service.proto
syntax = "proto3";

package proto;

option go_package = "backend/proto";

import "google/api/annotations.proto";

service MyService {
  rpc SayHello (HelloRequest) returns (HelloResponse) {
    option (google.api.http) = {
      post: "/v1/sayhello"
      body: "*"
    };
  }
  rpc StreamTime (TimeRequest) returns (stream TimeResponse) {
    option (google.api.http) = {
      post: "/v1/streamtime"
      body: "*"
    };
  }
}

message TimeRequest {}

message TimeResponse {
  string current_time = 1;
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}
```

## golang gRPC server and golang gRPC client


### Preparation
```shell
cd backend
go mod init backend
go get -u google.golang.org/grpc
go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go get -u github.com/improbable-eng/grpc-web/go/grpcweb
git clone https://github.com/googleapis/googleapis.git
```

### Converting proto files to go files and typescript files
```shell
npm install @protobuf-ts/plugin
protoc --proto_path=. --proto_path=./googleapis --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative proto/service.proto
mkdir -P ./frontend/src/backend/gRPC
npx protoc --ts_out ../frontend/src/backend/gRPC --proto_path=proto --proto_path=./googleapis proto/service.proto
```

#### Running golang server

``` shell
go run server.go
```

### Running golang client

```shell
go run client.go
```

``` log
2024/09/12 16:35:58 Response: Hello World
2024/09/12 16:35:58 Response: Hello Gopher
Done
```

## frontend (vite vue typescript)

```shell
npm create vite@latest frontend -- --template vue-ts
cd frontend
npm install @protobuf-ts/plugin @protobuf-ts/runtime @protobuf-ts/runtime-rpc @protobuf-ts/grpcweb-transport
npm install @improbable-eng/grpc-web @improbable-eng/grpc-web-node-http-transport
```







