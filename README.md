
```bash
./protoc -I=. --go-grpc_out=./pkg/manabuf/usermgmt --go_out=./pkg/manabuf/usermgmt ./proto/usermgmt/user.proto
```

```bash
go build -o $GOPATH/bin ./hexagen/hexagen.go
```

```bash
go generate ./...
```