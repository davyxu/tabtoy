using System;
using System.IO;

namespace TabtoyExample
{
    class Program
    {
        // 加载所有表
        static void LoadAllTable()
        {            
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
        }

        // 读取指定名字的表, 可根据实际需求调整该函数适应不同加载数据来源
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

        // 指定表读取例子
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

        static void Main(string[] args)
        {
            LoadAllTable();
            LoadSpecifiedTable();
        }
    }
}
