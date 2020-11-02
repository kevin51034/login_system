package controllers

import (
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevin51034/login_system/models"
)

type session struct {
	un       string
	lastActivity time.Time
}

var dbUsers = map[string]models.User{}      // map[user ID] user (struct)
var dbSessions = map[string]session{} // map[session ID] user ID
var dbSessionsCleaned time.Time
const sessionAge int = 30

func getUser(ctx *gin.Context) models.User {
	fmt.Println("getUser")
	// get cookie
	c, err := ctx.Request.Cookie("session")
	if err != nil {
		fmt.Println("session not found")
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name: "session",
			Value: sID.String(),
		}
	}
	c.MaxAge = sessionAge
	http.SetCookie(ctx.Writer, c)
	fmt.Println("session time update")

	//ctx.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)


	// if the user exists already, get user
	var u models.User
	if s, ok := dbSessions[c.Value]; ok {
		fmt.Println("session active update")
		s.lastActivity = time.Now()
		dbSessions[c.Value] = s
		u = dbUsers[s.un]
	}
	fmt.Println("return user")
	fmt.Println(u)
	return u
}

func alreadyLoggedIn(ctx *gin.Context) bool {
	c, err := ctx.Request.Cookie("session")
	if err != nil {
		return false
	}
	s, ok := dbSessions[c.Value]
	if ok {
		s.lastActivity = time.Now()
		dbSessions[c.Value] = s
	}

	_, ok = dbUsers[s.un]
	c.MaxAge = sessionAge
	http.SetCookie(ctx.Writer, c)
	return ok
}

func cleanSessions() {
	fmt.Println("BEFORE CLEAN") // for demonstration purposes
	showSessions()              // for demonstration purposes
	for k, v := range dbSessions {
		if time.Now().Sub(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
	fmt.Println("AFTER CLEAN") // for demonstration purposes
	showSessions()             // for demonstration purposes
}

// for demonstration purposes
func showSessions() {
	fmt.Println("******** dbSessions")
	for k, v := range dbSessions {
		fmt.Println(k, v.un)
	}
	fmt.Println("")
	fmt.Println("******** dbUsers")
	for k, v := range dbUsers {
		fmt.Println(k, v)
	}
	fmt.Println("********")

}
