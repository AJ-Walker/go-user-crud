package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	fmt.Println("lets go user crud api")

	// env load
	if err := godotenv.Load(); err != nil {
		fmt.Println("env err: ", err)
		return
	}

	// db conn
	if err := DBConnAndPing(); err != nil {
		fmt.Println("db err: ", err)
		return
	}

	// gin router
	r := gin.Default()

	// healthcheck route
	r.GET("/healthcheck", healthcheck)

	// users route
	r.GET("/users", getUsers)
	r.POST("/user", addUser)
	r.GET("/user/:id", getUserById)
	r.PUT("/user/:id", updateUser)
	r.DELETE("/user/:id", deleteUserById)

	// auth
	r.POST("/login", login)

	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}

// healthcheck func
func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, response(true, nil, "healthcheck works.. :)"))
}

// response common func
func response(status bool, data any, message string) gin.H {
	return gin.H{"status": status, "data": data, "message": message}
}

// getUsers
func getUsers(c *gin.Context) {

	res, err := GetUsersDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("GetUsersDB err: %v", err)))
		return
	}

	if len(res) == 0 {
		c.JSON(http.StatusOK, response(false, nil, "users not found"))
	} else {
		c.JSON(http.StatusOK, response(true, res, "users fetched"))
	}
}

// addUser
func addUser(c *gin.Context) {
	var body User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, "some error occured"))
		return
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "fields cannot be empty"))
		return
	}

	// unique id
	body.Id = uuid.NewV4().String()

	// hashed password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, "password hasing gone wrong"))
		return
	}

	body.Password = string(hashedPass)

	// check if email already exists
	resEmail, errEmail := GetUserByEmailDB(body.Email)
	if errEmail != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("add user err: %v", err)))
		return
	}

	if resEmail.Email == body.Email {
		c.JSON(http.StatusBadRequest, response(false, nil, "email already exists"))
		return
	}

	if err := AddUserDB(body); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("add user err: %v", err)))
		return
	}

	res, err := GetUserByIdDB(body.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, response(false, nil, fmt.Sprintf("add user err: %v", err)))
		return
	}

	c.JSON(http.StatusOK, response(true, res, "user added."))
}

// getUserById
func getUserById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "user id not provided"))
		return
	}

	res, err := GetUserByIdDB(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response(false, nil, "user not found"))
		return
	}

	c.JSON(http.StatusOK, response(true, res, "user found"))
}

// updateUser
func updateUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "user id not provided"))
		return
	}

	var body User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, "some error occured"))
		return
	}

	if body.Name == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "fields cannot be empty"))
		return
	}

	if err := UpdateUserDB(id, body); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("update user err: %v", err)))
		return
	}

	res, err := GetUserByIdDB(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response(false, nil, "user not found."))
		return
	}

	c.JSON(http.StatusOK, response(true, res, "user updated."))
}

// deleteUserById
func deleteUserById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "user id not provided"))
		return
	}

	_, err := GetUserByIdDB(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response(false, nil, "user not found"))
		return
	}

	res, err := DeleteUserDB(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("delete user err: %v", err)))
		return
	}

	if res {
		c.JSON(http.StatusOK, response(true, nil, "user deleted."))
		return
	}
	c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("delete user err: %v", err)))
}

// login
func login(c *gin.Context) {
	var loginBody Login

	if err := c.ShouldBindJSON(&loginBody); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, "some error occured."))
		return
	}

	if loginBody.Email == "" || loginBody.Password == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "fields cannot be empty"))
		return
	}

	res, err := GetUserByEmailDB(loginBody.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, response(false, nil, "user not found"))
		return
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(loginBody.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, response(false, nil, "unauthorized user."))
		return
	}

	// if correct password
	c.JSON(http.StatusOK, response(true, nil, "login success"))
}
