package routes

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupInventarioRoutes(r *gin.Engine) {
	r.GET("/inventario/data", getInventarioData)
	r.GET("/inventario/placas", getPlacasData)
}

func getInventarioData(c *gin.Context) {
	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Ensure the SQL query sums quantities for all identical 'placa' values
	rows, err := db.Query(`
			SELECT 
				inventario.placa, 
				SUM(inventario.cantidad) AS Cantidad
			FROM 
				inventario
			WHERE 
				inventario.placa LIKE 'PLACA%'
			GROUP BY 
				inventario.placa
		`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer rows.Close()

	// Process query results
	var results []map[string]interface{}
	for rows.Next() {
		var placa string
		var cantidad int
		if err := rows.Scan(&placa, &cantidad); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		results = append(results, gin.H{"placa": placa, "cantidad": cantidad})
	}

	c.JSON(http.StatusOK, results)
}

func getPlacasData(c *gin.Context) {
	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	rows, err := db.Query(`
			SELECT 
				inventario.placa, 
				SUM(inventario.cantidad) AS Cantidad
			FROM 
				inventario
			WHERE 
				inventario.placa LIKE 'PLACA%'
			GROUP BY 
				inventario.placa
		`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer rows.Close()

	// Process query results
	var results []map[string]interface{}
	for rows.Next() {
		var placa string
		var cantidad int
		if err := rows.Scan(&placa, &cantidad); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		results = append(results, gin.H{"placa": placa, "cantidad": cantidad})
	}

	c.JSON(http.StatusOK, results)
}
