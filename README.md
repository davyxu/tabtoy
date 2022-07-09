# tabtoy v3

高性能表格数据导出工具

![tabtoylogo](doc/logo.png)

# 特性
* 支持Xlsx/CSV作为表格数据混合输入

* 支持JSON/Golang/C#/Java/Lua/二进制 源码, 数据, 类型输出

* 自动单元格数据格式检查, 精确到单元格的报错

* 支持预定义枚举, 可使用中文枚举类型

* 支持表拆分, 支持多人协作

* 支持KV配置表, 方便将表格作为配置文件

* 多核并发导出, 缓存加速, 上百文件秒级导出

# 迭代历程

* 2020年6月: tabtoy v3
    支持Xlsx/CSV混合导出
    
    新的表格格式
        
    重构代码

* 2016年8月: 第六代导出器,tabtoy v2 调整为以电子表格为中心的方式, 支持v1 90%常用功能

	增加: 所有导出文件均为1个文件, 提高加载读取速度

	增加: 二进制合并导出(第五代导出器需要使用2个工具才能完成)
	
	增加: C#源码导出及索引创建,无需protobuf支持
	
	增加: proto格式导出, 支持v2,v3格式
		
	重构代码, 导出速度更快

* 2016年3月: 第五代导出器,tabtoy v1 在四代基础上重构,开源,支持并发导出	

* 2015年: 第四代导出器,基于Golang导出器,增加ID重复检查,数组格的多重写法, 支持a.b.c栏位导出, 导出速度大大提高

* 2013年: 第三代导出器,在二代基础上做到内容格式与导出器独立,但依然依赖csv前置导出,增加逗号分隔格子内容,导出速度慢

* 2012年: 第二代导出器,基于C++和Protobuf的导出器,内容格式与导出器混合编写,需要vbs导出csv,速度慢
	
* 2011年: 第一代导出器,基于VBA的表格内建导出器,速度慢,复用困难,容易错,不安全

# 导出第一个表

## 类型表
准备一个电子表格命名为: Type.xlsx

类型表用于定义表格中表头以及用到的类型

表格内容如下:

种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记 | 备注
---|---|---|---|---|---|---|---|---|---
表头 | MyData | ID | ID | int32|
表头 | MyData | 名称 | Name | string|

## 数据表
准备一个电子表格命名为: MyData.xlsx

表格内容如下:

ID | 名称
---|---
1 | 坦克
2 | 法师

## 索引表
* 准备一个电子表格命名为: Index.xlsx

模式 | 表类型 | 表文件名
---|---|---
类型表 |        | Type.xlsx
数据表 | MyData | MyData.xlsx

注意 数据表的表类型需要与类型表里的对象类型对应

## 编写导出shell


