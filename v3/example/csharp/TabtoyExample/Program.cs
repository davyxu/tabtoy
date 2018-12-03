using System;
using System.IO;

namespace TabtoyExample
{
    class Program
    {
        static void Main(string[] args)
        {
            var curr = Directory.GetCurrentDirectory();
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
                

                Console.WriteLine(tab.ExampleData[3].Name);

            }
        }
    }
}
