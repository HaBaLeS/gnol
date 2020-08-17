package storage

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/rs/xid"
	"golang.org/x/crypto/argon2"
)

var USER_BUCKET = []byte("USERS")

type User struct {
	*BaseEntity
	Name    string
	PwdHash []byte
	Salt    []byte
	MetadataList  []string
}

type UserStore struct {
	bs *BoltStorage
}

func newUserStore(bs *BoltStorage) *UserStore{
	return &UserStore{
		bs: bs,
	}
}



func (us *UserStore) CreateUser(name string, pass string) *User {
	//FIXME users can just override each others
	hash, salt := hashPassword(pass)
	u := &User{
		BaseEntity: CreateBaseEntity(USER_BUCKET),
		MetadataList: make([]string,0),
		Name:       name,
		PwdHash:    hash,
		Salt:       salt,
	}
	us.bs.Write(u)
	return u
}


func (us *UserStore) AddComic(uid string, metadata *Metadata) error {
	u := us.UserByID([]byte(uid))
	if u == nil {
		return fmt.Errorf("user with id %s does not exist", uid)
	}

	u.MetadataList = append(u.MetadataList, metadata.Id)
	metadata.owners = append(metadata.owners, uid)

	err := us.bs.Write(metadata)
	err = us.bs.Write(u)
	return err
}

func hashPassword(pass string) ([]byte, []byte) {
	salt := xid.New().Bytes()
	hash := argon2.Key([]byte(pass), salt, 3, 32*1024, 4, 32)
	return hash, salt
}

func checkPassword(salt []byte, dbhash []byte, pass string) error {
	hash := argon2.Key([]byte(pass), salt, 3, 32*1024, 4, 32)
	if bytes.Compare(hash, dbhash) != 0 {
		fmt.Println("Password do not match")
		return fmt.Errorf("Login Error")
	}
	return nil
}


func (us *UserStore) AuthUser(name string, pass string) (*User, error) {
	u := new(User)
	logError := us.bs.ReadRaw(func(tx *bolt.Tx) error {
		c := tx.Bucket(USER_BUCKET).Cursor()
		//spx := []byte(name)
		//FIXME introduce search instead of scan //&& bytes.HasPrefix(k, spx)
		for k, v := c.First(); k != nil ; k, v = c.Next() {
			err := loadFromJson(u, v)
			if err != nil {
				return err
			}
			if u.Name == name {
				return checkPassword(u.Salt, u.PwdHash, pass)
			}
		}
		fmt.Printf("Did not find: %s in DB\n", name)
		return fmt.Errorf("Login failed")
	})
	return u, logError
}

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
