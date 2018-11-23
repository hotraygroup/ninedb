package engine

type Row interface {
	GetID() int        //主键
	GetVer() int       //记录版本号
	Index() [][]string //索引
}
