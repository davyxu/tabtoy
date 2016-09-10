..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=exportorv2 ^
--pbt_outdir=. ^
--proto3_outdir=. ^
--json_outdir=. ^
--lua_outdir=. ^
--csharp_outdir=. ^
--binary_out=.\Config.bin ^
Sample.xlsx ^
Exp.xlsx

@IF %ERRORLEVEL% NEQ 0 pause