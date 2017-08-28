: 输出C#源码,二进制(例子中供C#读取), lua表, json格式
: 适用于csharp, golang, lua例子
..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=v2 ^
--csharp_out=.\csharp\Example\Config.cs ^
--binary_out=.\csharp\Example\Config.bin ^
--lua_out=.\lua\Config.lua ^
--luaenumintvalue=true ^
--go_out=.\golang\table\table_gen.go ^
--json_out=.\golang\Config.json ^
--combinename=Config ^
--lan=zh_cn ^
Globals.xlsx ^
Sample.xlsx

@IF %ERRORLEVEL% NEQ 0 pause