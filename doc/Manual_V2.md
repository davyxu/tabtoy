# tabtoy

# 导出步骤

## 准备电子表格文件

格式请参考:
	
	[Sample.xlsx]: (https://github.com/davyxu/tabtoy/blob/master/exportorv2/sample/Sample.xlsx)
	
	
## 准备tabtoy二进制

* 已经编译好的二进制:
	
	[Tabtoy Release]: (https://github.com/davyxu/tabtoy/releases)
	
* 手动编译
	
	go get github.com/davyxu/tabtoy
	
## 编写导出命令行

范例:
		
	```bat
	
	tabtoy --mode=exportorv2 --json_out=.\config.json Table.xlsx
	
	```

# 电子表格格式

电子表格文件的表单(Sheet)分为两种格式

## 数据表单(DataSheet)

* 首行首列单元格非空的表单, 被识别为数据表单

* 表单(Sheet)名字前带有#时, 表单内容不会被导出

* 数据表单必须拥有4条信息行

从上到下分别是:

### 字段名 行

* 以_或英文开头的标示符,不能包含中文

### 字段类型 行

支持以下类型
* int32

* int64

* uint32

* uint64

* string

* bool

* float

* (枚举类型)

	在类型表单中定义, 参见类型表单(TypeSheet)

* (结构体类型)

	在类型表单中定义, 参见类型表单(TypeSheet)


数组方式的值, 请在以上类型前添加repeated

例如:

	repeated int32
	
	repeated ActorType
	

### 字段特性 行

以pbt文本描述字段的特有功能, 如

* 字段数据重复性检查

	格式: RepeatCheck: true
	
	功能: 对单元格以字符串方式检查是否重复
	
	范例: 要求ID字段不能重复时, 设置重复性检查
	
* 值切割

	格式: ListSpliter: "分隔符"
	
	功能: 对repeated字段的单元格内容按照分隔符进行切割, 导出后以数组方式存储
	
	范例: 奖励包id通常是repeated的, 在特性中添加ListSpliter: ";" 单元格填写: 100;200时, 获得的将是包含 100, 200的数组
	
	注意: 整形数值类在电子表格中的分隔符推荐使用分号";" 切忌使用逗号"," 因逗号为电子表格默认的大数分割符
	
		 不能使用值切割对结构体字段进行操作,多个结构体请标记repeated类型并放置在多个单元格
	
* 默认值

	格式: Default: "YourDefaultValue"
	
	功能: 对字段设置默认值, 单元格不填写时, 导出获取数据后以默认值获取
	
	范例: bool值默认为false, 如果需要默认值为true时, 修改默认值

* 索引创建

	格式: MakeIndex: true
	
	功能: 数据输出时, 对MakeIndex所在的字段创建索引
	
	范例: 对Item表的ID字段添加索引, 在代码中可以通过ID直接获取Item指定的记录
	
* 字段别名

	格式: Alias: "字段中文名"
	
	使用区域: 字符串解析为结构体, 类型表单中的枚举命名
	
	功能: 设置一个字段的别名, 别名通常是中文;当字段需要在单元格表示时, 可以使用别名填写
	
	范例: Alias: "血量"
	
* 自定义tag

	格式: 与系统tag不冲突的tag类型名: "你的tag值"
	
	使用区域: @Types中的字段Meta, 数据表中的Meta均可
	
	功能: 添加一个自定义tag名, 为字段添加一个特殊的语义或者定义, 方便逻辑中使用
	
	范例: RandomType: "Pick"
	

### 字段描述 行

此行不解析, 但请保留, 并以中文编写注释, 方便查看



## 类型表单(TypeSheet)

* 被命名为@Types的表单被识别为类型表单

* 整个电子表格文件只允许有1个

* 类型表单必须包含2条信息行

从上到下分别是:

### 文件特性 行

以pbt文本描述字段的特有功能, 如

* 包名

	格式: Package: "gamedef"

	功能: 指定输出的包名
	
	作用: 影响代码生成的命名控件或者包名, 如C#,golang
	
	范例: 以gamedef命名的包名, C#输出的代码包含namespace gamedef{}
	
* 表名

	格式: TableName: "Item"
	
	功能: 指定输出的表名
	
	作用: 影响数据输出时的记录数组字段名, 所有列所在结构体的名称为 表名+Define
	
	范例: 以Item命名的表名, 输出数据在Item的数组中获取所有记录

* 输出标记匹配

	格式: OutputTag: [".pbt", ".proto"]

	功能: 当表相关类型和数据需要输出时, 当输出标记匹配和输出格式匹配时进行输出
	
	作用: 对数据, 类型, 索引都有效, 一切和本表有关系的信息在设定输出且匹配时, 会被过滤
	
	范例: 所有的表一次导出时, 服务器匹配.pbt, .proto  客户端匹配.cs, .bin时, 类型数据, 索引和数据将自动根据客户端服务器需求分离输出


### 类型字段描述 行

此行不解析, 但请保留, 并以中文编写注释, 方便查看

### 对象类型(ObjectType) 列

* 以_或英文开头的标示符,不能包含中文

* 对象类型可以是枚举或结构体, 通过是否有枚举值自动区分

### 字段名(FieldName) 列

* 以_或英文开头的标示符,不能包含中文

* 字段名归属对象类型

### 字段类型(FieldType) 列

* 参考 数据表单(DataSheet)的类型

* 结构体字段不能再次包含结构体, 但可以是枚举

* 枚举字段类型必须是int32

### 枚举值(Value) 列

* 此处填写时, 表示对象类型为枚举

* 枚举首值必须为0

* 枚举值不能重复

### 注释(Comment) 列

* 注释将出现在代码生成的代码注释中

### 特性(Meta) 列

* pbt文本描述, 参考 数据表单(DataSheet)的字段特性

* 结构体字段设置别名时, 将在数据表单中可以使用别名字段
	
	
# 命令行参数

## 指定文件格式输出	

### Protobuf格式

* 格式: --proto_out=path/to/out.proto

* 功能: 生成所有表中的类型信息的proto格式文本

* 范例: 生成proto后, 需要再生成pbt格式, 通过protoc将proto编译为你使用的语言代码后, 读取pbt文件

### Protobuf文本格式

* 格式: --pbt_out=path/to/out.pbt

* 功能: 生成所有表中的类型信息,数据信息的pbt格式文本

* 范例: 生成pbt后, 需要再生成proto格式, 通过protoc将proto编译为你使用的语言代码后, 读取pbt文件
	
### Lua格式

* 格式: --lua_out=path/to/out.lua

* 功能: 生成所有表中类型信息, 数据信息及索引的lua格式脚本

* 范例: 生成lua后, 通过lua.exe解释器, 或lua嵌入代码使用require进行lua文件读取
	
### Json格式

* 格式: --json_out=path/to/out.json

* 功能: 生成所有表中类型信息, 数据信息的json格式配置

* 范例: 通过各种语言提供的json库可直接读取文件
	
	
### C#格式

* 格式: --csharp_out=path/to/out.cs

* 功能: 生成所有表中类型信息, 索引的C#格式脚本

* 范例: 生成cs后, 需要再生成bin格式, 通过[C#读取器]: (https://github.com/davyxu/tabtoy/blob/master/exportorv2/csharp) 读取二进制数据

### 二进制格式

* 格式: --binary_out=path/to/out.bin

* 功能: 生成所有表中二进制数据

* 范例: 生成bin后, 需要再生成cs格式, 通过[C#读取器]: (https://github.com/davyxu/tabtoy/blob/master/exportorv2/csharp) 读取二进制数据
	
### go格式

* 格式: --go_out=path/to/out.go

* 功能: 生成所有表中索引信息的golang代码

* 范例: 生成go后, 配合github/golang/protobuf库读取pbt格式, 再使用生成的golang文件为数据建立索引

### 类型信息

* 格式: --type_out=path/to/out.json

* 功能: 生成所有表中类型信息并输出json

* 范例: 对于没有反射功能的语言, 例如C++, 想快速的遍历所有表格的类型信息, 可以通过这个选项自行解析读取

## 指定合并结构体名

* 格式: --combinename=YourStructName

* 功能: 生成代码时, 每个表格归属的结构体所在的结构的名称来自于合并结构体名

* 范例: Item表和Quest表一同输出时, 指定合并结构体名为Config, 则输出C#文件包含Config类, 其字段包含Item和Quest的记录集合

## 指定proto文件输出版本

* 格式: --protover=2

* 功能: 设置protobuf格式的版本

	2表示proto2语法(使用protoc v2, 不带syntax识别头)

	3表示proto3语法(使用protoc v3, 带syntax识别头)
	
## 设置输出语言

* 格式: --lan=en_us

* 功能: 通过语言名, 可以设置不同的输出日志, 方便非程序员导出查错

* 范例: 语言名支持en_us(默认),  zh_cn(简体中文)


## 多文件合并

* 格式: 在输入电子表格文件名中, 使用加号(+)将要合并的文件写出来,注意+号前后不能有空格

* 功能: 将格式相同的多个电子表格内容合并

* 范例: tabtoy --mode=exportorv2 Info.xlsx+Info2.xlsx OtherFile.xlsx

## 纵向表格导出

* 方法: 在@Types表中的1,1单元格位置添加Vertical: "true" 开启功能

* 功能: 将以行延伸的表格, 适用于配置

* 范例: https://github.com/davyxu/tabtoy/blob/master/exportorv2/sample/Vertical.xlsx

# FAQ
问：如何导出结构体数组？
答：参考例子https://github.com/davyxu/tabtoy/blob/master/exportorv2/sample/Sample.xlsx
中StrStruct字段
注意： 结构体数组要求每个数组的元素在一个独立的单元格


# 例子

参考文件夹:
	
	[范例]: (https://github.com/davyxu/tabtoy/blob/master/exportorv2/sample)	
	
文件夹说明:

* csharp

	C#版通过生成的C#源码读取二进制的例子
		
* goreadjson
	
	golang通过proto生成的结构体, 读取json的例子, 也可以手动编写结构体读取
	
* goreadpbt

	golang通过proto生成的结构体, 读取Protobuf文本格式
	
* lua

	lua通过生成的lua文件, 读取数据