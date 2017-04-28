..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=exportorv2 ^
--csharp_out=.\Config.cs ^
--binary_out=.\Config.bin ^
--proto_out=.\Config.proto ^
--json_out=.\Config.json ^
--lua_out=.\Config.lua ^
--go_out=.\Config.go ^
--type_out=.\Type.json ^
--combinename=Config ^
--luaenumintvalue=true ^
--lan=zh_cn ^
Globals.xlsx ^
Sample.xlsx ^
Vertical.xlsx ^
Info.xlsx

@IF %ERRORLEVEL% NEQ 0 pause

: 表索引
copy .\Config.go .\table\table_gen.go

: ..\..\..\..\..\..\bin\tabtoy.exe --mode=exportorv2 --json_out=.\Config.json --combinename=Config Globals.xlsx Sample.xlsx Vertical.xlsx Info.xlsx