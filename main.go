package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users []User = []User{
	{Id: "1", Name: "Abhay Jha", Email: "abhay.jha@gmail.com", Password: "pass123"},
	{Id: "2", Name: "John Doe", Email: "john.doe@gmail.com", Password: "pass123"},
}

func main() {
	fmt.Println("lets go user crud api")

	r := gin.Default()

	// healthcheck route
	r.GET("/healthcheck", healthcheck)

	// users route
	r.GET("/users", getUsers)
	r.POST("/user", addUser)
	r.GET("/user/:id", getUserById)
	r.PUT("/user/:id", updateUser)
	r.DELETE("/user/:id", deleteUserById)

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

	if len(users) == 0 {
		c.JSON(http.StatusOK, response(false, nil, "users not found"))
	} else {
		c.JSON(http.StatusOK, response(true, users, "users fetched"))
	}
}

// addUser
func addUser(c *gin.Context) {
	var body User

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, "some error occured"))
		return
	}

	if body.Id == "" || body.Name == "" || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "fields cannot be empty"))
		return
	}

	users = append(users, body)
	c.JSON(http.StatusOK, response(true, nil, "user added."))
}

// getUserById
func getUserById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "user id not provided"))
		return
	}

	for _, user := range users {
		if user.Id == id {
			c.JSON(http.StatusOK, response(true, user, "user found"))
			return
		}
	}

	c.JSON(http.StatusNotFound, response(false, nil, "user not found"))
}

// updateUser
func updateUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "user id not provided"))
		return
	}

	var body User

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, "some error occured"))
		return
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, response(false, body, "fields cannot be empty"))
		return
	}

	for i, user := range users {
		if user.Id == id {
			// user found

			users[i].Name = body.Name
			users[i].Email = body.Email
			users[i].Password = body.Password

			c.JSON(http.StatusOK, response(true, users[i], "user updated."))
			return

		}
	}

	c.JSON(http.StatusNotFound, response(false, nil, "user not found"))
}

// deleteUserById
func deleteUserById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "user id not provided"))
		return
	}

	for i, user := range users {
		if user.Id == id {
			// user found

			users = append(users[0:i], users[i+1:]...)
			c.JSON(http.StatusOK, response(true, nil, "user deleted."))
			return
		}
	}

	c.JSON(http.StatusNotFound, response(false, nil, "user not found"))
}
