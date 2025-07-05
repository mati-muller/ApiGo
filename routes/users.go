package routes

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var err error
	db, err = sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func handleError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func SetupUserRoutes(r *gin.Engine) {
	r.POST("/users/register", registerUser)
	r.POST("/users/login", loginUser)
	r.POST("/users/change-password", changePassword)
	r.GET("/users/:id", getUserByID)
}

func registerUser(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nombre   string `json:"nombre"`
		Apellido string `json:"apellido"`
		Rol      string `json:"rol"` // <-- Asegurarse de recibir el rol
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validar rol permitido
	validRoles := map[string]bool{"Superadmin": true, "Admin": true, "Operador": true}
	if !validRoles[user.Rol] {
		handleError(c, http.StatusBadRequest, "Rol inválido")
		return
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// El ID será autoincremental, no se pasa en el INSERT
	query := `INSERT INTO REPORTES.dbo.users (USERNAME, PASSWORD, NOMBRE, APELLIDO, ROL, procesos) 
	OUTPUT INSERTED.ID VALUES (@Username, @Password, @Nombre, @Apellido, @Rol, NULL)`
	var userID int
	err = db.QueryRow(query,
		sql.Named("Username", user.Username),
		sql.Named("Password", hashedPassword),
		sql.Named("Nombre", user.Nombre),
		sql.Named("Apellido", user.Apellido),
		sql.Named("Rol", user.Rol),
	).Scan(&userID)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to insert user: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       userID,
		"username": user.Username,
		"nombre":   user.Nombre,
		"apellido": user.Apellido,
		"rol":      user.Rol,
	})
}

func loginUser(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var user struct {
		ID       int
		Username string
		Password string
		Nombre   string
		Apellido string
		Procesos string
		Rol      string
	}
	err := db.QueryRow(
		"SELECT ID, USERNAME, PASSWORD, NOMBRE, APELLIDO, ISNULL(procesos, ''), ROL FROM REPORTES.dbo.users WHERE USERNAME = @Username",
		sql.Named("Username", credentials.Username),
	).Scan(&user.ID, &user.Username, &user.Password, &user.Nombre, &user.Apellido, &user.Procesos, &user.Rol)
	if err == sql.ErrNoRows {
		handleError(c, http.StatusBadRequest, "Invalid username or password")
		return
	} else if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to query user: "+err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid username or password")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"nombreCompleto": user.Nombre + " " + user.Apellido,
		"procesos":       user.Procesos,
		"rol":            user.Rol,
	})
}

func changePassword(c *gin.Context) {
	var request struct {
		Username    string `json:"username"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if request.Username == "" || request.NewPassword == "" {
		handleError(c, http.StatusBadRequest, "Username and new password are required")
		return
	}

	hashedPassword, err := hashPassword(request.NewPassword)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	query := `UPDATE REPORTES.dbo.users SET PASSWORD = @NewPassword WHERE USERNAME = @Username`
	result, err := db.Exec(query,
		sql.Named("NewPassword", hashedPassword),
		sql.Named("Username", request.Username),
	)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to update password: "+err.Error())
		return
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		handleError(c, http.StatusBadRequest, "User not found or no changes made")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")

	var user struct {
		ID       int
		Username string
		Nombre   string
		Apellido string
		Procesos string
		Rol      string
	}

	err := db.QueryRow(
		"SELECT ID, USERNAME, NOMBRE, APELLIDO, ISNULL(procesos, ''), ROL FROM REPORTES.dbo.users WHERE ID = @ID",
		sql.Named("ID", id),
	).Scan(&user.ID, &user.Username, &user.Nombre, &user.Apellido, &user.Procesos, &user.Rol)
	if err == sql.ErrNoRows {
		handleError(c, http.StatusNotFound, "User not found")
		return
	} else if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to query user: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"nombreCompleto": user.Nombre + " " + user.Apellido,
		"procesos":       user.Procesos,
		"rol":            user.Rol,
	})
}
