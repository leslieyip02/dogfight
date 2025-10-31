# Getting Started


## Generate Code
```
protoc --proto_path=protos --go_out=server protos/*.proto
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
