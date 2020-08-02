package router

import (
	"fmt"
	"github.com/HaBaLeS/gnol/server/dao"
	"github.com/go-chi/chi"
	"net/http"
)

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
	})
}

func (ah *AppHandler) serveTemplate(t string, data interface{}) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl, err := ah.getTemplate(t)
		if err != nil {
			panic(err)
		}
		err2 := renderTemplate(tpl, w, r, data)
		if err2 != nil {
			panic(err2)
		}
	})
}

func (ah *AppHandler) listUsers() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access listUsers")
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
			tpl, err := ah.getTemplate("create_user.gohtml")
			if err != nil {
				panic(err)
			}
			err2 := renderTemplate(tpl, w, r, est)
			if err2 != nil {
				panic(err2)
			}
		}

		//TODO check for username

		user := ah.dao.CreateUser(name, pass)
		us := getUserSession(r.Context())
		us.UserName = user.Name

	})
}

func (ah *AppHandler) loginUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("username")
		pass := r.FormValue("pass")
		user, loginerr := ah.dao.AuthUser(name, pass)
		if loginerr != nil {
			//FIXME extract render templete method
			tpl, err := ah.getTemplate("login_user.gohtml")
			if err != nil {
				panic(err)
			}
			err2 := renderTemplate(tpl, w, r, loginerr)
			if err2 != nil {
				panic(err2)
			}
			//FIXME to here
			return
		}
		us := getUserSession(r.Context())
		us.UserName = user.Name
	})
}

func restoreSetting(user *dao.User) {

}

//activate user

//login user
