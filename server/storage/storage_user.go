package storage

import (
	"bytes"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/rs/xid"
	"golang.org/x/crypto/argon2"
)



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



type Uxer struct {
	creds []webauthn.Credential

}


func (user *Uxer) WebAuthnID() []byte {
	return []byte("sdas")
}

func (user *Uxer) WebAuthnName() string {
	return "newUser"
}

func (user *Uxer) WebAuthnDisplayName() string {
	return "New User"
}

func (user *Uxer) WebAuthnIcon() string {
	return "https://pics.com/avatar.png"
}

func (user *Uxer) WebAuthnCredentials() []webauthn.Credential {
	if user.creds == nil {
		user.creds = []webauthn.Credential{}
	}
	return user.creds
}

func (user *Uxer) AddCredential(credential webauthn.Credential){
	user.WebAuthnCredentials() //make sure the array exists
	user.creds = append(user.creds, credential)
}