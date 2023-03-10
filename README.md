### Setup
Add GOPATH to bash profile
```bash
export GOPATH=$(go env GOPATH)
```

### Generate go code from grpc
```bash
./protoc -I=. --go-grpc_out=./pkg/manabuf/usermgmt --go_out=./pkg/manabuf/usermgmt ./proto/usermgmt/user.proto
```

### Build hexagon binary
```bash
go build -o $GOPATH/bin ./hexagen/hexagen.go
```

### Generate code
```bash
go generate ./...
```