package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/go-chi/render"
	//"github.com/duo-labs/webauthn/protocol"
	//"github.com/duo-labs/webauthn/webauthn"
	"net/http"
)



//GET -> USer + params
func (ah *AppHandler) BeginRegistration(w http.ResponseWriter, r *http.Request) {
	//user := datastore.GetUser() // Find or create the new user
	user := &storage.Uxer{}//webauthn.User{}
	options, sessionData, err := ah.web.BeginRegistration(user)

	s := getUserSession(r.Context())
	s.WebAuthnSession = sessionData
	s.WebAuthnUser = user

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", sessionData)
	// handle errors if present
	// store the sessionData values
	//JSONResponse(w, options, http.StatusOK) // return the options generated
	// options.publicKey contain our registration options
	render.JSON(w, r, options)
}

//renderIndex
func (ah *AppHandler) webAuthnIndex(w http.ResponseWriter, r *http.Request) {
	ah.renderTemplate("webauthn.gohtml", w, r, nil)
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
	render.JSON(w, r, "Registration Success")
}


func (ah *AppHandler) BeginLogin(w http.ResponseWriter, r *http.Request) {
	//user := datastore.GetUser() // Find the user
	user := getUserSession(r.Context()).WebAuthnUser
	options, sessionData, err := ah.web.BeginLogin(user)

	if err != nil {
		panic(err)
	}

	getUserSession(r.Context()).WebAuthnSession = sessionData

	// handle errors if present
	// store the sessionData values
	render.JSON(w, r, options)
	//JSONResponse(w, options, http.StatusOK) // return the options generated

	// options.publicKey contain our registration options
}


func (ah *AppHandler) FinishLogin(w http.ResponseWriter, r *http.Request) {
	//user := datastore.GetUser() // Get the user
	user := getUserSession(r.Context()).WebAuthnUser
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

	user.AddCredential(*credential)

	// Handle validation or input errors
	// If login was successful, handle next steps
	//JSONResponse(w, "Login Success", http.StatusOK)
	render.JSON(w, r, "Login Success")
}