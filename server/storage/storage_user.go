package storage

import (
	"bytes"
	"github.com/rs/xid"
	"golang.org/x/crypto/argon2"
)


/*
type UserStore struct {
	bs *BoltStorage
}

func newUserStore(bs *BoltStorage) *UserStore{
	return &UserStore{
		bs: bs,
	}
}*/



/*
func (us *UserStore) AddComic(uid string, metadata *Metadata) error {
	/*u := us.UserByID([]byte(uid))
	if u == nil {
		return fmt.Errorf("user with id %s does not exist", uid)
	}

	u.MetadataList = append(u.MetadataList, metadata.Id)
	metadata.owners = append(metadata.owners, uid)

	err := us.bs.Write(metadata)
	err = us.bs.Write(u)
	return err
	return nil
}*/


/*
func (us *UserStore) UserByID(uid []byte) *User {
	user := &User{}
	err := us.bs.ReadRaw(func(tx *bolt.Tx) error {
		b := tx.Bucket(USER_BUCKET)
		k, v := b.Cursor().Seek(uid)
		if k != nil && bytes.Equal(k, uid) {
			der := loadFromJson(user, v)
			return der
		} else {
			return fmt.Errorf("user with Id %s not found", uid)
		}
	})
	if err != nil {
		fmt.Printf("Error Loading User. %s\n", err)
		return nil
	}
	return  user
}
*/

func (dao *DAO) AuthUser(name string, pass string) *User {
	user := new(User)
	err := dao.DB.Get(user,"select * from gnoluser where name = $1",name)
	if err != nil {
		dao.log.Printf("Error querying for user: %v", err)
		return nil
	}
	auth := checkPassword(user.Salt, user.PasswordHash, pass)
	if !auth {
		return nil
	} else {
		return user
	}
}



func (dao *DAO)AddUser(name,  password string) bool {
	password_hash, salt := hashPassword(password)
	tx := dao.DB.MustBegin()
	_, err := tx.Exec("INSERT INTO gnoluser (name, password_hash, salt) VALUES ($1, $2, $3)", name, password_hash, salt)
	if err != nil {
		dao.log.Printf("Could not insert user. ")
		tx.Rollback()
		return false
	}
	err = tx.Commit()
	if err !=  nil {
		dao.log.Printf("Could not insert user. %v", err)
		return false
	}
	return true
}

func hashPassword(pass string) ([]byte, []byte) {
	salt := xid.New().Bytes()
	hash := argon2.Key([]byte(pass), salt, 3, 32*1024, 4, 32)
	return hash, salt
}


func checkPassword(salt []byte, dbhash []byte, pass string) bool {
	hash := argon2.Key([]byte(pass), salt, 3, 32*1024, 4, 32)
	if bytes.Compare(hash, dbhash) != 0 { //FIXME introduce constant time comparison
		return false
	}
	return true
}


func AddComic(name string) {
	//"insert into comic (a,b,v) values ($1,$2, $3)"
}

func AddComicToUser(c Comic, u User) {
	//"insert into user_to_comic (a,b,v) values ($1,$2, $3)"

}

func ListComicsForUser(u User) *[]Comic{
	//"select * comic joine to user hwre user = ?"
	return nil
}