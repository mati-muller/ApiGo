package routes

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupInventarioRoutes(r *gin.Engine) {
	r.GET("/inventario/data", getInventarioData)
	r.GET("/inventario/total", getInventarioTotal)
	r.GET("/inventario/placas", getPlacasData)
	r.GET("/inventario/all", getAllInventario)
	r.POST("/inventario/addplaca", addPlacas)
}

func getInventarioData(c *gin.Context) {
	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Add detailed error logging to capture database issues
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
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

func getInventarioTotal(c *gin.Context) {
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Add detailed error logging to capture database issues
	rows, err := db.Query(`
			SELECT 
				inventario.placa, 
				SUM(inventario.cantidad) AS Cantidad
			FROM 
				inventario
			GROUP BY 
				inventario.placa
		`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
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
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Add detailed error logging to capture database issues
	rows, err := db.Query(`
			SELECT 
				inventario.placa
			FROM 
				inventario
			WHERE 
				inventario.placa LIKE 'PLACA%'
			GROUP BY 
				inventario.placa
		`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		return
	}
	defer rows.Close()

	// Process query results
	var results []map[string]interface{}
	for rows.Next() {
		var placa string
		if err := rows.Scan(&placa); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		results = append(results, gin.H{"placa": placa})
	}

	c.JSON(http.StatusOK, results)
}

func addPlacas(c *gin.Context) {
	// Parse input JSON
	var input struct {
		Placa       string  `json:"placa"`
		Fecha       string  `json:"fecha"`
		PrecioPP    float64 `json:"preciopp"`
		PrecioTotal float64 `json:"precio_total"`
		Cantidad    int     `json:"cantidad"`
		OC          string  `json:"oc"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Get current date in Chile timezone (CLT/CLST)
	loc, _ := time.LoadLocation("America/Santiago")
	currentDate := time.Now().In(loc).Format("02/01/2006") // dd/mm/yyyy format

	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Insert data into the database
	_, err = db.Exec(
		"INSERT INTO inventario (placa, fecha_compra, precio_pp, precio_total, cantidad, oc) VALUES (@p1, @p2, @p3, @p4, @p5, @p6)",
		input.Placa, currentDate, input.PrecioPP, input.PrecioTotal, input.Cantidad, input.OC,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully"})
}

func getAllInventario(c *gin.Context) {
	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Select all records with all columns
	rows, err := db.Query(`
		SELECT 
			placa, 
			fecha_compra, 
			precio_pp, 
			precio_total, 
			cantidad, 
			oc
		FROM 
			inventario
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		return
	}
	defer rows.Close()
	// Process query results
	var results []map[string]interface{}
	for rows.Next() {
		var placa sql.NullString
		var fechaCompra sql.NullString
		var precioPP sql.NullFloat64
		var precioTotal sql.NullFloat64
		var cantidad sql.NullInt64
		var oc sql.NullString

		if err := rows.Scan(&placa, &fechaCompra, &precioPP, &precioTotal, &cantidad, &oc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row", "details": err.Error()})
			return
		}
		results = append(results, gin.H{
			"placa":        placa.String,
			"fecha_compra": fechaCompra.String,
			"precio_pp":    precioPP.Float64,
			"precio_total": precioTotal.Float64,
			"cantidad":     cantidad.Int64,
			"oc":           oc.String,
		})
	}

	c.JSON(http.StatusOK, results)
}
