set CURRDIR=%cd%
cd ../../../../../../..
set GOPATH=%cd%
cd %CURRDIR%
go run main.go
@IF %ERRORLEVEL% NEQ 0 pause