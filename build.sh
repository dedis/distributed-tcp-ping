pwd
/usr/local/go/bin/go mod vendor
/usr/local/go/bin/go mod tidy
protoc --go_out=. --go-grpc_out=. dummy/src/messages.proto
/usr/local/go/bin/go build -v -o ./dummy/bin/dummy ./dummy/src/