[下载tabtoy](https://github.com/davyxu/tabtoy/releases)

```bash
tabtoy.exe -mode=v3 -index=Index.xlsx -json_out=table_gen.json
```

[完整例子文件](https://github.com/davyxu/tabtoy/tree/master/v3/example/tutorial)

# 导出数据/源码/类型

## Golang使用表格导出的JSON数据

导出命令行:
```bash
tabtoy.exe -mode=v3 -index=Index.xlsx -package=main -go_out=table_gen.json -json_out=table_gen.json
```

读取数据源码:

```go
	var Tab = NewTable()

	// 表加载前清除之前的手动索引和表关联数据
	Tab.RegisterPreEntry(func(tab *Table) error {
		fmt.Println("tab pre load clear")
		return nil
	})

	// 表加载和构建索引后，需要手动处理数据的回调
	Tab.RegisterPostEntry(func(tab *Table) error {
		fmt.Println("tab post load done")
		fmt.Printf("%+v\n", tab.ExampleDataByID[200])

		fmt.Println("KV: ", tab.GetKeyValue_ExampleKV().ServerIP)
		return nil
	})

	err := tabtoy.LoadFromFile(Tab, "../json/table_gen.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
```
[完整Golang例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/golang)

## C#使用表格导出二进制数据

导出命令行:
```bash
tabtoy.exe -mode=v3 -index=Index.xlsx -package=main -csharp_out=table_gen.cs -binary_out=table_gen.bin
   ```

读取数据源码:

```cs
using (var stream = new FileStream("../../../../binary/table_gen.bin", FileMode.Open))
{
    stream.Position = 0;

    var reader = new tabtoy.TableReader(stream);


    var tab = new main.Table();

    try
    {
        tab.Deserialize(reader);
    }
    catch (Exception e)
    {
        Console.WriteLine(e);
        throw;
    }
    
    // 建立所有数据的索引
    tab.IndexData();

    // 表遍历
    foreach (var kv in tab.ExampleData) 
    {
        Console.Write("{0} {1}\n",kv.ID, kv.Name);
    }

    // 直接取值
    Console.WriteLine(tab.ExtendData[1].Additive);

    // KV配置
    Console.WriteLine(tab.GetKeyValue_ExampleKV().ServerIP);
}
```

* C#源码出于性能考虑, 默认读取tabtoy专用二进制格式

* C#也可以读取JSON数据格式, 由于C#第三方JSON不统一, 请自行使用生成的源码与第三方源码对接

[完整C#例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/csharp)

## Java使用表格导出的JSON数据

导出命令行:
```bash
tabtoy.exe -mode=v3 -index=Index.xlsx -package=main -java_out=Table.java -json_out=table_gen.json
```

读取数据源码:

```java
import main.Table;
import com.alibaba.fastjson.JSON;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Map;

public class Main {

    // 从文件读取数据
    private static String readFileAsString(String fileName)throws Exception
    {
        return new String(Files.readAllBytes(Paths.get(fileName)));
    }
    public static void main(String[] args) throws Exception {

        // 从文件读取配置表
        String data = null;
        try {
            data = readFileAsString("table_gen.json");
        } catch (Exception e) {
            e.printStackTrace();
        }

        // 表格数据
        Table tab;

        // 从json序列化出对象
        tab = JSON.parseObject(data, Table.class);

        if(tab == null){
            throw new Exception("parse table failed");
        }

        // 构建索引
        tab.BuildData();

        // 测试输出
        for(Map.Entry<Integer, Table.ExampleData> def : tab.ExampleDataByID.entrySet()){
            System.out.println(def.getValue().Name);
        }
    }
}
```

[完整Java例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/java)

## Lua使用表格导出的Lua数据(测试中)

导出命令行:
```bash
tabtoy.exe -mode=v3 -index=Index.xlsx -lua_out=table_gen.lua
```

读取数据源码:

```lua
    -- 加载
    local tab = {}
    require("table_gen").init(tab)
    
    -- 遍历表
    print("Iterate lua table by order:")
    for _, v in ipairs(tab.ExampleData) do
        print(v.ID, v.Name)
    end

    -- 通过索引访问
    print("Access index table data:")
    print(tab.ExampleDataByID[300].ID)

    -- 枚举类型访问
    print("Use generated enum:")
    print(tab.ActorType.Pharah,  tab.ActorType[3])
```

[完整Lua例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/lua)

## 将表格类型信息导出为JSON格式

导出命令行:
```bash
tabtoy -mode=v3 -index=Index.xlsx -jsontype_out=type_gen.json 
```

## 导出为Protobuf格式
tabtoy可以将表类型及结构输出为Google Protobuf的proto格式, 同时输出与之对应的二进制格式(*.pbb)

使用Protobuf的SDK即可方便的将表数据提供给所有Protobuf支持的语言

以下例子展示Golang使用Protobuf读取表格输出文件

* 导出proto文件:
```bash
tabtoy -mode=v3 -index=Index.xlsx -proto_out=table.proto 
```

* 导出proto二进制数据文件:
```bash
tabtoy -mode=v3 -index=Index.xlsx -pbbin_out=all.pbb
```

* Protobuf编译器protoc下载

下载地址: https://github.com/protocolbuffers/protobuf/releases

* 安装Golang的Protobuf生成插件
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

* 将proto文件生成代码
```bash
protoc --go_out=. ./table.proto -I .
```

[完整Golang使用Protobuf例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/protobuf/golang)

# 按表导出

tabtoy默认情况下, 均是将数据, 源码一次性导出.出于以下原因,tabtoy支持按表导出数据

* 某些语言在读取大量数据时, 会出现兼容性问题. 例如: lua的local和const限制等

* 按需读取数据, 降低内存需求

* 按需更新数据, 减少模块耦合

## Golang按需读取JSON数据

导出命令行:
```bash
tabtoy -mode=v3 -index=Index.xlsx -package=main -go_out=table_gen.json -json_dir=.
```

读取数据源码:

```go
	var TabData = NewTable()
	err := tabtoy.LoadTableFromFile(TabData, "../jsondir/ExampleData.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("load specified table: ExampleData")
	for k, v := range TabData.ExampleDataByID {
		fmt.Println(k, v)
	}

	// 分表加载时, 不会触发pre/post Handler
	var TabKV = NewTable()
	err = tabtoy.LoadTableFromFile(TabKV, "../jsondir/ExampleKV.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("load specified table: ExampleKV")
	for k, v := range TabKV.ExampleKV {
		fmt.Println(k, v)
	}
```
[完整Golang例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/golang)

## Lua按需读取Lua数据

导出命令行:
```bash
tabtoy -mode=v3 -index=Index.xlsx -lua_dir=.
```

读取数据源码:

```lua
    local tabData = {}
    require("ExampleData").init(tabData)
    require("ExtendData").init(tabData)

    print("Load 2 tables into one lua table:")
    for _, v in ipairs(tabData.ExampleData) do
        print(v.ID, v.Name)
    end
    for _, v in ipairs(tabData.ExtendData) do
        print(v.Additive)
    end

    print("Load kv table into single lua table:")
    local kvData = {}
    require("ExampleKV").init(kvData)
    for _, v in ipairs(kvData.ExampleKV) do
        print(v.ServerIP, v.ServerPort)
    end

    -- lua枚举是可选功能, 根据需要加载
    local tabType = {}
    require("_TableType").init(tabType)
    print("Use generated enum:")
    print(tabType.ActorType.Pharah,  tabType.ActorType[3])
```

[完整Lua例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/lua)
[导出的Lua表](https://github.com/davyxu/tabtoy/tree/master/v3/example/luasrc)

## C#按需读取二进制数据
导出命令行:
```bash
tabtoy -mode=v3 -index=Index.xlsx -package=main -csharp_out=table_gen.cs -binary_dir=.
   ```

读取数据源码:

```cs
 static void LoadTableByName(main.Table tab,  string tableName)
{
    using (var stream = new FileStream(string.Format("../../../../binary/{0}.bin", tableName), FileMode.Open))
    {
        stream.Position = 0;

        var reader = new tabtoy.TableReader(stream);
        try
        {
            tab.Deserialize(reader);
        }
        catch (Exception e)
        {
            Console.WriteLine(e);
            throw;
        }
    }
    
    tab.IndexData(tableName);
}

static void LoadSpecifiedTable()
{
    var tabData = new main.Table();

    LoadTableByName(tabData, "ExampleData");
    LoadTableByName(tabData, "ExtendData");

    Console.WriteLine("Load table merged into one class");
    // 表遍历
    foreach (var kv in tabData.ExampleData)
    {
        Console.Write("{0} {1}\n", kv.ID, kv.Name);
    }
    // 表遍历
    foreach (var kv in tabData.ExtendData)
    {
        Console.Write("{0}\n", kv.Additive);
    }

    Console.WriteLine("Load KV table into one class");
    var tabKV = new main.Table();
    LoadTableByName(tabKV, "ExampleKV");

    // KV配置
    Console.WriteLine(tabKV.GetKeyValue_ExampleKV().ServerIP);
}
```

[完整C#例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/csharp)

## Golang使用Protobuf按需读取二进制数据
[Golang例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/protobuf/golang)



# 特色功能

## 定义和使用枚举


* 在类型表中定义枚举

种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记 | 备注
---|---|---|---|---|---|---|---|---|---
枚举 | ActorType |   | None | int32|  | 0
枚举 | ActorType | 法鸡 | Pharah | int32|  | 1
枚举 | ActorType | 狂鼠 | Junkrat | string| | 2
枚举 | ActorType | 源氏 | Genji | int32|  | 3
枚举 | ActorType | 天使 | Mercy | string| | 4
表头 | ExampleData | 类型 | Type | ActorType

* 在数据表中使用枚举

类型 |
--- |
狂鼠 |
Genji |

* 在数据表的枚举字段中, 枚举 字段名或标识名都会自动识别对应枚举值

* 枚举只有枚举数值会被导出. 枚举标识名, 字段名均不会出现在数据中

## 使用数组
种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记 | 备注
---|---|---|---|---|---|---|---|---|---
表头 | ExampleData | 技能列表| Skill | int32 | <code>&#124;</code>   |

技能列表 |
--- |
<code>2&#124;3</code> |
1 |

输出:

 [2, 3]
 
 [ 1 ]

## 使用多列数组

种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记 | 备注
---|---|---|---|---|---|---|---|---|---
表头 | ExampleData | 技能列表| Skill | int32 | <code>&#124;</code>   |

技能列表 | 技能列表
--- | --- |
2 | 3
1 | 

输出:

 [2, 3 ]
 
 [ 1, 0 ]
 
 * 多列数组单元格所有数据会被自动切割并合并
 
 * 当数组字段拆分为多个同名列时, 导出数组将为空单元格默认填充类型默认值, 保证多列导出后, 数组数量统一
 
 * 切勿在拆分表中使用多列数组, 导出数据可能存在歧义

## 为字段建立索引
种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记 | 备注
---|---|---|---|---|---|---|---|---|---
表头 | ExampleData | ID| ID | int32 |  | |是| 


生成代码中, 会自动对数据创建索引, 例如:
```go
ExampleDataByID map[int32]*ExampleData
```

## 表拆分

将ExampleData表, 拆为Data.csv, Data2.csv表

模式 | 表类型 | 表文件名
---|---|---
类型表 |        | Type.xlsx
数据表 | ExampleData | Data.csv
数据表 | ExampleData | Data2.csv

每个表中的字段可按需填写

## KV表

准备类型表:

模式 | 表类型 | 表文件名
---|---|---
类型表 |        | Type.xlsx
数据表 | ExampleKV | KV.csv

准备KV表:

字段名 | 字段类型 | 标识名 |  值|  数组切割 | 标记
---|---|---|---|---|---|
ServerIP | string | 服务器IP | 8.8.8.8
ServerPort | uint16 | 服务器端口 | 1024  

## 空行分割

表格数据如下:

ID | 名称
---|---
1 | 坦克
2 | 法师
(空行)  |
3 | 治疗

导出数据
* 1 坦克
* 2 法师

导表工具在识别到空行后, 空行后的数据将被忽略

## 行数据注释

表格数据如下:

ID | 名称
---|---
1 | 坦克
#2 | 法师
3 | 治疗

导出数据
* 1 坦克
* 3 治疗

在任意表的首列单元格中首字符为#时，该行所有数据不会被导出

## 列数据注释

表格数据如下:

ID | #名称
---|---
1 | 坦克
2 | 法师
3 | 治疗

导出数据
* 1 
* 2
* 3 

表头中, 列字段首字符为#时，该列所有数据按默认值导出 

## 不导出指定表
实现此功能需要使用到TagAction, 参考下面例子配置:

在Index表中:

模式 | 表类型 | 表文件名 | 标记
---|---|---|---|
数据表 | Effect | Effect.csv | client
数据表 | Password | Server.csv | server

* 客户端数据导出
导出参数中新增参数
```shell script
--tag_action=nogentab:server
```
表示, 不导出带有server标记的所有表格

* 服务器数据导出
导出参数中新增参数
```shell script
--tag_action=nogentab:client
```
表示, 不导出带有client标记的所有表格

## 不输出指定列数据
实现此功能需要使用到TagAction, 参考下面例子配置:

在Type表中:

种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记 | 备注
---|---|---|---|---|---|---|---|---|---
表头 | ExampleData | 特效ID| EffectID | int32 |  | | | client 
表头 | ExampleData | 概率| Rate | float |  | | | server

表中的特效ID, 只希望客户端导出数据中包含EffectID, 同时服务器导出数据中只包含Rate, 不希望将Rate字段导入客户端数据
客户端导出为二进制, 服务器导出为json

此时在相应字段所在的Type表中的"标记" 一列增加如表所示标记(如标记列不存在, 请新建)

将原有导出流程拆分为客户端导出和服务器导出, 分两次分别导出不同需求的数据

* 客户端数据导出
导出参数中新增参数
```shell script
--tag_action=nogenfield_binary:server
```
表示: server标记的字段不导出到二进制

* 服务器数据导出  
导出参数中新增参数
```shell script
--tag_action=nogenfield_json:client
```
表示: client标记的字段不导出到json完整文件

## TagAction参考说明

### 格式
```shell script
--tag_action=action1:tag1+tag2|action2:tag1+tag3
```
* | 表示多个action
* 被标记的tag, 将被对应action处理

### action类型
action | 适用范围 | 功能
---|---|---|
nogenfield_json | Type表 | 被标记的字段不导出到json完整文件中
nogenfield_jsondir| Type表 | 被标记的字段不导出到每个表文件json
nogenfield_binary| Type表 | 被标记的字段不导出到二进制中
nogenfield_pbbin| Type表 | 被标记的字段不导出到Protobuf二进制中
nogenfield_lua| Type表 | 被标记的字段不导出到Lua中
nogenfield_csharp| Type表 | 被标记的字段不导出到C#中
nogentab| Index表 | 被标记的表不会导出到任何输出中

## 启用缓冲
命令行中加入-usecache=true, 将启用缓存功能, 加速导出速度

-cachedir参数设定缓存目录, 默认换出到tabtoy当前目录下的.tabtoycache目录

# FAQ

* 怎么让客户端和服务器通过标记分别导出

    请为客户端和服务器分别编写两个tabtoy导出过程分别导出

# V2版本说明


* tabtoy 同时支持V2, V3版本导出

* V2版本将不再获得更新

* [V2文档](https://github.com/davyxu/tabtoy/blob/master/README_v2.md)

# 备注

感觉不错请star, 谢谢!

知乎: [http://www.zhihu.com/people/sunicdavy](http://www.zhihu.com/people/sunicdavy)

提交bug及特性: [https://github.com/davyxu/tabtoy/issues](https://github.com/davyxu/tabtoy/issues)
