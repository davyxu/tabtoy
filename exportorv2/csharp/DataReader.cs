using System;
using System.IO;
using System.Text;
using System.Collections.Generic;

namespace tabtoy
{
    enum FieldType
    {
        None    = 0,
	    Int32   = 1,
	    Int64   = 2,
	    UInt32  = 3,
	    UInt64  = 4,
	    Float   = 5,
	    String  = 6,
	    Bool    = 7,
	    Enum    = 8,
	    Struct  = 9,
	    Bytes   = 10, // 暂时为binaryfile输出使用
    }

    public interface DataObject
    {
        void Deserialize(DataReader reader);
    }

    public class DataReader
    {
        BinaryReader _reader;

        public DataReader( Stream straem )
        {
            _reader = new BinaryReader(straem );
        }

       
        const int CombineFileVersion = 1;

        public bool ReadHeader()
        {
            var tag = RawReadString();
            if (tag != "TABTOY")
            {
                return false;
            }

            var ver = _reader.ReadInt32();
            if (ver != CombineFileVersion)
            {
                return false;
            }

            return true;
        }

        static readonly UTF8Encoding encoding = new UTF8Encoding();
        public string RawReadString( )
        {
            var len = _reader.ReadInt32();
            return encoding.GetString(_reader.ReadBytes(len));
        }

        int pedingtag;
        public bool MatchTag( int expect )
        {
            if (pedingtag == 0 )
            {
                pedingtag = _reader.ReadInt32();                
            }

            if ( expect == pedingtag )
            {
                pedingtag = 0;
                return true;
            }

            return false;
        }
   
        public int ReadInt32( )
        {
            return _reader.ReadInt32();
        }

        public long ReadInt64( )
        {
            return _reader.ReadInt64();
        }

        public uint ReadUInt32( )
        {
            return _reader.ReadUInt32();
        }

        public ulong ReadUInt64( )
        {
            return _reader.ReadUInt64();
        }

        public float ReadFloat( )
        {
            return _reader.ReadSingle();
        }

        public bool ReadBool( )
        {
            return _reader.ReadBoolean();
        }

        public string ReadString( )
        {
            return RawReadString();
        }

        public T ReadEnum<T>( )
        {
            return (T)Enum.ToObject(typeof(T), _reader.ReadInt32());                
        }


        public void ReadList_Int32( List<int> list)
        {
            var c = _reader.ReadInt32();

            for (int i = 0; i < c; i++)
            {
                var element = _reader.ReadInt32();

                list.Add(element);
            }
        }

        public void ReadList_Int64(List<long> list)
        {
            var c = _reader.ReadInt32();

            for (int i = 0; i < c; i++)
            {
                var element = _reader.ReadInt64();

                list.Add(element);
            }
        }

        public void ReadList_UInt32( List<uint> list)
        {
            var c = _reader.ReadInt32();

            for (int i = 0; i < c; i++)
            {
                var element = _reader.ReadUInt32();

                list.Add(element);
            }
        }

        public void ReadList_UInt64( List<ulong> list)
        {
            var c = _reader.ReadInt32();

            for (int i = 0; i < c; i++)
            {
                var element = _reader.ReadUInt64();

                list.Add(element);
            }
        }

        public void ReadList_Float(List<float> list)
        {
            var c = _reader.ReadInt32();

            for (int i = 0; i < c; i++)
            {
                var element = _reader.ReadSingle();

                list.Add(element);
            }
        }

        public void ReadList_String(List<string> list)
        {
            var c = _reader.ReadInt32();

            for (int i = 0; i < c; i++)
            {
                var element = RawReadString();

                list.Add(element);
            }
        }


        public void ReadList_Bool(List<bool> list)
        {
            var c = _reader.ReadInt32();

            for (int i = 0; i < c; i++)
            {
                var element = _reader.ReadBoolean();

                list.Add(element);
            }
        }

        public void ReadList_Enum<T>(List<T> list)
        {
            var c = _reader.ReadInt32();            

            for (int i = 0; i < c; i++)
            {
                var element = (T)Enum.ToObject( typeof(T), _reader.ReadInt32() );

                list.Add(element);
            }
        }

        public void ReadList_Struct<T>(List<T> list) where T : tabtoy.DataObject
        {
            var c = _reader.ReadInt32();

            for (int i = 0; i < c; i++)
            {
                var element = Activator.CreateInstance<T>();

                element.Deserialize(this);

                list.Add(element);
            }
        }



        
    }
}
