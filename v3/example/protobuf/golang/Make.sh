# protoc 下载 https://github.com/protocolbuffers/protobuf/releases
go install google.golang.org/protobuf/cmd/protoc-gen-go
./protoc --go_out=. ../table.proto -I ../
go run main.go table.pb.go