..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=exportorv2 ^
--csharp_outdir=. ^
--binary_outdir=. ^
--combinename=Config ^
Sample.xlsx ^
Exp.xlsx

@IF %ERRORLEVEL% NEQ 0 pause

:--pbt_outdir=. ^
:--proto3_outdir=. ^
:--json_outdir=. ^
:--lua_outdir=. ^