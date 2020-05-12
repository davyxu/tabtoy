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
        void Deserialize(TableReader reader);
    }
    

    public class TableReader
    {
        BinaryReader _binaryReader;
        long _boundPos;

        public TableReader(Stream stream)
        {
            _binaryReader = new BinaryReader(stream);
            _boundPos = stream.Length;
        }

        public TableReader(Stream stream, long boundpos)
        {
            _binaryReader = new BinaryReader(stream);
            _boundPos = boundpos;
        }

        public TableReader(TableReader reader, long boundpos)
        {
            _binaryReader = reader._binaryReader;
            _boundPos = boundpos;
        }

        void ConsumeData(uint size)
        {
            if (!IsDataEnough(size))
            {
                throw new Exception("Out of struct bound");
            }
        }

        bool IsDataEnough(uint size)
        {
            return _binaryReader.BaseStream.Position + size <= _boundPos;
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

            UInt32 ver = 0;
            ReadUInt32(ref ver);
            if (ver != FileVersion)
            {
                throw new Exception("Invalid tabtoy version");
            }            
        }

        public bool ReadTag(ref UInt32 v)
        {
            if (IsDataEnough(sizeof(UInt32)))
            {
                v = _binaryReader.ReadUInt32();
                return true;
            }

            return false;
        }

        public void SkipFiled(UInt32 tag)
        {
            var fieldIndex = tag & 0xffff;
            switch (tag >> 16)
            {
                case 1: // int16
                    {
                        Int16 dummy = 0;
                        this.ReadInt16(ref dummy);
                    }
                    break;
                case 2: // int32
                case 10: // enum
                    {
                        Int32 dummy = 0;
                        this.ReadInt32(ref dummy);
                    }
                    break;
                case 3: // int64
                    {
                        Int64 dummy = 0;
                        this.ReadInt64(ref dummy);
                    }
                    break;
                case 4: // uint16
                    {
                        UInt16 dummy = 0;
                        this.ReadUInt16(ref dummy);
                    }
                    break;
                case 5: // uint32
                    {
                        UInt32 dummy = 0;
                        this.ReadUInt32(ref dummy);
                    }
                    break;
                case 6: // uint64
                    {
                        UInt64 dummy = 0;
                        this.ReadUInt64(ref dummy);
                    }
                    break;
                case 7: // float32
                    {
                        float dummy = 0;
                        this.ReadFloat(ref dummy);                        
                    }
                    break;
                case 8: // string
                    {
                        string dummy = string.Empty;
                        this.ReadString(ref dummy);
                    }
                    break;
                case 9: // bool
                    {
                        bool dummy = false;
                        this.ReadBool(ref dummy);
                    }
                    break;


                case 101: // int16
                    {
                        var dummy = new List<Int16>();
                        this.ReadInt16(ref dummy);
                    }
                    break;
                case 102: // int32
                case 110: // enum
                    {
                        var dummy = new List<Int32>();
                        this.ReadInt32(ref dummy);
                    }
                    break;
                case 103: // int64
                    {
                        var dummy = new List<Int64>();
                        this.ReadInt64(ref dummy);
                    }
                    break;
                case 104: // uint16
                    {
                        var dummy = new List<UInt16>();
                        this.ReadUInt16(ref dummy);
                    }
                    break;
                case 105: // uint32
                    {
                        var dummy = new List<UInt32>();
                        this.ReadUInt32(ref dummy);
                    }
                    break;
                case 106: // uint64
                    {
                        var dummy = new List<UInt64>();
                        this.ReadUInt64(ref dummy);
                    }
                    break;
                case 107: // float32
                    {
                        var dummy = new List<float>();
                        this.ReadFloat(ref dummy);
                    }
                    break;
                case 108: // string
                    {
                        var dummy = new List<string>();
                        this.ReadString(ref dummy);
                    }
                    break;
                case 109: // bool
                    {
                        var dummy = new List<bool>();
                        this.ReadBool(ref dummy);
                    }
                    break;
                default:
                    throw new Exception("Invalid tag type");
            }
        }


        static readonly UTF8Encoding encoding = new UTF8Encoding();

        public void ReadInt16(ref Int16 v)
        {
            ConsumeData(sizeof(Int16));

            v = _binaryReader.ReadInt16();
        }

        public void ReadInt32(ref Int32 v)
        {
            ConsumeData(sizeof(Int32));

            v = _binaryReader.ReadInt32();
        }
       

        public void ReadInt64(ref Int64 v)
        {
            ConsumeData(sizeof(Int64));

            v = _binaryReader.ReadInt64();
        }

        public void ReadUInt16(ref UInt16 v)
        {
            ConsumeData(sizeof(UInt16));

            v = _binaryReader.ReadUInt16();
        }

        public void ReadUInt32(ref UInt32 v)
        {
            ConsumeData(sizeof(UInt32));

            v = _binaryReader.ReadUInt32();
        }

        public void ReadUInt64(ref UInt64 v)
        {
            ConsumeData(sizeof(UInt64));

            v = _binaryReader.ReadUInt64();
        }

        public void ReadFloat(ref float v)
        {
            ConsumeData(sizeof(float));

            v = _binaryReader.ReadSingle();
        }

        public void ReadBool(ref bool v)
        {
            ConsumeData(sizeof(bool));

            v = _binaryReader.ReadBoolean();
        }

        public void ReadString(ref string v)
        {
            UInt32 len = 0;
            ReadUInt32(ref len);

            ConsumeData(sizeof(Byte) * len);

            v = encoding.GetString(_binaryReader.ReadBytes((int)len));
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
            UInt32 bound = 0;
            ReadUInt32(ref bound);

            v = new T();

            // 避免不同结构体跨越式影响其他数据二进制边界
            v.Deserialize(new TableReader(this, _binaryReader.BaseStream.Position + bound));
        }

        public void ReadStruct<T>(ref List<T> v) where T : ITableSerializable, new()
        {
            UInt32 len = 0;
            ReadUInt32(ref len);

            for (int i = 0; i < len; i++)
            {
                T element = default(T);
                ReadStruct<T>(ref element);
                v.Add(element);
            }
        }
    }
}