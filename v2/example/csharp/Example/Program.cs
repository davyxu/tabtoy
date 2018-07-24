using System;
using System.IO;

namespace csharptest
{
    class Program
    {

        static void Main(string[] args)
        {
            var dir = Directory.GetCurrentDirectory();
            using (var stream = new FileStream("../../Config.bin", FileMode.Open))
            {
                stream.Position = 0;

                var reader = new tabtoy.DataReader(stream);
                
                if ( !reader.ReadHeader(  ) )
                {
                    Console.WriteLine("combine file crack!");
                    return;
                }

                var config = new table.Config();
                table.Config.Deserialize(config, reader);                

                // ֱ��ͨ���±��ȡ�����
                var directFetch = config.Sample[2];

                // �����־������Զ������
                config.TableLogger.AddTarget( new tabtoy.DebuggerTarget() );

                // ȡ��ʱ, ��Ĭ��ֵ��Ϊ��ʱ, �����־
                var nullFetchOutLog = config.GetSampleByID( 0 );
                
            }
           
        }

    }
}
