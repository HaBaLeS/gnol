package dao

import (
	"bytes"

	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"golang.org/x/crypto/argon2"
)

func (dao *DAO) AllUsers() []*storage.User {
	retList := make([]*storage.User, 0)
	err := dao.DB.Select(&retList, "select id, name  from gnoluser")
	if err != nil {
		panic(err)
	}
	return retList
}

func (dao *DAO) GetUserForApiToken(gt string) (error, int) {
	var uid int
	err := dao.DB.Get(&uid, "select us.id from gnoluser us, apitoken at where us.id = at.user_id and at.token = $1", gt)
	if err != nil {
		return err, -1
	}
	return nil, uid
}

func (dao *DAO) GetOrCreateAPItoken(id int) []string {
	var res []string
	err := dao.DB.Select(&res, "select token from apitoken where user_id = $1", id)
	if err != nil {
		panic(err)
	}

	if len(res) == 0 {
		newToken := uuid.New().String()
		dao.DB.MustExec("insert into apitoken (user_id,token) values ($1,$2)", id, newToken)

		return dao.GetOrCreateAPItoken(id)
	}
	return res
}

func (dao *DAO) GetUser(userId int) (*storage.User, error) {
	user := new(storage.User)
	err := dao.DB.Get(user, "select * from gnoluser where id = $1", userId)
	return user, err
}

func (dao *DAO) AuthUser(name string, pass string) *storage.User {
	user := new(storage.User) //TODO why use new?
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

func (dao *DAO) AddWebAuthnUser(user *storage.User) bool {

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

func (dao *DAO) GetWebAuthnUser(username string) *storage.User {
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

// FIXME move somewhere out of DAO
func hashPassword(pass string) ([]byte, []byte) {
	salt := xid.New().Bytes()
	hash := argon2.Key([]byte(pass), salt, 3, 32*1024, 4, 32)
	return hash, salt
}

// FIXME move somewhere out of DAO
func checkPassword(salt []byte, dbhash []byte, pass string) bool {
	hash := argon2.Key([]byte(pass), salt, 3, 32*1024, 4, 32)
	if bytes.Compare(hash, dbhash) != 0 { //FIXME introduce constant time comparison
		return false
	}
	return true
}
