package session

import "github.com/rs/xid"

var sessionMap = make(map[string]*UserSession, 100)

type UserSession struct {
	SessionID     string
	UserName      string
	authenticated bool
	D             interface{}
}

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

func UserSessionByID(id string) *UserSession {
	return sessionMap[id]
}
