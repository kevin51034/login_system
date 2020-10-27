package main

import (
	"net/http"
	"html/template"
	"log"
	//"golang.org/x/crypto/bcrypt"
	"github.com/satori/go.uuid"
)

type user struct {
	UserName string
	Password string
	First    string
	Last     string
}

// *Template in template package
var tpl *template.Template
var dbUsers = map[string]user{}      // map[user ID] user (struct)
var dbSessions = map[string]string{} // map[session ID] user ID


func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	// generate a test user
	//bs, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	//dbUsers["test@test.com"] = user{"test@test.com", bs, "Chen", "Kevin"}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/bar", bar)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(req)
	tpl.ExecuteTemplate(w, "index.gohtml", u)
}

func bar(w http.ResponseWriter, req *http.Request) {
	u := getUser(req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}

func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	// process the form submission
	if req.Method == http.MethodPost {
		// get form values
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")

		// check username
		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			// need to return
			return
		}

		// create session
		sID, _ := uuid.NewV4()
		// set cookie
		c := &http.Cookie{
			Name: "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un

		// store user in dbUsers
		u := user{un, p, f, l}
		dbUsers[un] = u

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}