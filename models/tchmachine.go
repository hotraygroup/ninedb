package models

type TchMachine struct {
	ID      int
	Version int
	UID     int
	GID     int
}

func (t *TchMachine) GetID() int {
	return t.ID
}

func (t *TchMachine) GetVer() int {
	return t.Version
}

//索引列的值暂不支持修改，如需修改，可以先删除记录，再插入新记录
func (t *TchMachine) Index() [][]string {
	return [][]string{
		{"GID", "UID"},
		{"GID"},
		{"UID"},
	}

}

func (t *TchMachine) TableName() string {
	return "TchMachine"
}
