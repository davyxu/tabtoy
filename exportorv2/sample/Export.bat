..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=exportorv2 ^
--csharp_out=.\Config.cs ^
--binary_out=.\Config.bin ^
--pbt_out=.\Config.pbt ^
--proto_out=.\Config.proto ^
--json_out=.\Config.json ^
--lua_out=.\Config.lua ^
--go_out=.\Config.go ^
--combinename=Config ^
--lan=zh_cn ^
Globals.xlsx ^
Sample.xlsx ^
Info.xlsx

@IF %ERRORLEVEL% NEQ 0 pause
