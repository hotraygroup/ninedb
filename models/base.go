package models

type Base struct {
	ID int
}

func (b *Base) GetID() int {
	return b.ID
}
