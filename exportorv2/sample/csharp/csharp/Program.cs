using System;
using System.IO;

namespace csharptest
{
    class Program
    {

        static void Main(string[] args)
        {
            using (var stream = new FileStream("../../../../Config.bin", FileMode.Open))
            {
                stream.Position = 0;

                var reader = new tabtoy.DataReader(stream);
                
                if ( !reader.ReadHeader( ) )
                {
                    Console.WriteLine("combine file crack!");
                    return;
                }

                var config = new table.Config();
                config.Deserialize(reader);

                // 直接通过下标获取或遍历
                var directFetch = config.Sample[2];

                // 根据索引取

                // 取不存在的元素时, 返回给定的默认值, 避免空

                // 添加日志输出或自定义输出
                config.TableLogger.AddTarget( new tabtoy.DebuggerTarget() );

                // 取空时, 当默认值不为空时, 输出日志
                var nullFetchOutLog = config.GetSampleByID( 0 );

            }
            
            
            
        }

    }
}
