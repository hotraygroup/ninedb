package engine

type Row interface {
	GetID() int        //主键
	Index() [][]string //索引
}
