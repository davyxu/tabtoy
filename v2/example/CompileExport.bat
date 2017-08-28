set TOOL_DIR=%cd%
cd ..\..\..\..\..\..
set GOPATH=%cd%
go build -v -o %GOPATH%\bin\tabtoy.exe github.com/davyxu/tabtoy

cd %TOOL_DIR%

call Export.bat