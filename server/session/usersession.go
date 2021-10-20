package session

import (
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/rs/xid"
)

var sessionMap = make(map[string]*UserSession, 100)

//UserSession is bound to gonl cookie and has server side date about the user
type UserSession struct {
	SessionID     string
	UserName      string
	UserID		  int
	ComicList	  *[]storage.Comic
	SeriesList	  *[]storage.Series
	authenticated bool
	D             interface{}
	WebAuthnSession *webauthn.SessionData
	WebAuthnUser *storage.Uxer
}

//NewUserSession creates a session for Anon and stores it in the Session map
func NewUserSession() *UserSession {
	us := &UserSession{
		SessionID:     xid.New().String(),
		UserName:      "Anon",
		authenticated: false,
	}
	sessionMap[us.SessionID] = us
	return us
}

//TODO serialize session, so that we service a restart or can do loadbalancing
func (us *UserSession) save() {
	//Serialize
}

//Invalidate removes UserSession from cache
func (us *UserSession) Invalidate() {
	us.authenticated = false
	delete(sessionMap, us.SessionID)
}

func (us *UserSession) IsLoggedIn() bool {
	return us.authenticated
}

func (us *UserSession) AuthSession() {
	us.authenticated = true
}

//UserSessionByID helps to get a session which is in memory
func UserSessionByID(id string) *UserSession {
	return sessionMap[id]
}

