using System;
using System.IO;
using tabtoy;

namespace csharptest
{
    class Program
    {

        static void Main(string[] args)
        {            
            using (var stream = new FileStream("../../Config.bin", FileMode.Open))
            {
                stream.Position = 0;

                var reader = new tabtoy.DataReader(stream);

                var config = new table.Config();

                var result = reader.ReadHeader(config.GetBuildID());
                if ( result != FileState.OK)
                {
                    Console.WriteLine("combine file crack!");
                    return;
                }

                
                table.Config.Deserialize(config, reader);

                // 直接通过下标获取或遍历
                var directFetch = config.Sample[2];

                // 添加日志输出或自定义输出
                config.TableLogger.AddTarget(new tabtoy.DebuggerTarget());

                // 取空时, 当默认值不为空时, 输出日志
                var nullFetchOutLog = config.GetSampleByID(0);

            }
           
        }

    }
}
