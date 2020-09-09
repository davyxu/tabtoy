-- 添加搜索路径
package.path = package.path .. ";../luadir/?.lua"

-- 一次性加载所有表
function LoadAllTable()
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
end



-- 加载指定名称的表
-- P.S. 考虑一些lua的运行限制(如luajit)的const, local限制, 应尽量将lua按文件拆分读取
function LoadSpecifiedTable()
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
end

LoadAllTable()
LoadSpecifiedTable()