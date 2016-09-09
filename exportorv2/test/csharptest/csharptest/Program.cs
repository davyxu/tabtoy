using System.IO;
using tabtoy;

namespace csharptest
{
    class Program
    {

        static void Main(string[] args)
        {
            using (var stream = new FileStream("../../../../Sample.bin", FileMode.Open))
            {
                stream.Position = 0;

                var reader = new DataReader(stream);
                

                var file = new gamedef.SampleFile();
                file.Deserialize(reader);

            }
            
            
            
        }

    }
}
