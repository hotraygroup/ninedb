package models

type User struct {
	GID     int
	UID     int
	Version int
	TCC     int
}

func (u *User) GetID() int {
	return u.UID
}

func (u *User) GetVer() int {
	return u.Version
}

func (u *User) Index() [][]string {
	return nil
}
