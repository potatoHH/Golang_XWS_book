package migrator

type Entity interface {
	ID() int64         //需要返回的id
	TableName() string //表名
	//CompareTo dst 必然也是Entity
	CompareTo(dst Entity) bool //比较两个entity
	Columns() []string         //返回字段
}
