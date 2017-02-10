
namespace table
{
    public class ConfigMemQL
    {
        MemQL.Table Sample = new MemQL.Table(typeof(SampleDefine));

        public ConfigMemQL(Config config)
        {
            foreach (var record in config.Sample)
            {
                Sample.AddRecord(record);
            }

            Sample.GenFieldIndex("Type", ">=", 0, 3);
            Sample.GenFieldIndex("Type", "<=", 0, 3);
        }

        public MemQL.Query NewQuerySample( )
        {
            return new MemQL.Query(Sample);
        }


    }
}
