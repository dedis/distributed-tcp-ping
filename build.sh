pwd
go mod vendor
go mod tidy
go build -v -o ./dummy/bin/dummy ./dummy/
go build -v -o ./stats/bin/stats ./stats/

