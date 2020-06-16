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

种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记
---|---|---|---|---|---|---|---|---
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

## Golang使用表格数据

导出命令行:
```bash
tabtoy.exe -mode=v3 -index=Index.xlsx -package=main -go_out=table_gen.json -json_out=table_gen.json
```

读取数据源码:

```go
// 重新加载指定文件名的表
func ReloadTable(filename string) {

	// 根据需要从你的源数据读取，这里从指定文件名的文件读取
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 重置数据，这里会触发Prehandler
	Tab.ResetData()

	// 使用json反序列化
	err = json.Unmarshal(data, Tab)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 构建数据和索引，这里会触发PostHandler
	Tab.BuildData()
}

```
[完整Golang例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/golang)

## C#使用表格数据


导出命令行:
```bash
tabtoy.exe -mode=v3 -index=Index.xlsx -package=main -csharp_out=table_gen.cs -binary_out=table_gen.bin
   ```

读取数据源码:

```cs
using System;
using System.IO;

using (var stream = new FileStream("table_gen.bin", FileMode.Open))
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
    

    Console.WriteLine(tab.ExampleData[3].Name);

}
```

* C#源码出于性能考虑, 默认读取tabtoy专用二进制格式

* C#也可以读取JSON数据格式, 由于C#第三方JSON不统一, 请自行使用生成的源码与第三方源码对接

[完整C#例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/csharp)

## Java使用表格数据

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

## Lua使用表格数据

导出命令行:
```bash
tabtoy.exe -mode=v3 -index=Index.xlsx -lua_out=table_gen.lua
```

读取数据源码:

```lua
-- 添加搜索路径
package.path = package.path .. ";../?.lua"

-- 加载
local t = require "table_gen"

-- 直接访问原始数据, 此处输出为UTF8格式, windows命令行下会出现乱码是正常现象
print(t.ExampleData[2].Name)

-- 通过索引访问
print(t.ExampleDataByID[300].ID)
```

[完整Lua例子](https://github.com/davyxu/tabtoy/tree/master/v3/example/lua)

## 导出表格类型信息

导出命令行:
```bash
tabtoy.exe -mode=v3 -index=Index.xlsx -jsontype_out=type_gen.json 
```

# 特色功能

## 定义和使用枚举


* 在类型表中定义枚举

种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记
---|---|---|---|---|---|---|---|---
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
种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记
---|---|---|---|---|---|---|---|---
表头 | ExampleData | 技能列表| Skill | int32 | <code>&#124;</code>   | 

技能列表 |
--- |
<code>2&#124;3</code> |
1 |

输出:

 [2, 3]
 
 [ 1 ]


## 使用多列数组

种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记
---|---|---|---|---|---|---|---|---
表头 | ExampleData | 技能列表| Skill | int32 | <code>&#124;</code>   | 

技能列表 | 技能列表
--- | --- |
<code>2&#124;3</code> | 4
1 | 

输出:

 [2, 3, 4]
 
 [ 1, 0 ]
 
 * 多列数组单元格所有数据会被自动切割并合并
 
 * 当数组字段拆分为多个同名列时, 导出数组将为空单元格默认填充类型默认值
 
 * 切勿在被拆分表中使用多列数组, 导出数据可能存在歧义

## 为字段建立索引
种类 | 对象类型 | 标识名 | 字段名 | 字段类型 | 数组切割| 值 | 索引 | 标记
---|---|---|---|---|---|---|---|---
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


## 使用标记

在类型表标记中添加字符串, 使用|做默认分割, 将在-jsontype导出json中获得字段关联的标记

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
