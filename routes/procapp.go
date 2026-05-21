package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupProcAppRoutes(r *gin.Engine) {
	r.GET("/app/troquelado", getTroquelado)
	r.GET("/app/troquelado2", getTroquelado2)
	r.GET("/app/encolado", getEncolado)
	r.GET("/app/encolado2", getEncolado2)
	r.GET("/app/multiple", getMultiple)
	r.GET("/app/multiple2", getMultiple2)
	r.GET("/app/pegado", getPegado)
	r.GET("/app/trozado", getTrozado)
	r.GET("/app/impresion", getImpresion)
	r.GET("/app/calado", getCalado)
	r.GET("/app/plizado", getPlizado)
	r.GET("/app/emplacado", getEmplacado)
	r.PUT("/app/procesos/fecha-entrega", updateFechaEntrega)

	// Endpoints de bloqueo de procesos
	r.POST("/app/lock-process", lockProcess)
	r.GET("/app/check-lock/:tableName/:id", checkProcessLock)
}

func queryDatabase(c *gin.Context, query string) {
	log.Println("[DEBUG] queryDatabase: ejecutando query")
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("[ERROR] queryDatabase: error ejecutando query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	log.Printf("[DEBUG] queryDatabase: columnas: %v", columns)
	results := []map[string]interface{}{}

	for rows.Next() {
		row := make([]interface{}, len(columns))
		rowPointers := make([]interface{}, len(columns))
		for i := range row {
			rowPointers[i] = &row[i]
		}

		if err := rows.Scan(rowPointers...); err != nil {
			log.Printf("[ERROR] queryDatabase: error en Scan: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := map[string]interface{}{}
		for i, col := range columns {
			result[col] = row[i]
		}
		results = append(results, result)
		log.Printf("[DEBUG] queryDatabase: row: %v", result)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] queryDatabase: rows.Err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[DEBUG] queryDatabase: total rows: %d", len(results))
	c.JSON(http.StatusOK, results)
}

func getTroquelado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN TROQUELADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func getTroquelado2(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN TROQUELADO2 p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func getEncolado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN ENCOLADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func getEncolado2(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN ENCOLADO2 p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func getMultiple(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN MULTIPLE p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func getMultiple2(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN MULTIPLE2 p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getPegado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN PEGADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getTrozado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN TROZADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getImpresion(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN IMPRESION p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getCalado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN CALADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getPlizado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN PLIZADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getEmplacado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD,
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN EMPLACADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func updateFechaEntrega(c *gin.Context) {
	var reqBody struct {
		ID           int    `json:"ID" binding:"required"`
		FechaEntrega string `json:"FECHA_ENTREGA" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if strings.TrimSpace(reqBody.FechaEntrega) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FECHA_ENTREGA is required"})
		return
	}

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	updateQuery := `UPDATE procesos SET FECHA_ENTREGA = @fecha WHERE ID = @id`
	result, err := db.Exec(updateQuery, sql.Named("fecha", reqBody.FechaEntrega), sql.Named("id", reqBody.ID))
	if err != nil {
		log.Printf("Error updating FECHA_ENTREGA: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating delivery date"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error reading rows affected: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error confirming update"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Process not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Fecha de entrega actualizada",
		"ID":            reqBody.ID,
		"FECHA_ENTREGA": reqBody.FechaEntrega,
	})
}

// checkProcessLock verifica si un proceso está bloqueado
func checkProcessLock(c *gin.Context) {
	tableName := c.Param("tableName")
	id := c.Param("id")

	// Validar nombre de tabla para prevenir SQL injection
	validTables := []string{"TROQUELADO", "TROQUELADO2", "ENCOLADO", "ENCOLADO2", "MULTIPLE", "MULTIPLE2", "PEGADO", "TROZADO", "IMPRESION", "CALADO", "PLIZADO", "EMPLACADO"}
	isValid := false
	for _, valid := range validTables {
		if strings.EqualFold(tableName, valid) {
			tableName = strings.ToUpper(tableName)
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table name"})
		return
	}

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	// Verificar si la columna PROCESO_BLOQUEADO existe
	var columnExists bool
	checkColumnQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_NAME = '%s' AND COLUMN_NAME = 'PROCESO_BLOQUEADO'
	`, tableName)

	err = db.QueryRow(checkColumnQuery).Scan(&columnExists)
	if err != nil {
		log.Printf("Error checking column existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking column"})
		return
	}

	// Si la columna no existe, asumir que no está bloqueado
	if !columnExists {
		c.JSON(http.StatusOK, gin.H{
			"bloqueado": false,
			"mensaje":   "Proceso disponible",
		})
		return
	}

	// Verificar el estado del bloqueo
	var bloqueado sql.NullBool
	query := fmt.Sprintf("SELECT PROCESO_BLOQUEADO FROM %s WHERE ID = @id", tableName)
	err = db.QueryRow(query, sql.Named("id", id)).Scan(&bloqueado)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Process not found"})
		return
	} else if err != nil {
		log.Printf("Error checking lock status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking lock"})
		return
	}

	// Si es NULL o false, no está bloqueado
	isBloqueado := bloqueado.Valid && bloqueado.Bool

	c.JSON(http.StatusOK, gin.H{
		"bloqueado": isBloqueado,
		"mensaje":   map[bool]string{true: "Proceso bloqueado por otro dispositivo", false: "Proceso disponible"}[isBloqueado],
	})
}

// lockProcess bloquea un proceso para evitar que se inicie desde otro dispositivo
func lockProcess(c *gin.Context) {
	var reqBody struct {
		TableName string `json:"tableName" binding:"required"`
		ID        int    `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Validar nombre de tabla
	validTables := []string{"TROQUELADO", "TROQUELADO2", "ENCOLADO", "ENCOLADO2", "MULTIPLE", "MULTIPLE2", "PEGADO", "TROZADO", "IMPRESION", "CALADO", "PLIZADO", "EMPLACADO"}
	isValid := false
	tableName := strings.ToUpper(reqBody.TableName)
	for _, valid := range validTables {
		if tableName == valid {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table name"})
		return
	}

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	// Verificar si la columna existe, si no, crearla
	var columnExists int
	checkColumnQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_NAME = '%s' AND COLUMN_NAME = 'PROCESO_BLOQUEADO'
	`, tableName)

	err = db.QueryRow(checkColumnQuery).Scan(&columnExists)
	if err != nil {
		log.Printf("Error checking column existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking column"})
		return
	}

	// Si la columna no existe, crearla con ALTER TABLE
	if columnExists == 0 {
		alterQuery := fmt.Sprintf("ALTER TABLE %s ADD PROCESO_BLOQUEADO BIT DEFAULT 0", tableName)
		_, err = db.Exec(alterQuery)
		if err != nil {
			log.Printf("Error creating column: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating lock column"})
			return
		}
		log.Printf("Column PROCESO_BLOQUEADO created in table %s", tableName)
	}

	// Verificar si ya está bloqueado
	var bloqueado sql.NullBool
	checkQuery := fmt.Sprintf("SELECT PROCESO_BLOQUEADO FROM %s WHERE ID = @id", tableName)
	err = db.QueryRow(checkQuery, sql.Named("id", reqBody.ID)).Scan(&bloqueado)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Process not found"})
		return
	} else if err != nil {
		log.Printf("Error checking lock: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking lock"})
		return
	}

	// Si ya está bloqueado, retornar error
	if bloqueado.Valid && bloqueado.Bool {
		c.JSON(http.StatusConflict, gin.H{
			"error":     "Proceso ya bloqueado",
			"mensaje":   "Este proceso ya está siendo usado en otro dispositivo",
			"bloqueado": true,
		})
		return
	}

	// Bloquear el proceso
	updateQuery := fmt.Sprintf("UPDATE %s SET PROCESO_BLOQUEADO = 1 WHERE ID = @id", tableName)
	_, err = db.Exec(updateQuery, sql.Named("id", reqBody.ID))
	if err != nil {
		log.Printf("Error locking process: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error locking process"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Proceso bloqueado exitosamente",
		"bloqueado": true,
	})
}
