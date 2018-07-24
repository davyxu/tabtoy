using System;
using System.Collections.Generic;
using System.IO;
using System.Text;

namespace tabtoy
{
    internal enum FieldType
    {
        None = 0,
        Int16 = 1,
        Int32 = 2,
        Int64 = 3,
        UInt16 = 4,
        UInt32 = 5,
        UInt64 = 6,
        Float = 7,
        String = 8,
        Bool = 9,
        Enum = 10,        
    }

    public interface ITableSerializable
    {
        void Deserialize(DataReader reader);
    }
    

    public class DataReader
    {
        BinaryReader _reader;
        long _boundPos;

        public DataReader(Stream stream)
        {
            _reader = new BinaryReader(stream);
            _boundPos = stream.Length;
        }

        public DataReader(Stream stream, long boundpos)
        {
            _reader = new BinaryReader(stream);
            _boundPos = boundpos;
        }

        public DataReader(DataReader reader, long boundpos)
        {
            _reader = reader._reader;
            _boundPos = boundpos;
        }

        void ConsumeData(int size)
        {
            if (!IsDataEnough(size))
            {
                throw new Exception("Out of struct bound");
            }
        }

        bool IsDataEnough(int size)
        {
            return _reader.BaseStream.Position + size <= _boundPos;
        }

        const int FileVersion = 3;

        public void ReadHeader()
        {
            string header = string.Empty;
            ReadString(ref header);
            if (header != "TABTOY")
            {
                throw new Exception("Invalid tabtoy file");
            }

            Int32 ver = 0;
            ReadInt32(ref ver);
            if (ver != FileVersion)
            {
                throw new Exception("Invalid tabtoy version");
            }            
        }

        public bool ReadTag(ref Int32 v)
        {
            if (IsDataEnough(sizeof(Int32)))
            {
                v = _reader.ReadInt32();
                return true;
            }

            return false;
        }

        static readonly UTF8Encoding encoding = new UTF8Encoding();

        public void ReadInt16(ref Int16 v)
        {
            ConsumeData(sizeof(Int16));

            v = _reader.ReadInt16();
        }

        public void ReadInt32(ref Int32 v)
        {
            ConsumeData(sizeof(Int32));

            v = _reader.ReadInt32();
        }
       

        public void ReadInt64(ref Int64 v)
        {
            ConsumeData(sizeof(Int64));

            v = _reader.ReadInt64();
        }

        public void ReadUInt16(ref UInt16 v)
        {
            ConsumeData(sizeof(UInt16));

            v = _reader.ReadUInt16();
        }

        public void ReadUInt32(ref UInt32 v)
        {
            ConsumeData(sizeof(UInt32));

            v = _reader.ReadUInt32();
        }

        public void ReadUInt64(ref UInt64 v)
        {
            ConsumeData(sizeof(UInt64));

            v = _reader.ReadUInt64();
        }

        public void ReadFloat(ref float v)
        {
            ConsumeData(sizeof(float));

            v = _reader.ReadSingle();
        }

        public void ReadBool(ref bool v)
        {
            ConsumeData(sizeof(bool));

            v = _reader.ReadBoolean();
        }

        public void ReadString(ref string v)
        {
            Int32 len = 0;
            ReadInt32(ref len);

            ConsumeData(sizeof(Byte) * len);

            v = encoding.GetString(_reader.ReadBytes(len));
        }

        public void ReadEnum<T>(ref T v)
        {
            Int32 value = 0;
            ReadInt32(ref value);

            v = (T) Enum.ToObject(typeof(T), value);
        }

        public void ReadInt16(ref List<Int16> v)
        {
            Int16 value = 0;
            ReadInt16(ref value);
            v.Add(value);
        }

        public void ReadInt32(ref List<Int32> v)
        {
            Int32 value = 0;
            ReadInt32(ref value);
            v.Add(value);
        }

        public void ReadInt64(ref List<Int64> v)
        {
            Int64 value = 0;
            ReadInt64(ref value);
            v.Add(value);
        }

        public void ReadUInt16(ref List<UInt16> v)
        {
            UInt16 value = 0;
            ReadUInt16(ref value);
            v.Add(value);
        }

        public void ReadUInt32(ref List<UInt32> v)
        {
            UInt32 value = 0;
            ReadUInt32(ref value);
            v.Add(value);
        }

        public void ReadUInt64(ref List<UInt64> v)
        {
            UInt64 value = 0;
            ReadUInt64(ref value);
            v.Add(value);
        }

        public void ReadBool(ref List<bool> v)
        {
            bool value = false;
            ReadBool(ref value);
            v.Add(value);
        }

        public void ReadString(ref List<string> v)
        {
            string value = string.Empty;
            ReadString(ref value);
            v.Add(value);
        }

        public void ReadFloat(ref List<float> v)
        {
            float value = 0;
            ReadFloat(ref value);
            v.Add(value);
        }

        public void ReadEnum<T>(ref List<T> v)
        {
            T value = default(T);
            ReadEnum(ref value);
            v.Add(value);
        }

        public void ReadStruct<T>(ref T v) where T : ITableSerializable, new()
        {
            Int32 bound = 0;
            ReadInt32(ref bound);

            v = new T();

            // 避免不同结构体跨越式影响其他数据二进制边界
            v.Deserialize(new DataReader(this, _reader.BaseStream.Position + bound));
        }

        public void ReadStruct<T>(ref List<T> v) where T : ITableSerializable, new()
        {
            Int32 len = 0;
            ReadInt32(ref len);

            for (int i = 0; i < len; i++)
            {
                T element = default(T);
                ReadStruct<T>(ref element);
                v.Add(element);
            }
        }
    }
}