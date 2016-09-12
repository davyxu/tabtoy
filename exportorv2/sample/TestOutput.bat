..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=exportorv2 ^
--csharp_outdir=. ^
--binary_outdir=. ^
--pbt_outdir=. ^
--proto_outdir=. ^
--json_outdir=. ^
--lua_outdir=. ^
--combinename=Config ^
Sample.xlsx ^
Exp.xlsx

@IF %ERRORLEVEL% NEQ 0 pause



