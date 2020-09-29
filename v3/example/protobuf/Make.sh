# protoc 下载
# wget https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protoc-3.13.0-win64.zip
go install google.golang.org/protobuf/cmd/protoc-gen-go
./protoc --go_out=. ./table.proto