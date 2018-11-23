package models

type User struct {
	Base
	GID int
	TCC int
}

func (u *User) Index() [][]string {
	return nil
}
