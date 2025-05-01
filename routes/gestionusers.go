package routes

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupUserDataRoutes(r *gin.Engine) {
	r.GET("/users/data", getUsers)
}

func getUsers(c *gin.Context) {
	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to connect to database: "+err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT ID, NOMBRE, APELLIDO FROM REPORTES.dbo.users")
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to fetch users: "+err.Error())
		return
	}
	defer rows.Close()

	var users []struct {
		ID       int    `json:"id"`
		Nombre   string `json:"nombre"`
		Apellido string `json:"apellido"`
	}

	for rows.Next() {
		var user struct {
			ID       int    `json:"id"`
			Nombre   string `json:"nombre"`
			Apellido string `json:"apellido"`
		}
		if err := rows.Scan(&user.ID, &user.Nombre, &user.Apellido); err != nil {
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
