package router

import (
	"fmt"
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

		user := ah.bs.User.CreateUser(name, pass)
		us := getUserSession(r.Context())
		us.UserName = user.Name

		http.Redirect(w,r,"/comics",303)
	}
}

func (ah *AppHandler) loginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("username")
		pass := r.FormValue("pass")
		user, loginerr := ah.bs.User.AuthUser(name, pass)
		if loginerr != nil {
			ah.renderTemplate("login_user.gohtml", w, r, loginerr)
			return
		}
		us := getUserSession(r.Context())
		us.AuthSession()
		us.UserName = user.Name
		us.UserID = string(user.IdBytes())

		http.Redirect(w,r,"/comics",303)
	}
}