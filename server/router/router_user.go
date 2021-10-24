package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/command"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)



func (ah *AppHandler) serveTemplate(t string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ah.renderTemplate(t, w, r, data)
	}
}

func (ah *AppHandler) listUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access listUsers")
	}
}

func (ah *AppHandler) logoutUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getUserSession(r.Context()).Invalidate()
		http.Redirect(w,r,"/comics",303)
	}
}


func (ah *AppHandler) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access deleteUser")
	}
}

func (ah *AppHandler) updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access updateUser")
	}
}
func (ah *AppHandler) getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access getUser")
	}
}
func (ah *AppHandler) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("username")
		pass := r.FormValue("pass")
		repass := r.FormValue("repass")
		est := ""
		if name == "" {
			est += "Keine Name? "
		}

		if pass == "" {
			est += "Kein Password! "
		}

		if pass != repass {
			est += "Passwörter nicht gleich! "
		}
		if est != "" {
			ah.renderTemplate("create_user.gohtml", w, r, est)
		}

		//TODO check for username

		if ok := ah.dao.AddUser(name,pass); !ok {
			est = "Duplicate Username"
			ah.renderTemplate("create_user.gohtml", w, r, est)
		} else {
			us := getUserSession(r.Context())
			us.UserName = name
			http.Redirect(w,r,"/comics",303)
		}
	}
}

func (ah *AppHandler) loginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("username")
		pass := r.FormValue("pass")
		user := ah.dao.AuthUser(name, pass)
		if user == nil {
			ah.renderTemplate("login_user.gohtml", w, r, "Login Failed")
			return
		}
		us := getUserSession(r.Context())
		us.AuthSession()
		us.UserName = user.Name
		us.UserID = user.Id

		http.Redirect(w,r,"/comics",303)
	}
}




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
		render.JSON(w, r, command.NewRedirectCommand("/comics"))
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
	render.JSON(w, r, command.NewRedirectCommand("/comics"))
}