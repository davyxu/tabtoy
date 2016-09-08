md tool
"protoc.exe" --plugin=protoc-gen-go=protoc-gen-go.exe --go_out .\tool --proto_path "." tool.proto
"protoc.exe" --plugin=protoc-gen-go=protoc-gen-go.exe --go_out ..\exportorv1\test\test --proto_path=..\exportorv2\test ..\exportorv1\test\test.proto
"protoc.exe" --plugin=protoc-gen-go=protoc-gen-go.exe --go_out ..\exportorv2\test\test --proto_path=..\exportorv1\test ..\exportorv2\test\test.proto