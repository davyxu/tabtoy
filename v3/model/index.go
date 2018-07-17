package model

type TableKind int32

const (
	TableKind_None     TableKind = iota //
	TableKind_Type                      // 类型表
	TableKind_Data                      // 数据表
	TableKind_KeyValue                  // 键值表
)

type IndexDefine struct {
	Kind          TableKind `tb_name:"模式"`
	TableType     string    `tb_name:"表类型"`
	TableFileName string    `tb_name:"表文件名"`
}
