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
    }

    public delegate void DeserializeHandler<T>(T ins, DataReader reader);

    public class DataReader
    {
        BinaryReader _reader;
        long _boundPos  = -1;

        public DataReader(Stream stream )
        {
            _reader = new BinaryReader(stream);
            _boundPos = stream.Length;
        }

        public DataReader(Stream stream, long boundpos)
        {
            _reader = new BinaryReader(stream );
            _boundPos = boundpos;
        }

        public DataReader(DataReader reader, long boundpos )
        {
            _reader = reader._reader;
            _boundPos = boundpos;
        }

        void ConsumeData(int size)
        {          
            if ( !IsDataEnough( size ) )
            {
                throw new Exception("Out of struct bound");
            }            
        }

        bool IsDataEnough(int size)
        {            
            return _reader.BaseStream.Position + size <= _boundPos;
        }

        const int CombineFileVersion = 2;

        public bool ReadHeader( )
        {            
            var tag = ReadString();
            if (tag != "TABTOY")
            {
                return false;
            }

            var ver = ReadInt32();
            if (ver != CombineFileVersion)
            {
                return false;
            }
           
            return true;
        }

        static readonly UTF8Encoding encoding = new UTF8Encoding();

        public int ReadTag()
        {
            if ( IsDataEnough(sizeof(Int32) ) )
            {
                return ReadInt32( );
            }

            return -1;
        }
   
        public int ReadInt32( )
        {
            ConsumeData(sizeof(Int32));

            return _reader.ReadInt32();
        }

        public long ReadInt64( )
        {
            ConsumeData(sizeof(Int64));

            return _reader.ReadInt64();
        }

        public uint ReadUInt32( )
        {
            ConsumeData(sizeof(UInt32));

            return _reader.ReadUInt32();
        }

        public ulong ReadUInt64( )
        {
            ConsumeData(sizeof(UInt64));

            return _reader.ReadUInt64();
        }

        public float ReadFloat( )
        {
            ConsumeData(sizeof(float));

            return _reader.ReadSingle();
        }

        public bool ReadBool( )
        {
            ConsumeData(sizeof(bool));

            return _reader.ReadBoolean();
        }

        public string ReadString()
        {
            var len = ReadInt32();

            ConsumeData(sizeof(Byte) * len);

            return encoding.GetString(_reader.ReadBytes(len));
        }

        public T ReadEnum<T>( )
        {
            return (T)Enum.ToObject(typeof(T), ReadInt32());                
        }

        public T ReadStruct<T>(DeserializeHandler<T> handler) where T : class
        {
            var bound = _reader.ReadInt32();

            var element = Activator.CreateInstance<T>();

            handler(element, new DataReader(this, _reader.BaseStream.Position + bound));

            return element;
        }
        
    }
}
