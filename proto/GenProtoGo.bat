md tool
"protoc.exe" --plugin=protoc-gen-go=protoc-gen-go.exe --go_out .\tool --proto_path "." tool.proto
@IF %ERRORLEVEL% NEQ 0 pause
"protoc.exe" --plugin=protoc-gen-go=protoc-gen-go.exe --go_out ..\exportorv1\test\test --proto_path=..\exportorv1\test ..\exportorv1\test\test.proto