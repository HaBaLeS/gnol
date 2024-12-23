package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/storage"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (ah *AppHandler) serveTemplate(t string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gctx, ok := ctx.Get("gnol-context")
		if ok {
			ah.renderTemplate(t, ctx, gctx.(*GnolContext))
		} else {
			ah.renderTemplate(t, ctx, nil)
		}

	}
}

func (ah *AppHandler) listUsers(ctx *gin.Context) {
	fmt.Println("Access listUsers")
}

func (ah *AppHandler) logoutUser(ctx *gin.Context) {
	us := getGnolContext(ctx).Session
	ah.dao.DB.MustExec("delete from gnol_session where session_id = $1", us.SessionId)
	ctx.Redirect(302, "/users/login")
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
	us := getGnolContext(ctx).Session
	token := ah.dao.GetOrCreateAPItoken(us.UserId)
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
		ctx.Redirect(303, "/login")
	}
}

func (ah *AppHandler) loginUser(ctx *gin.Context) {
	name := ctx.PostForm("username")
	pass := ctx.PostForm("pass")
	user := ah.dao.AuthUser(name, pass)
	sid := sessions.Default(ctx).Get("gnol-session-id")

	if user == nil {
		gs := &storage.GnolSession{}
		rc := NewGnolContext(gs)
		rc.Flash = "Login Failed"
		ah.renderTemplate("login_user.gohtml", ctx, rc)
		return
	}

	ah.dao.DB.MustExec("delete from gnol_session where session_id = $1", sid)
	ah.dao.DB.MustExec("insert into gnol_session values ($1,$2,$3)", sid, time.Now().Add(24*time.Hour), user.Id)

	ctx.Redirect(303, "/series")
}

// renderIndex
func (ah *AppHandler) webAuthnIndex(ctx *gin.Context) {
	ah.renderTemplate("webauthn.gohtml", ctx, nil)
}
