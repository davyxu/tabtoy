package table;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public final class Table {

	public class ExampleData{
		public int ID;//任务ID
		public String Name;//名称
		public float Rate;//比例
		public ActorType ActorType;//类型
		public int[] Skill;//技能列表
	}


	class ExampleKV {
		public String ServerIP;//服务器IP
		public int ServerPort; //服务器端口
		public int[] GroupID;//分组`
	}

	 public  enum ActorType {

		 None(0),
		 Pharah(1),// 法鸡
		 Junkrat(2),// 狂鼠
		 Genji(3),// 源氏
		 Mercy(4),// 天使
	     ;

       private ActorType(int ActorType) {
           this.ActorType = ActorType;
       }

	   public int getActorType() {
           return ActorType;
       }

       private final int ActorType;
	    }


	public List<ExampleData> ExampleData = new ArrayList<ExampleData>();
	// table: ExampleKV
	public List<ExampleKV> ExampleKV = new ArrayList<ExampleKV>();

	// Indices
	public Map<Integer,ExampleData> ExampleDataByID = new HashMap<Integer,ExampleData>();

	public  ExampleKV GetKeyValue_ExampleKV() {
		return ExampleKV.get(0);
	}

	public interface TableEvent{
		void OnPreProcess( );
		void OnPostProcess( );
	}

	// Handlers
	private List<TableEvent> eventHandlers = new ArrayList<TableEvent>();

	// 注册加载后回调
	public void RegisterEventHandler(TableEvent ev ){

		eventHandlers.add(ev);
	}

	// 清除索引和数据, 在处理前调用OnPostProcess, 可能抛出异常
	public void ResetData()  {
		for( TableEvent ev : eventHandlers){
			ev.OnPreProcess();
		}
		ExampleData.clear();
		ExampleKV.clear();
		ExampleDataByID.clear();
	}

	// 构建索引, 需要捕获OnPostProcess可能抛出的异常
	public void  BuildData()  {
		// 遍历赋值 TOOD
		for(ExampleData v:ExampleData) {
			ExampleDataByID.put(v.ID, v);
		}

		for( TableEvent ev : eventHandlers){
			ev.OnPostProcess();
		}
	}

}
