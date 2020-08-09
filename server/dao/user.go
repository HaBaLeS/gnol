package dao

var USER_BUCKET = []byte("USERS")

type User struct {
	BaseEntity
	Name    string
	PwdHash []byte
	Salt    []byte
	Comics  []*UserComic
}

type UserComic struct {
	Name string
	Tags []string
	MetaDataID string
}

func NewUserComic(metadata Metadata) *UserComic{
	return &UserComic{
		Name: metadata.Name,
		Tags: nil,
		MetaDataID: metadata.Id,
	}
}

func (u *User) AddComic(metadata Metadata){
	cm := NewUserComic(metadata)
	u.Comics = append(u.Comics, cm)
	u.Save()
}

func (u *User) Save() {

}

func LoadUser(userID string) *User{

}