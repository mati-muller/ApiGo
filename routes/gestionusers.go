package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupUserDataRoutes(r *gin.Engine) {
	r.GET("/users/data", getUsers)
	r.POST("/users/procesos", procesosUser)
	r.POST("/users/delete", deleteUser)
}

func getUsers(c *gin.Context) {
	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to connect to database: "+err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT ID, NOMBRE, APELLIDO, USERNAME, ISNULL(procesos, '[]') as procesos FROM REPORTES.dbo.users")
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to fetch users: "+err.Error())
		return
	}
	defer rows.Close()

	var users []struct {
		ID       int    `json:"id"`
		Nombre   string `json:"nombre"`
		Apellido string `json:"apellido"`
		Username string `json:"username"`
		Procesos string `json:"procesos"` // Added field for processes
	}

	for rows.Next() {
		var user struct {
			ID       int    `json:"id"`
			Nombre   string `json:"nombre"`
			Apellido string `json:"apellido"`
			Username string `json:"username"`
			Procesos string `json:"procesos"` // Added field for processes
		}
		if err := rows.Scan(&user.ID, &user.Nombre, &user.Apellido, &user.Username, &user.Procesos); err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to scan user: "+err.Error())
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		handleError(c, http.StatusInternalServerError, "Error iterating over users: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func procesosUser(c *gin.Context) {
	// Parse request body
	var requestBody struct {
		UserID   int      `json:"user_id"`
		Procesos []string `json:"procesos"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to connect to database: "+err.Error())
		return
	}
	defer db.Close()

	// Convert procesos to JSON string
	procesosJSON, err := json.Marshal(requestBody.Procesos)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to encode procesos: "+err.Error())
		return
	}

	// Update the user's procesos in the database
	query := "UPDATE REPORTES.dbo.users SET procesos = @procesos WHERE ID = @userID"
	_, err = db.Exec(query, sql.Named("procesos", string(procesosJSON)), sql.Named("userID", requestBody.UserID))
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to update procesos: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Procesos updated successfully"})
}

func deleteUser(c *gin.Context) {
	// Parse request body
	var requestBody struct {
		UserID int `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to connect to database: "+err.Error())
		return
	}
	defer db.Close()

	// Delete the user from the database
	query := "DELETE FROM REPORTES.dbo.users WHERE ID = @userID"
	_, err = db.Exec(query, sql.Named("userID", requestBody.UserID))
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
