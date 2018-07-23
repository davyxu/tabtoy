-- 添加搜索路径
package.path = package.path .. ";../?.lua"

-- 加载
local t = require "table_gen"

-- 直接访问原始数据, 此处输出为UTF8格式, windows命令行下会出现乱码是正常现象
print(t.ExampleData[2].Name)

-- 通过索引访问
print(t.ExampleDataByID[300].ID)