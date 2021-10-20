package router

import (
	"fmt"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/go-chi/render"
	//"github.com/duo-labs/webauthn/protocol"
	//"github.com/duo-labs/webauthn/webauthn"
	"net/http"
)

type Uxer struct {


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
	return []webauthn.Credential{}
}


//GET -> USer + params
func (ah *AppHandler) BeginRegistration(w http.ResponseWriter, r *http.Request) {
	//user := datastore.GetUser() // Find or create the new user
	user := Uxer{}//webauthn.User{}
	options, sessionData, err := ah.web.BeginRegistration(&user)
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

/*
func (ah *AppHandler) FinishRegistration(w http.ResponseWriter, r *http.Request) {
	user := datastore.GetUser() // Get the user
	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	sessionData := store.Get(r, "registration-session")
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(r.Body)
	credential, err := ah.web.CreateCredential(&user, sessionData, parsedResponse)
	// Handle validation or input errors
	// If creation was successful, store the credential object
	JSONResponse(w, "Registration Success", http.StatusOK) // Handle next steps
}


func (ah *AppHandler) BeginLogin(w http.ResponseWriter, r *http.Request) {
	user := datastore.GetUser() // Find the user
	options, sessionData, err := webauthn.BeginLogin(&user)
	// handle errors if present
	// store the sessionData values
	JSONResponse(w, options, http.StatusOK) // return the options generated
	// options.publicKey contain our registration options
}

func (ah *AppHandler) FinishLogin(w http.ResponseWriter, r *http.Request) {
	user := datastore.GetUser() // Get the user
	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	sessionData := store.Get(r, "login-session")
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	credential, err := webauthn.ValidateLogin(&user, sessionData, parsedResponse)
	// Handle validation or input errors
	// If login was successful, handle next steps
	JSONResponse(w, "Login Success", http.StatusOK)
}*/