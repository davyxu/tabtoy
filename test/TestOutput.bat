..\proto\protoc.exe test.proto --plugin=protoc-gen-meta=..\..\..\..\..\bin\protoc-gen-meta.exe --proto_path "." --meta_out test.pb:.
..\..\..\..\..\bin\tabtoy.exe --mode=xls2pbt --pb=test.pb --outdir=. --fmt=pbt Actor.xlsx
..\..\..\..\..\bin\tabtoy.exe --mode=xls2pbt --pb=test.pb --outdir=. --fmt=lua Actor.xlsx
..\..\..\..\..\bin\tabtoy.exe --mode=xls2pbt --pb=test.pb --outdir=. --fmt=json Actor.xlsx