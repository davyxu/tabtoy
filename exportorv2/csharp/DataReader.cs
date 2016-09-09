using System.IO;
using System.Text;

namespace tabtoy
{
    public class DataReader : BinaryReader
    {
        public DataReader( Stream stream )
            : base( stream )
        {            
        }

        static readonly UTF8Encoding encoding = new UTF8Encoding();

        public string ReadUTF8String( )
        {
            var len = this.ReadInt32();
            return encoding.GetString(this.ReadBytes(len));
        }
    }
}
