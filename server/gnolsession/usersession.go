package gnolsession

import (
	"github.com/rs/xid"
)

// UserSession is bound to gonl cookie and has server side date about the user
type UserSession struct {
	SessionID     string
	UserName      string
	UserID        int
	Authenticated bool
	//WebAuthnSession *webauthn.SessionData //fixme add to registration Session
	//WebAuthnUser *storage.User //fixme add to registration Session
}

// NewUserSession creates a gnolsession for Anon and stores it in the Session map
func NewUserSession() *UserSession {
	us := &UserSession{
		SessionID:     xid.New().String(),
		UserName:      "Anon",
		Authenticated: false,
	}
	return us
}

// TODO serialize gnolsession, so that we service a restart or can do loadbalancing
func (us *UserSession) save() {
	//Serialize
}

// Invalidate removes UserSession from cache
func (us *UserSession) Invalidate() {
	us.Authenticated = false
}

func (us *UserSession) IsLoggedIn() bool {
	return us.Authenticated
}

func (us *UserSession) AuthSession() {
	us.Authenticated = true
}
