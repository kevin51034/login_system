package controllers

import (
	"context"
	"log"
	"fmt"
	"time"
	"net/http"
	"html/template"


	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/kevin51034/login_system/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var	db *mongo.Database
var Collection *mongo.Collection

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("views/*"))
	dbSessionsCleaned = time.Now()
	Client = Connect()
	Collection = Client.Database("login_system_golang").Collection("users")
	fmt.Println("database connected")
}

func GetAllUser(ctx *gin.Context) {
	fmt.Println("GetUser")
	//collection := Client.Database("login_system_golang").Collection("users")
	fmt.Println(Collection)

	cur, err := Collection.Find(context.Background(), bson.D{})
	if err != nil { log.Fatal(err) }
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
	// To decode into a struct, use cursor.Decode()
	result := struct{
		Foo string
		Bar int32
	}{}
	err := cur.Decode(&result)
	if err != nil { log.Fatal(err) }
	// do something with result...

	// To get the raw bson bytes use cursor.Current
	raw := cur.Current
	// do something with raw...
	fmt.Println(raw)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

}

func Index(ctx *gin.Context) {
	u := getUser(ctx)
	showSessions() // for demonstration purposes
	tpl.ExecuteTemplate(ctx.Writer, "index.gohtml", u)
}

func Bar(ctx *gin.Context) {
	u := getUser(ctx)
	if !alreadyLoggedIn(ctx) {
		http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusSeeOther)
		return
	}
	if u.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Permission denied"})
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	showSessions() // for demonstration purposes
	tpl.ExecuteTemplate(ctx.Writer, "bar.gohtml", u)
}


func Signuppage(ctx *gin.Context) {
	if alreadyLoggedIn(ctx) {
		http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(ctx.Writer, "signup.gohtml", nil)
}


func Signup(ctx *gin.Context) {
	//collection := Client.Database("login_system_golang").Collection("users")
	if alreadyLoggedIn(ctx) {
		http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusSeeOther)
		return
	}
	un := ctx.PostForm("username")
	p := ctx.PostForm("password")
	fn := ctx.PostForm("firstname")
	ln := ctx.PostForm("lastname")
	r := ctx.PostForm("role")

	// check username
	var user models.User
	err := Collection.FindOne(context.TODO(), bson.D{{"username", un}}).Decode(&user)
	if err == nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Username already taken"})
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	bs, errr := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	newuser := models.User{un, bs, fn, ln, r}
	res, errrr := Collection.InsertOne(context.Background(), newuser)
	if errrr != nil { log.Fatal(err) }
	id := res.InsertedID
	fmt.Println(id)

	// create session
	sID, _ := uuid.NewV4()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	http.SetCookie(ctx.Writer, c)
	dbSessions[c.Value] = session{un, time.Now()}
	dbUsers[un] = newuser
	http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusSeeOther)
	return
}

func GetHello(ctx *gin.Context) {
	fmt.Println("hello")
}

func Loginpage(ctx *gin.Context) {
	if alreadyLoggedIn(ctx) {
		http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(ctx.Writer, "login.gohtml", nil)
}

func Login(ctx *gin.Context) {
	fmt.Println("login")
	//tpl.ExecuteTemplate(ctx.Writer, "login.gohtml", nil)

	if alreadyLoggedIn(ctx) {
		http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusSeeOther)
		return
	}
	un := ctx.PostForm("username")
	fmt.Println(un)
	p := ctx.PostForm("password")
	//
	var user models.User //bson.D
	err := Collection.FindOne(context.TODO(), bson.D{{"username", un}}).Decode(&user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Username and/or password do not match"})
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("found document %v", user)
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(p))
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Username and/or password do not match"})
		ctx.AbortWithStatus(http.StatusForbidden)
		return	
	}
	//
	fmt.Println("login succeed")

	sID, _ := uuid.NewV4()
	c := &http.Cookie{
		Name: "session",
		Value: sID.String(),
	}
	c.MaxAge = sessionAge
	http.SetCookie(ctx.Writer, c)
	dbSessions[c.Value] = session{un, time.Now()}

	http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusSeeOther)
	return
}

func Logout(ctx *gin.Context) {
	if !alreadyLoggedIn(ctx) {
		http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusSeeOther)
		return
	}

	c, _ := ctx.Request.Cookie("session")
	delete(dbSessions, c.Value)
	c = &http.Cookie{
		Name: "session",
		Value: "",
		MaxAge: -1,
	}
	http.SetCookie(ctx.Writer, c)
	// clean up dbSessions
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}
	http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusSeeOther)
	return
}