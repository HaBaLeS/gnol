package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ah *AppHandler) serveTemplate(t string, data interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ah.renderTemplate(t, ctx, data)
	}
}

func (ah *AppHandler) listUsers(ctx *gin.Context) {
	fmt.Println("Access listUsers")
}

func (ah *AppHandler) logoutUser(ctx *gin.Context) {
	getUserSession(ctx).Invalidate()
	ctx.Redirect(303, "/users/login")
}

func (ah *AppHandler) deleteUser(ctx *gin.Context) {
	fmt.Println("Access deleteUser")
}

func (ah *AppHandler) updateUser(ctx *gin.Context) {
	fmt.Println("Access updateUser")
}
func (ah *AppHandler) getUser(ctx *gin.Context) {
	fmt.Println("Access getUser")
}

func (ah *AppHandler) APIToken(ctx *gin.Context) {
	id := getUserSession(ctx).UserID
	token := ah.dao.GetOrCreateAPItoken(id)
	ctx.JSON(http.StatusOK, token)
}

func (ah *AppHandler) createUser(ctx *gin.Context) {
	name := ctx.PostForm("username")
	pass := ctx.PostForm("pass")
	repass := ctx.PostForm("repass")
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
		ah.renderTemplate("create_user.gohtml", ctx, est)
	}

	//TODO check for username

	if ok := ah.dao.AddUser(name, pass); !ok {
		est = "Duplicate Username"
		ah.renderTemplate("create_user.gohtml", ctx, est)
	} else {
		us := getUserSession(ctx)
		us.UserName = name
		ctx.Redirect(303, "/series")
	}
}

func (ah *AppHandler) loginUser(ctx *gin.Context) {
	name := ctx.PostForm("username")
	pass := ctx.PostForm("pass")
	user := ah.dao.AuthUser(name, pass)
	if user == nil {
		ah.renderTemplate("login_user.gohtml", ctx, "Login Failed")
		return
	}
	us := getUserSession(ctx)
	us.AuthSession()
	us.UserName = user.Name
	us.UserID = user.Id
	updateUSerSession(ctx, us)

	ctx.Redirect(303, "/series")
}

// renderIndex
func (ah *AppHandler) webAuthnIndex(ctx *gin.Context) {
	ah.renderTemplate("webauthn.gohtml", ctx, nil)
}

// GET -> USer + params
// called first
// check if user exists
/*
func (ah *AppHandler) BeginRegistration(ctx *gin.Context) {

	username := ctx.Param("userID")
	tempUser := &storage.User{}
	tempUser.Name = username
	options, sessionData, err := ah.web.BeginRegistration(tempUser)

	s := getUserSession(ctx)
	s.WebAuthnSession = sessionData
	s.WebAuthnUser = tempUser

	if err != nil {
		panic(err)
	}

	ctx.JSON(200, options)
}

func (ah *AppHandler) FinishRegistration(ctx *gin.Context) {
	//user := datastore.GetUser() // Get the user
	s := getUserSession(ctx)
	user := s.WebAuthnUser

	// Get the gnolsession data stored from the function above
	// using gorilla/sessions it could look like this
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(ctx.Request.Body)
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
		ctx.JSON(200, command.NewRedirectCommand("/series"))
	} else {
		ctx.JSON(200, "Registration FAILED")
	}

}

// Start of auth
// Check for user in DB
func (ah *AppHandler) BeginLogin(ctx *gin.Context) {
	//user := datastore.GetUser() // Find the user
	user := ah.dao.GetWebAuthnUser(ctx.Param("userID"))
	options, sessionData, err := ah.web.BeginLogin(user)

	if err != nil {
		panic(err)
	}

	getUserSession(ctx).WebAuthnSession = sessionData
	getUserSession(ctx).WebAuthnUser = user
	// handle errors if present
	// store the sessionData values
	ctx.JSON(200, options)
	//JSONResponse(w, options, http.StatusOK) // return the options generated

	// options.publicKey contain our registration options
}

func (ah *AppHandler) FinishLogin(ctx *gin.Context) {
	//user := datastore.GetUser() // Get the user
	us := getUserSession(ctx)
	user := us.WebAuthnUser
	// Get the gnolsession data stored from the function above
	// using gorilla/sessions it could look like this
	//sessionData := store.Get(r, "login-gnolsession")
	sessionData := getUserSession(ctx).WebAuthnSession

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(ctx.Request.Body)
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
	ctx.JSON(200, command.NewRedirectCommand("/series"))
}
*/
