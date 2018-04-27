
#include <iostream>
#include <vector>
#include <cstdio>
#include <memory>

namespace tabtoy
{
    enum class LogLevel
    {
        Debug,
        Info,
        Warnning,
        Error,
    };

    class LogTarget
    {
    public:
         virtual void WriteLog(LogLevel level, const char* msg)
        {

        }

    public:
     static std::string LevelToString( LogLevel level )
        {
            switch (level)
            {
                case LogLevel::Debug:
                    return "tabtoy [Debug] ";
                case LogLevel::Info:
                    return "tabtoy [Info] ";
                case LogLevel::Warnning:
                    return "tabtoy [Warn] ";
                case LogLevel::Error:
                    return "tabtoy [Error] ";
            }

            return "tabtoy [Unknown] ";
        }
	};



	class DebuggerTarget : public LogTarget
    {
	public:
		void virtual WriteLog(LogLevel level, const char* msg) override
        {
            std::cout <<  (LevelToString(level) + msg).c_str() << "\n";
        }
	};


    class Logger
    {
        std::vector<std::unique_ptr<LogTarget>> _targets;

	public:
		void AddTarget(LogTarget* tgt)
        {
            _targets.push_back(std::unique_ptr<LogTarget>(tgt));
        }

	public:
		void ClearTargets( )
        {
            _targets.clear();
        }

        void WriteLine(LogLevel level, const char* msg)
        {
            for(auto& tgt :_targets )
            {
                tgt->WriteLog(level, msg);
            }
        }

	public:
		template <typename ...Args>
		void DebugLine(const char* fmt, Args... args)
        {
			int sz = std::snprintf(nullptr, 0, fmt, args...);
			std::vector<char> buf(sz + 1);
			std::snprintf(&buf[0], buf.size(), args...);

            WriteLine(LogLevel::Debug, &buf[0]);
        }

	public:
		template <typename ...Args>
		void InfoLine(const char* fmt, Args... args)
        {
			int sz = std::snprintf(nullptr, 0, fmt, args...);
			std::vector<char> buf(sz + 1);
			std::snprintf(&buf[0], buf.size(), fmt, args...);

            WriteLine(LogLevel::Info, &buf[0]);
        }

	public:
		template <typename ...Args>
		void WarningLine(const char* fmt, Args... args)
        {
			int sz = std::snprintf(nullptr, 0, fmt, args...);
			std::vector<char> buf(sz + 1);
			std::snprintf(&buf[0], buf.size(), fmt, args...);

            WriteLine(LogLevel::Warnning, &buf[0]);
        }

	public:
		template <typename ...Args>
		void ErrorLine(const char* fmt, Args... args)
        {
			int sz = std::snprintf(nullptr, 0, fmt, args...);
			std::vector<char> buf(sz + 1);
			std::snprintf(&buf[0], buf.size(), fmt, args...);

            WriteLine(LogLevel::Error, &buf[0]);
        }
	};
}
