package api

import (
 	"context"
 	"net/http"
 	"time"


 	"github.com/gin-gonic/gin"
 	"golang.org/x/crypto/bcrypt"

)


func SignUp() gin.HandlerFunc {
 	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
  		defer cancel()

  		var creds struct {
   			Username 	string 		`json:"username"`
   			Password 	string 		`json:"password"`
  		}

  		if err := c.ShouldBindJSON(&creds); err != nil {
   			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
   			return
  		}

  		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
  		if err != nil {
   			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
   			return
  		}


  		query := "INSERT INTO users (username, password) VALUES ($1, $2)"
  		_, err = conn.Exec(ctx, query, creds.Username, string(hashedPassword))
  		if err != nil {
   			c.JSON(http.StatusConflict, gin.H{"error": "Email уже зарегистрирован"})
   			return
  		}

  		c.JSON(http.StatusCreated, gin.H{"message": "Регистрация завершена"})
 	}
}

func Login() gin.HandlerFunc {
 	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
  		defer cancel()
  		var creds struct {
   			Username 	string 		`json:"username"`
   			Password 	string 		`json:"password"`
  		}

  		if err := c.ShouldBindJSON(&creds); err != nil {
   			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
   			return
  		}


 	 	var storedPassword string
  		query := "SELECT password FROM users WHERE username = $1"
  		err := conn.QueryRow(ctx, query, creds.Username).Scan(&storedPassword)
  		if err != nil {
    		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
   			return
  		}

  		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
   			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неправильный пароль"})
   			return
  		}

  		token, err := GenerateJWT(creds.Username)
  		if err != nil {
   			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
   			return
 		}

  		c.SetCookie("token", token, 3600*24, "/", "", true, true)

  		c.JSON(http.StatusOK, gin.H{"message": "Успешный вход"})
 	}
}

func LogoutUser() gin.HandlerFunc {
 	return func(c *gin.Context) {
  		c.SetCookie("token", "", -1, "/", "", true, true)
  		c.JSON(http.StatusOK, gin.H{"message": "Успешный выход"})
 	}
}