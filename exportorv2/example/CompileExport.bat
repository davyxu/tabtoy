set TOOL_DIR=%cd%
cd ..\..\..\..\..\..
set GOPATH=%cd%
go install github.com/davyxu/tabtoy

cd %TOOL_DIR%

call Export.bat