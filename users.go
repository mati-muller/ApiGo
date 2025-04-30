package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SetupUserRoutes(r *gin.Engine) {
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r.POST("/register", func(c *gin.Context) {
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Nombre   string `json:"nombre"`
			Apellido string `json:"apellido"`
			Rol      string `json:"rol"`
		}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if user.Username == "" || user.Password == "" || user.Nombre == "" || user.Apellido == "" || user.Rol == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		query := `INSERT INTO REPORTES.dbo.users (USERNAME, PASSWORD, NOMBRE, APELLIDO, ROL, procesos) 
				OUTPUT INSERTED.ID VALUES (@Username, @Password, @Nombre, @Apellido, @Rol, NULL)`
		var userID int
		err = db.QueryRow(query,
			sql.Named("Username", user.Username),
			sql.Named("Password", string(hashedPassword)),
			sql.Named("Nombre", user.Nombre),
			sql.Named("Apellido", user.Apellido),
			sql.Named("Rol", user.Rol),
		).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":       userID,
			"username": user.Username,
			"nombre":   user.Nombre,
			"apellido": user.Apellido,
			"rol":      user.Rol,
		})
	})

	r.POST("/login", func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		var user struct {
			Username string
			Password string
			Nombre   string
			Apellido string
		}
		err := db.QueryRow(
			"SELECT USERNAME, PASSWORD, NOMBRE, APELLIDO FROM REPORTES.dbo.users WHERE USERNAME = @Username",
			sql.Named("Username", credentials.Username),
		).Scan(&user.Username, &user.Password, &user.Nombre, &user.Apellido)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user: " + err.Error()})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username":       user.Username,
			"nombreCompleto": user.Nombre + " " + user.Apellido,
		})
	})

	r.POST("/change-password", func(c *gin.Context) {
		var request struct {
			Username    string `json:"username"`
			NewPassword string `json:"newPassword"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if request.Username == "" || request.NewPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and new password are required"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		query := `UPDATE REPORTES.dbo.users SET PASSWORD = @NewPassword WHERE USERNAME = @Username`
		result, err := db.Exec(query,
			sql.Named("NewPassword", string(hashedPassword)),
			sql.Named("Username", request.Username),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password: " + err.Error()})
			return
		}

		affected, err := result.RowsAffected()
		if err != nil || affected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found or no changes made"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
	})
}
