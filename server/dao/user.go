package dao

var USER_BUCKET = []byte("USERSd")

type User struct {
	BaseEntity
	Name    string
	PwdHash []byte
	Salt    []byte
}

func CreateUser() *User {
	return &User{
		BaseEntity: CreateBaseEntity(),
	}
}
