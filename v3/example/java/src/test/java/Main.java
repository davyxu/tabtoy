import table.Table;
import com.alibaba.fastjson.JSON;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Map;

public class Main {

    // 从文件读取数据
    private static String readFileAsString(String fileName)throws Exception
    {
        return new String(Files.readAllBytes(Paths.get(fileName)));
    }
    public static void main(String[] args) throws Exception {

        // 从文件读取配置表
        String data = null;
        try {
            data = readFileAsString("./cfg/table_gen.json");
        } catch (Exception e) {
            e.printStackTrace();
        }

        // 表格数据
        Table tab;

        // 从json序列化出对象
        tab = JSON.parseObject(data, Table.class);

        if(tab == null){
            throw new Exception("parse table failed");
        }

        // 构建索引
        tab.BuildData();

        // 测试输出
        for(Map.Entry<Integer, Table.ExampleData> def : tab.ExampleDataByID.entrySet()){
            System.out.println(def.getValue().Name);
        }
    }
}