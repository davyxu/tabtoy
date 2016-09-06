md tool
"protoc.exe" --plugin=protoc-gen-go=protoc-gen-go.exe --go_out .\tool --proto_path "." tool.proto
"protoc.exe" --plugin=protoc-gen-go=protoc-gen-go.exe --go_out ..\test\test --proto_path=..\test ..\test\test.proto
"protoc.exe" --plugin=protoc-gen-go=protoc-gen-go.exe --go_out ..\testv2\test --proto_path=..\testv2 ..\testv2\test.proto