package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

// User struct to map MongoDB documents
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// MongoDB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	userCollection = client.Database("testdb").Collection("users")

	// Set up Gin router
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// Register routes
	router.POST("/register", registerUser)
	router.POST("/login", loginUser)
	router.GET("/welcome", welcomePage)

	// Serve the frontend files
	router.StaticFile("/", "./index.html")

	router.Run(":8080")
}

func registerUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input"})
		return
	}

	// Check if user already exists
	var existingUser User
	err := userCollection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Username already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error hashing password"})
		return
	}

	// Store user in MongoDB
	user.Password = string(hashedPassword)
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error registering user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Registration successful"})
}

func loginUser(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input"})
		return
	}

	var user User
	err := userCollection.FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid username or password"})
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successful"})
}

func welcomePage(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the dashboard!")
}
