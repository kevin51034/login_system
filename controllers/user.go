package controllers

import (
	//"context"
	"log"
	"fmt"
	//"time"
	//"net/http"

	//"github.com/kevin51034/login_system/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

/*
type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}*/

func GetUser(ctx *gin.Context) {
	fmt.Println("GetUser")

	//fmt.Println(Client)

	
	//var user models.User
	//filter := bson.D{{"Name", "Nic"}}
	//db = client.Database("login_system_golang")
	/*
	u1 := models.User{"kevin", []byte("password"), "Kevin", "Chen", "admin"}

	collection = Client.Database("login_system_golang").Collection("users")
	c, _:= context.WithTimeout(context.Background(), 10*time.Second)
	insertResult, err := collection.InsertOne(c, u1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
*/

	collection = Client.Database("login_system_golang").Collection("users")
	/*c, _:= context.WithTimeout(context.Background(), 10*time.Second)	
	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		log.Fatal(err)
	}*/
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cursor)

	fmt.Printf("Found a single document: %+v\n", cursor)

	ctx.JSON(200, gin.H{"message": "Get User successfully!"})
}

func GetHello(ctx *gin.Context) {
	fmt.Println("hello")
}