#include <fstream>
#include <string>
#include <exception>

namespace tabtoy
{
	enum class FieldType
	{
		None = 0,
		Int32 = 1,
		Int64 = 2,
		UInt32 = 3,
		UInt64 = 4,
		Float = 5,
		String = 6,
		Bool = 7,
		Enum = 8,
		Struct = 9,
	};	

	class DataReader
    {       
		std::streamoff _boundPos  = -1;

		std::ifstream& _reader;

	public:
		DataReader(std::ifstream& stream )
			: _reader(stream)
        {
			stream.seekg(0, std::ios::end);
            _boundPos = stream.tellg();
			stream.seekg(0);
        }

	public:
		DataReader(std::ifstream& stream, long boundpos)
			: _reader(stream)
        {
            _boundPos = boundpos;
        }

	public:
		DataReader(DataReader& reader, long boundpos )
			: _reader(reader._reader)
        {
            _boundPos = boundpos;
        }

	private:
        void ConsumeData(int size)
        {          
            if ( !IsDataEnough( size ) )
            {
                throw new std::out_of_range("Out of struct bound");
            }            
        }

        bool IsDataEnough(int size)
        {            
            return (int)_reader.tellg() + size <= _boundPos;
        }

        const int CombineFileVersion = 2;

	public:
		bool ReadHeader( )
        {            
            auto tag = ReadString();
            if (tag != "TABTOY")
            {
                return false;
            }

            auto ver = ReadInt32();
            if (ver != CombineFileVersion)
            {
                return false;
            }
           
            return true;
        }

	public:
		int ReadTag()
        {
            if ( IsDataEnough(sizeof(int) ) )
            {
                return ReadInt32( );
            }

            return -1;
        }
   
	public:
		int ReadInt32( )
        {
            ConsumeData(sizeof(int));
			int n;
			_reader.read((char*)&n, sizeof(int));
            return n;
        }

	public:
		long long ReadInt64( )
        {
            ConsumeData(sizeof(long long ));
			long long n;
			_reader.read((char*)&n, sizeof(long long));
            return n;
        }

	public:
		unsigned int ReadUInt32( )
        {
            ConsumeData(sizeof(unsigned int));
			unsigned int n;
			_reader.read((char*)&n, sizeof(unsigned int));
            return n;
        }

	public:
		unsigned long long ReadUInt64( )
        {
            ConsumeData(sizeof(unsigned long long));
			unsigned long long n;
			_reader.read((char*)&n, sizeof(unsigned long long));
            return n;
        }

	public:
		float ReadFloat( )
        {
            ConsumeData(sizeof(float));
			float f;
			_reader.read((char*)&f, sizeof(float));
			return f;
        }

	public:
		bool ReadBool( )
        {
            ConsumeData(sizeof(bool));
			bool b;
			_reader.read((char*)&b, sizeof(bool));
            return b;
        }

	public:
		std::string ReadString()
        {
            auto len = ReadInt32();

            ConsumeData(sizeof(char) * len);
			
			std::string str;
			str.resize(len);
			_reader.read((char*)str.data(), len);
            return str;
        }

	public:
		template <typename T>
		T ReadEnum( )
        {
            return (T)ReadInt32();
        }

	public:
		template <typename T>
		T ReadStruct(void (*handler)(T&, DataReader&))
        {
            auto bound = ReadInt32();

            T element;

            DataReader reader(*this, (int)_reader.tellg() + bound);
            handler(element,reader);

            return element;
        }
        
	};
};
