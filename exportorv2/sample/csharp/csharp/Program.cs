using System;
using System.IO;
using tabtoy;

namespace csharptest
{
    class Program
    {

        static void Main(string[] args)
        {
            using (var stream = new FileStream("../../../../Config.bin", FileMode.Open))
            {
                stream.Position = 0;

                var reader = new DataReader(stream);
                
                if ( !reader.ReadHeader( ) )
                {
                    Console.WriteLine("combine file crack!");
                    return;
                }

                var file = new gamedef.Config();
                file.Deserialize(reader);
            }
            
            
            
        }

    }
}
