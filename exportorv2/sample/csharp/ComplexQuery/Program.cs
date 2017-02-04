using System;
using System.IO;

namespace complexquery
{
    class Program
    {
        // 本例子需要第三方库支持 https://github.com/davyxu/MemQLSharp
        // 将第三方库与tabtoy放在同层即可

        static void Main(string[] args)
        {
            using (var stream = new FileStream("../../../../Config.bin", FileMode.Open))
            {
                stream.Position = 0;

                var reader = new tabtoy.DataReader(stream);

                if (!reader.ReadHeader())
                {
                    Console.WriteLine("combine file crack!");
                    return;
                }

                // 读取配置
                var config = new table.Config();
                config.Deserialize(reader);

                // 简单范围查询
                var cmql = new table.ConfigMemQL(config);

                foreach (var v in cmql.NewQuerySample()
                    .Where("ID", ">=", (Int64)101)
                    .Where("ID", "<=", (Int64)102)
                    .Result())
                {
                    var def = v as table.SampleDefine;
                    Console.WriteLine(def.ID);
                }
                


                // 枚举查询

                foreach (var v in cmql.NewQuerySample()
                    .Where("Type", ">=", table.ActorType.Monkey)
                    .Where("Type", "<=", table.ActorType.Hammer)
                    .Result())
                {
                    var def = v as table.SampleDefine;
                    Console.WriteLine(def.Type);
                }


            }



        }

    }
}
