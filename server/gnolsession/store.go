package gnolsession

/*type GnolSessionStore struct {
	sessions.Store
	session *sessions.Session
}

func (gs GnolSessionStore) Options(options sessions2.Options) {
	//TODO implement me
	panic("implement me")
}

func NewGnolSessionStore() GnolSessionStore {
	gs := GnolSessionStore{}
	gs.session = sessions.NewSession(gs, "egal")
	return gs
}

// Get should return a cached session.
func (gs GnolSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {

	return gs.session, nil
}

// New should create and return a new session.
//
// Note that New should never return a nil session, even in the case of
// an error if using the Registry infrastructure to cache the session.
func (gs GnolSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return nil, nil
}

// Save should persist session to the underlying store implementation.
func (gs GnolSessionStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return nil
}
*/
