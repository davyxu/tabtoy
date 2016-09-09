..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=exportorv2 ^
--pbt_outdir=. ^
--proto3_outdir=. ^
--proto2_outdir=. ^
--json_outdir=. ^
--lua_outdir=. ^
--csharp_outdir=. ^
--binary_outdir=. ^
Sample.xlsx

@IF %ERRORLEVEL% NEQ 0 pause

..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=exportorv2 ^
--csharp_outdir=. ^
Exp.xlsx


@IF %ERRORLEVEL% NEQ 0 pause