package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

type CustomClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func main() {
	// Set log flags (adds timestamp)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Print("lets go user crud api")

	// env load
	if err := godotenv.Load(); err != nil {
		log.Fatalf("env err: %v", err)
	}

	// db conn
	if err := DBConnAndPing(); err != nil {
		log.Fatal(err)
	}

	// gin router
	r := gin.Default()

	// healthcheck route
	r.GET("/healthcheck", healthcheck)

	// auth
	r.POST("/login", login)

	userRouter := r.Group("/users")
	userRouter.Use(AuthMiddleware())
	// users route
	userRouter.GET("", getUsers)
	userRouter.POST("", addUser)
	userRouter.GET("/:id", getUserById)
	userRouter.PUT("/:id", updateUser)
	userRouter.DELETE("/:id", deleteUserById)

	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}

// AuthMiddleware handles JWT token validation
func AuthMiddleware() gin.HandlerFunc {
	log.Print("Inside AuthMiddleware")
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Print(authHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, response(false, nil, "Authorization header missing."))
			return
		}

		if strings.Split(authHeader, " ")[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, response(false, nil, "Wrong type of authorization header."))
			return
		}

		jwtToken := strings.Split(authHeader, " ")[1]
		log.Print(jwtToken)

		validToken, err := validateJwtToken(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response(false, nil, fmt.Sprintf("%v", err)))
			return
		}

		if claims, ok := validToken.Claims.(*CustomClaims); ok {
			log.Printf("User Id: %v", claims.UserId)
			log.Printf("Email: %v", claims.Issuer)

			user, err := GetUserByIdDB(claims.UserId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusNotFound, response(false, nil, "user not found"))
				return
			}

			if user.Email != claims.Issuer {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response(false, nil, "wrong data provided."))
				return
			}
		}

		c.Next()

	}
}

// healthcheck func
func healthcheck(c *gin.Context) {
	log.Print("Inside healthcheck")
	c.JSON(http.StatusOK, response(true, nil, "healthcheck works.. :)"))
}

// response common func
func response(status bool, data any, message string) gin.H {
	return gin.H{"status": status, "data": data, "message": message}
}

// Generate JWT token
func generateJwtToken(user User) (string, error) {
	log.Print("Inside generateJwtToken")
	secretKey := os.Getenv("JWT_SECRET_KEY")

	claims := CustomClaims{
		user.Id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // token valid for 1 hour
			Issuer:    user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// Validate JWT token
func validateJwtToken(jwtToken string) (*jwt.Token, error) {
	log.Print("Inside validateJwtToken")
	secretKey := os.Getenv("JWT_SECRET_KEY")

	token, err := jwt.ParseWithClaims(jwtToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("bad signed method received")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("validateJwtToken err: %v", err)
		return nil, err
	}
	log.Printf("validateJwtToken: token valid")
	return token, nil
}

// getUsers
func getUsers(c *gin.Context) {
	log.Print("Inside getUsers")
	res, err := GetUsersDB()
	if err != nil {
		log.Printf("getUsers: %v", err)
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("getUsers: %v", err)))
		return
	}

	if len(res) == 0 {
		log.Printf("Users list is empty")
		c.JSON(http.StatusOK, response(false, nil, "users not found"))
	} else {
		log.Printf("users found")
		c.JSON(http.StatusOK, response(true, res, "users fetched"))
	}
}

// addUser
func addUser(c *gin.Context) {
	log.Print("Inside addUser")
	var body User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("addUser: %v", err)))
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
		c.JSON(http.StatusBadRequest, response(false, nil, "password hashing gone wrong"))
		return
	}

	body.Password = string(hashedPass)

	// check if email already exists
	resEmail, errEmail := GetUserByEmailDB(body.Email)
	if errEmail != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("addUser: %v", err)))
		return
	}

	if resEmail.Email == body.Email {
		c.JSON(http.StatusBadRequest, response(false, nil, "email already exists"))
		return
	}

	if err := AddUserDB(body); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("addUser: %v", err)))
		return
	}

	res, err := GetUserByIdDB(body.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, response(false, nil, fmt.Sprintf("addUser: %v", err)))
		return
	}

	c.JSON(http.StatusCreated, response(true, res, "user added."))
}

// getUserById
func getUserById(c *gin.Context) {
	log.Print("Inside getUserById")
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
	log.Print("Inside updateUser")
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "user id not provided"))
		return
	}

	var body User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("updateUser: %v", err)))
		return
	}

	if body.Name == "" {
		c.JSON(http.StatusBadRequest, response(false, nil, "fields cannot be empty"))
		return
	}

	if err := UpdateUserDB(id, body); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("updateUser: %v", err)))
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
	log.Print("Inside deleteUserById")
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
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("deleteUserById: %v", err)))
		return
	}

	if res {
		c.JSON(http.StatusOK, response(true, nil, "user deleted."))
		return
	}
	c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("deleteUserById: %v", err)))
}

// login
func login(c *gin.Context) {
	log.Print("Inside login")
	var loginBody Login

	if err := c.ShouldBindJSON(&loginBody); err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("login: %v", err)))
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

	jwtToken, err := generateJwtToken(res)
	if err != nil {
		c.JSON(http.StatusBadRequest, response(false, nil, fmt.Sprintf("token error: %v", err)))
		return
	}

	data := make(map[string]string, 0)
	data["token"] = jwtToken

	// if correct password
	c.JSON(http.StatusOK, response(true, data, "login success"))
}
