..\..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=v2 ^
--json_out=CombineConfig.json ^
--combinename=Config ^
--lan=zh_cn ^
Item.xlsx+Item_Equip.xlsx+Item_Pet.xlsx

@IF %ERRORLEVEL% NEQ 0 pause