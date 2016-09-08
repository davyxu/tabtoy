..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=exportorv2 ^
--pbt_outdir=. ^
--proto3_outdir=. ^
--proto2_outdir=. ^
Sample.xlsx

@IF %ERRORLEVEL% NEQ 0 pause