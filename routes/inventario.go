package routes

import (
	"database/sql"
	"fmt"
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
	r.GET("/inventario/oc", getAllOC)
	r.GET("/inventario/oc/:oc", getOCItems)
	r.PUT("/inventario/:id", updateInventarioItem)
	r.GET("/inventario/export/csv", exportInventarioCSV)
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
			id,
			placa, 
			fecha_compra, 
			precio_pp, 
			precio_total, 
			cantidad, 
			oc
		FROM 
			inventario
		Where cantidad > 0
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		return
	}
	defer rows.Close()
	// Process query results
	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var placa sql.NullString
		var fechaCompra sql.NullString
		var precioPP sql.NullFloat64
		var precioTotal sql.NullFloat64
		var cantidad sql.NullInt64
		var oc sql.NullString
		if err := rows.Scan(&id, &placa, &fechaCompra, &precioPP, &precioTotal, &cantidad, &oc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row", "details": err.Error()})
			return
		}
		results = append(results, gin.H{
			"id":           id,
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

// getAllOC obtiene todas las órdenes de compra únicas
func getAllOC(c *gin.Context) {
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT DISTINCT 
			oc,
			COUNT(*) AS cantidad_items
		FROM 
			inventario
		WHERE 
			oc IS NOT NULL AND oc != ''
		GROUP BY 
			oc
		ORDER BY 
			oc DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var oc sql.NullString
		var cantidadItems int

		if err := rows.Scan(&oc, &cantidadItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row", "details": err.Error()})
			return
		}
		results = append(results, gin.H{
			"oc":             oc.String,
			"cantidad_items": cantidadItems,
		})
	}

	c.JSON(http.StatusOK, results)
}

// getOCItems obtiene todos los items de una orden de compra específica
func getOCItems(c *gin.Context) {
	oc := c.Param("oc")

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

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
		WHERE 
			oc = @p1
		ORDER BY 
			placa
	`, oc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var placa sql.NullString
		var fechaCompra sql.NullString
		var precioPP sql.NullFloat64
		var precioTotal sql.NullFloat64
		var cantidad sql.NullInt64
		var ocValue sql.NullString

		if err := rows.Scan(&placa, &fechaCompra, &precioPP, &precioTotal, &cantidad, &ocValue); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row", "details": err.Error()})
			return
		}
		results = append(results, gin.H{
			"placa":        placa.String,
			"fecha_compra": fechaCompra.String,
			"precio_pp":    precioPP.Float64,
			"precio_total": precioTotal.Float64,
			"cantidad":     cantidad.Int64,
			"oc":           ocValue.String,
		})
	}

	c.JSON(http.StatusOK, results)
}

// updateInventarioItem actualiza un item específico por su ID
func updateInventarioItem(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Cantidad    int     `json:"cantidad" binding:"required"`
		PrecioPP    float64 `json:"precio_pp"`
		PrecioTotal float64 `json:"precio_total"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Recalcular precio_total = cantidad * precio_pp
	precioTotal := float64(input.Cantidad) * input.PrecioPP

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	result, err := db.Exec(
		`UPDATE inventario 
		 SET cantidad = @p1, precio_pp = @p2, precio_total = @p3
		 WHERE id = @p4`,
		input.Cantidad, input.PrecioPP, precioTotal, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data", "details": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows affected", "details": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found", "id": id})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "Item updated successfully",
		"id":            id,
		"cantidad":      input.Cantidad,
		"precio_pp":     input.PrecioPP,
		"precio_total":  precioTotal,
		"rows_affected": rowsAffected,
	})
}

func exportInventarioCSV(c *gin.Context) {
	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()
	// Query to get inventory grouped by placa with calculated totals (handling NULL values)
	rows, err := db.Query(`
		SELECT 
			placa,
			SUM(ISNULL(cantidad, 0)) AS cantidad_total,
			MAX(ISNULL(precio_pp, 0)) AS precio_pp,
			SUM(ISNULL(CAST(cantidad AS float), 0) * ISNULL(precio_pp, 0)) AS valor_total
		FROM inventario
		GROUP BY placa
		ORDER BY placa
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		return
	}
	defer rows.Close()

	// Set response headers for CSV download
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=inventario_export.csv")

	// Write CSV header
	c.Writer.WriteString("Placa,Cantidad Total,Precio PP,Valor Total\n")

	// Write CSV rows
	for rows.Next() {
		var placa string
		var cantidadTotal int
		var precioPP, valorTotal float64

		if err := rows.Scan(&placa, &cantidadTotal, &precioPP, &valorTotal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row", "details": err.Error()})
			return
		}

		// Format and write CSV row
		csvLine := fmt.Sprintf("%s,%d,%.2f,%.2f\n", placa, cantidadTotal, precioPP, valorTotal)
		c.Writer.WriteString(csvLine)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading rows", "details": err.Error()})
		return
	}
}
