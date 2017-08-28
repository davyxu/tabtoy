..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=exportorv2 ^
--go_out=.\types_gen.go ^
--combinename=Builtin ^
--lan=zh_cn ^
BuiltinTypes.xlsx

@IF %ERRORLEVEL% NEQ 0 pause
