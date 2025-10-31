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

### Start Server
```
cd server
go run main.go
```
