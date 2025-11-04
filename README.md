# Getting Started

## Generate Code
```
protoc \
    --plugin=client/node_modules/.bin/protoc-gen-ts_proto \
    --proto_path=protos \
    --ts_proto_out=client/src/pb \
    --go_out=server \
    protos/*.proto
```

## Start Client
```
cd client
npm run build
```

## Start Server
```
cd server

# load balancer
go run cmd/master/main.go

# game server
go run cmd/worker/main.go
```
