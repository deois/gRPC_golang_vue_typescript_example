# gRPC, gRPC-Gateway, gRPC-Web, Vue3, Vite, Typescript Example
## golang echo/v4 적용
- 
- Windows 11 기반 개발 환경
- gRPC 서버 및 클라이언트 구현 (50011)
- gRPC-Gateway 서버 구현 (8080)
- gRPC-Reflection 서버 구현 (50051)
- gRPC-Web 클라이언트 구현 
- Vue3, Vite, Typescript 프론트엔드 구현
- gRPC 서버, gRPC-Gateway 서버, gRPC-Web 클라이언트, Vue3 프론트엔드를 통합하여 실행

## 동작 예시

### golang gRPC 서버
- MyService 내에 SayHello RPC 메서드 구현
  - SayHello : String 요청 시, "Hello" + String 응답
  - StreamTime : 요청 시, 현재 시간을 1초 간격으로 응답

### golang gRPC 클라이언트
- MyService 내에 SayHello RPC 메서드 호출
  - "World" 요청 시, "Hello World" 응답
  - "Gopher" 요청 시, "Hello Gopher" 응답

### vue3 프론트엔드
- typescript 프론트엔드에서 gRPC-Web 클라이언트를 통해 gRPC 서버에 요청 
  - sayHello "World" 요청 시, 디버그 콘솔에 "Hello World" 출력
  - streamTime 요청 시, 현재 시간을 1초 간격으로 디버그 콘솔에 출력


## protoc 설치 (windows 11)

- https://github.com/protocolbuffers/protobuf/releases 에서 win64 바이너리 다운로드
- ex, https://github.com/protocolbuffers/protobuf/releases/download/v28.1/protoc-28.1-win64.zip
- protoc 실행파일 PATH 설정
- protoc --version 확인

### protoc 실행 후 버전 확인
```shell
protoc --version
```

```log
libprotoc 28.1
```

## proto 코드 작성

mkdir -p backend/proto

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

## golang gRPC 서버 및 golang gRPC 클라이언트


### 준비
```shell
cd backend
go mod init backend
go get -u google.golang.org/grpc
go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go get -u github.com/improbable-eng/grpc-web/go/grpcweb
git clone https://github.com/googleapis/googleapis.git
```

### proto 파일을 go 파일, typescript 파일로 변환
```shell
npm install @protobuf-ts/plugin
protoc --proto_path=. --proto_path=./googleapis --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative proto/service.proto
mkdir -P ../frontend/src/backend/gRPC
npx protoc --ts_out ../frontend/src/backend/gRPC --proto_path=proto --proto_path=./googleapis proto/service.proto
```

#### golang 서버 실행

``` shell
go run server.go
```

### golang 클라이언트 실행

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







