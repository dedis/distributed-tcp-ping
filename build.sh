pwd
/usr/local/go/bin/go mod vendor
/usr/local/go/bin/go mod tidy
/usr/local/go/bin/go build -v -o ./dummy/bin/dummy ./dummy/src/
/usr/local/go/bin/go build -v -o ./stats/bin/stats ./stats/