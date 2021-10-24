package router

import (
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

//renderIndex
func (ah *AppHandler) webAuthnIndex(w http.ResponseWriter, r *http.Request) {
	ah.renderTemplate("webauthn.gohtml", w, r, nil)
}



//GET -> USer + params
//called first
// check if user exists
func (ah *AppHandler) BeginRegistration(w http.ResponseWriter, r *http.Request) {

	username := chi.URLParam(r, "userID")
	tempUser := &storage.User{}
	tempUser.Name = username
	options, sessionData, err := ah.web.BeginRegistration(tempUser)

	s := getUserSession(r.Context())
	s.WebAuthnSession = sessionData
	s.WebAuthnUser = tempUser

	if err != nil {
		panic(err)
	}

	render.JSON(w, r, options)
}


func (ah *AppHandler) FinishRegistration(w http.ResponseWriter, r *http.Request) {
	//user := datastore.GetUser() // Get the user
	s := getUserSession(r.Context())
	user := s.WebAuthnUser

	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(r.Body)
	if err != nil {
		panic(err)
	}

	sessionData := s.WebAuthnSession
	cred, err := ah.web.CreateCredential(user, *sessionData, parsedResponse)
	if err != nil {
		panic(err)
	}

	// Handle validation or input errors
	// If creation was successful, store the credential object
	user.AddCredential(*cred)

	//JSONResponse(w, "Registration Success", http.StatusOK) // Handle next steps

	if ah.dao.AddWebAuthnUser(user) {
		s.AuthSession()
		s.UserName = user.Name
		s.UserID = user.Id
		render.JSON(w, r, "Registration Success")
	} else {
		render.JSON(w, r, "Registration FAILED")
	}


}

//Start of auth
//Check for user in DB
func (ah *AppHandler) BeginLogin(w http.ResponseWriter, r *http.Request) {
	//user := datastore.GetUser() // Find the user
	user := ah.dao.GetWebAuthnUser(chi.URLParam(r, "userID"))
	options, sessionData, err := ah.web.BeginLogin(user)

	if err != nil {
		panic(err)
	}

	getUserSession(r.Context()).WebAuthnSession = sessionData
	getUserSession(r.Context()).WebAuthnUser = user
	// handle errors if present
	// store the sessionData values
	render.JSON(w, r, options)
	//JSONResponse(w, options, http.StatusOK) // return the options generated

	// options.publicKey contain our registration options
}


func (ah *AppHandler) FinishLogin(w http.ResponseWriter, r *http.Request) {
	//user := datastore.GetUser() // Get the user
	us := getUserSession(r.Context())
	user := us.WebAuthnUser
	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	//sessionData := store.Get(r, "login-session")
	sessionData := getUserSession(r.Context()).WebAuthnSession

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	if err != nil {
		panic(err)
	}
	credential, err := ah.web.ValidateLogin(user, *sessionData, parsedResponse)
	if err != nil {
		panic(err)
	}

	//FIXME update credentials ... check for clones
	user.AddCredential(*credential)
	us.AuthSession()
	us.UserName = user.Name
	us.UserID = user.Id
	// Handle validation or input errors
	// If login was successful, handle next steps
	//JSONResponse(w, "Login Success", http.StatusOK)
	render.JSON(w, r, "Login Success")
}