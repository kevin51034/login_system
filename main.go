package main

import (
	//"net/http"
	//"html/template"
	//"log"

	//"time"

	//"github.com/kevin51034/login_system/models"
	"github.com/kevin51034/login_system/controllers"

	"github.com/gin-gonic/gin"
)


/*
type session struct {
	un       string
	lastActivity time.Time
}

// *Template in template package
var tpl *template.Template
var dbUsers = map[string]models.User{}      // map[user ID] user (struct)
var dbSessions = map[string]session{} // map[session ID] user ID
var dbSessionsCleaned time.Time
const sessionAge int = 30
*/

func init() {
	//tpl = template.Must(template.ParseGlob("views/*"))
	//dbSessionsCleaned = time.Now()
	//controllers.Initdb()
	// generate a test user
	// bs, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	// dbUsers["test@test.com"] = user{"test@test.com", bs, "Chen", "Kevin"}
}

func main() {
	/* replace with GIN router
	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/bar", bar)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
*/
	// use GIN router
	router := gin.Default()
	router.GET("/", controllers.Index)
	router.GET("bar", controllers.Bar)
	router.GET("/user", controllers.GetAllUser)
	//router.GET("/hello", controllers.GetHello)
	router.GET("/signup", controllers.Signuppage)
	router.POST("/signup", controllers.Signup)
	router.GET("/login", controllers.Loginpage)
	router.POST("/login", controllers.Login)
	router.GET("/logout", controllers.Logout)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}