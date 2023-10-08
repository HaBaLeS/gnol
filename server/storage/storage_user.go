package storage

import (
	"bytes"
	"github.com/rs/xid"
	"golang.org/x/crypto/argon2"
)

var SELECT_WEBAUTN_CRED = "select " +
	"wa.aagu_id, wa.signcount, wa.clonewarning, " +
	"wc.id, wc.publicKey, wc.attestationType " +
	"from webauthn_credential wc, webauthn_authenticator wa " +
	"where wc.user_id = $1 and wc.authenticator_id = wa.id"

func (dao *DAO) AuthUser(name string, pass string) *User {
	user := new(User)
	err := dao.DB.Get(user, "select * from gnoluser where name = $1", name)
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

func (dao *DAO) AddUser(name, password string) bool {
	password_hash, salt := hashPassword(password)
	tx := dao.DB.MustBegin()
	_, err := tx.Exec("INSERT INTO gnoluser (name, password_hash, salt) VALUES ($1, $2, $3)", name, password_hash, salt)
	if err != nil {
		dao.log.Printf("Could not insert user. %v", err)
		tx.Rollback()
		return false
	}
	err = tx.Commit()
	if err != nil {
		dao.log.Printf("Could not insert user. %v", err)
		return false
	}
	return true
}

func (dao *DAO) AddWebAuthnUser(user *User) bool {

	/*creds := user.creds[0]

	tx := dao.DB.MustBegin()
	res, err := tx.Exec("INSERT INTO gnoluser (name, password_hash, salt, webauthn) VALUES ($1, $2, $3, $4)", user.Name, "", "", true)
	if err != nil {
		dao.log.Printf("Could not insert user. %v", err) //xx
		tx.Rollback()
		return false
	}
	uid, _ := res.LastInsertId()
	aid, _ := tx.MustExec("insert into webauthn_authenticator (aagu_id, signcount) values ($1,$2)", creds.Authenticator.AAGUID, creds.Authenticator.SignCount).LastInsertId()
	tx.MustExec("insert into webauthn_credential (id, publicKey, attestationType, authenticator_id, user_id) values ($1, $2, $3, $4, $5 )", creds.ID, creds.PublicKey, creds.AttestationType, aid, uid)

	err = tx.Commit()
	if err != nil {
		dao.log.Printf("Could not insert user. %v", err)
		return false
	}
	user.Id = int(uid)

	return true*/
	return false
}

func (dao *DAO) GetWebAuthnUser(username string) *User {
	/*user := new(User)
	err := dao.DB.Get(user, "select * from gnoluser where name =$1 and webauthn = true", username)
	if err != nil {
		return nil
	}
	row := dao.DB.QueryRow(SELECT_WEBAUTN_CRED, user.Id)
	if row.Err() != nil {
		return nil
	}

	a := webauthn.Authenticator{}
	c := webauthn.Credential{}
	err = row.Scan(&a.AAGUID, &a.SignCount, &a.CloneWarning, &c.ID, &c.PublicKey, &c.AttestationType)
	if err != nil {
		panic(err)
	}
	c.Authenticator = a
	user.creds = make([]webauthn.Credential, 0)
	user.creds = append(user.creds, c)

	*/
	return nil
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

type User struct {
	Id           int
	Name         string
	PasswordHash []byte `db:"password_hash"`
	Salt         []byte
	WebAuthn     bool `db:"webauthn"`
	//creds        []webauthn.Credential
}

func (user *User) WebAuthnID() []byte {
	return []byte(user.Name)
}

func (user *User) WebAuthnName() string {
	return user.Name
}

func (user *User) WebAuthnDisplayName() string {
	return user.Name
}

func (user *User) WebAuthnIcon() string {
	return "https://pics.com/avatar.png"
}

/*func (user *User) WebAuthnCredentials() []webauthn.Credential {
	if user.creds == nil {
		user.creds = []webauthn.Credential{}
	}
	return user.creds
}

func (user *User) AddCredential(credential webauthn.Credential) {
	user.WebAuthnCredentials() //make sure the array exists
	user.creds = append(user.creds, credential)
}*/
