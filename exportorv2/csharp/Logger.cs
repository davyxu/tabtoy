using System.Collections.Generic;
using System.Diagnostics;

namespace tabtoy
{
    public enum LogLevel
    {
        Debug,
        Info,
        Warnning,
        Error,        
    }

    public class LogTarget
    {
        public virtual void WriteLog(LogLevel level, string msg)
        {

        }

        public static string LevelToString( LogLevel level )
        {
            switch (level)
            {
                case LogLevel.Debug:
                    return "tabtoy [Debug] ";                    
                case LogLevel.Info:
                    return "tabtoy [Info] ";                    
                case LogLevel.Warnning:
                    return "tabtoy [Warn] ";
                case LogLevel.Error:
                    return "tabtoy [Error] ";
            }

            return "tabtoy [Unknown] ";
        }
    }



    public class DebuggerTarget : LogTarget
    {        
        public override void WriteLog(LogLevel level, string msg)
        {   
            Debug.WriteLine( LevelToString(level) + msg );
        }
    }


    public class Logger
    {
        List<LogTarget> _targets = new List<LogTarget>();

        public void AddTarget(LogTarget tgt)
        {
            _targets.Add(tgt);
        }

        public void ClearTargets( )
        {
            _targets.Clear();
        }

        void WriteLine(LogLevel level, string msg)
        {
            foreach(var tgt in _targets )
            {
                tgt.WriteLog(level, msg);
            }
        }

        public void DebugLine(string fmt, params object[] args)
        {
            var text = string.Format(fmt, args);

            WriteLine(LogLevel.Debug, text);
        }

        public void InfoLine(string fmt, params object[] args)
        {
            var text = string.Format(fmt, args);

            WriteLine(LogLevel.Info, text);
        }

        public void WarningLine(string fmt, params object[] args)
        {
            var text = string.Format(fmt, args);

            WriteLine(LogLevel.Warnning, text);
        }

        public void ErrorLine(string fmt, params object[] args)
        {
            var text = string.Format(fmt, args);

            WriteLine(LogLevel.Error, text);
        }
    }
}
