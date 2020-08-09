package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/dao"
	"github.com/go-chi/chi"
	"net/http"
)

//SetupUserRouting defines all routes for /user and below.
//this path cares about UserManagement
func (ah *AppHandler) SetupUserRouting() {
	ah.Router.Route("/users", func(r chi.Router) {
		r.Get("/", ah.listUsers())
		r.Post("/", ah.createUser())
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", ah.getUser())
			r.Put("/", ah.updateUser())
			r.Delete("/", ah.deleteUser())
		})
		r.Get("/create", ah.serveTemplate("create_user.gohtml", nil))
		r.Get("/login", ah.serveTemplate("login_user.gohtml", nil))
		r.Post("/login", ah.loginUser())
		r.Get("/logout", ah.logoutUser())
	})
}

func (ah *AppHandler) serveTemplate(t string, data interface{}) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ah.renderTemplate(t, w, r, data)
	})
}

func (ah *AppHandler) listUsers() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access listUsers")
	})
}

func (ah *AppHandler) logoutUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getUserSession(r.Context()).Invalidate()
		http.Redirect(w,r,"/comics",303)
	})
}


func (ah *AppHandler) deleteUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access deleteUser")
	})
}

func (ah *AppHandler) updateUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access updateUser")
	})
}
func (ah *AppHandler) getUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access getUser")
	})
}
func (ah *AppHandler) createUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		user := ah.dao.CreateUser(name, pass)
		us := getUserSession(r.Context())
		us.UserName = user.Name

		http.Redirect(w,r,"/comics",303)
	})
}

func (ah *AppHandler) loginUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("username")
		pass := r.FormValue("pass")
		user, loginerr := ah.dao.AuthUser(name, pass)
		if loginerr != nil {
			ah.renderTemplate("login_user.gohtml", w, r, loginerr)
			return
		}
		us := getUserSession(r.Context())
		us.AuthSession()
		us.UserName = user.Name
		us.UserID = user.Id

		http.Redirect(w,r,"/comics",303)
	})
}

func restoreSetting(user *dao.User) {

}

//activate user

//login user
