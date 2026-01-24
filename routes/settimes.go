package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupSetupTimesRoutes(r *gin.Engine) {
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create table if not exists
	createTableSQL := `
	IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'dbo' AND TABLE_NAME = 'SETUP_TIMES')
	BEGIN
		CREATE TABLE dbo.SETUP_TIMES (
			ID INT PRIMARY KEY IDENTITY(1,1),
			ID_PROCESO INT DEFAULT NULL,
			NVNUMERO NVARCHAR(MAX) DEFAULT NULL,
			PRODUCTO NVARCHAR(MAX) DEFAULT NULL,
			PROCESO NVARCHAR(MAX) DEFAULT NULL,
			TIEMPO_MINUTOS DECIMAL(10,2) DEFAULT 0,
			USUARIO NVARCHAR(MAX) DEFAULT NULL,
			FECHA DATETIME DEFAULT GETDATE()
		);
	END
	`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Printf("Warning: Could not ensure SETUP_TIMES table exists: %v", err)
	}

	// POST endpoint to save setup time
	r.POST("/app/setup-time", func(c *gin.Context) {
		var payload struct {
			IDProceso     int         `json:"id_proceso"`
			NVNumero      interface{} `json:"nv_numero"` // Aceptar número o string
			Producto      string      `json:"producto"`
			Proceso       string      `json:"proceso"`
			TiempoMinutos float64     `json:"tiempo_minutos"`
			Usuario       string      `json:"usuario"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON payload",
				"details": err.Error(),
			})
			return
		}

		// Convertir NVNumero a string si es necesario
		nvNumeroStr := ""
		switch v := payload.NVNumero.(type) {
		case string:
			nvNumeroStr = v
		case float64:
			nvNumeroStr = fmt.Sprintf("%.0f", v) // Convertir número a string sin decimales
		case int:
			nvNumeroStr = fmt.Sprintf("%d", v)
		default:
			nvNumeroStr = fmt.Sprintf("%v", v)
		}

		// Validate required fields
		if payload.IDProceso == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID_PROCESO es requerido",
			})
			return
		}

		// Validar que el tiempo de seteo sea mayor a 0
		if payload.TiempoMinutos == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "No se envió el tiempo de seteo",
				"message":   "Por favor, reintente enviando el tiempo de seteo en minutos",
				"reintentar": true,
			})
			return
		}

		// Insert into SETUP_TIMES table
		_, err := db.Exec(
			`INSERT INTO dbo.SETUP_TIMES (ID_PROCESO, NVNUMERO, PRODUCTO, PROCESO, TIEMPO_MINUTOS, USUARIO, FECHA)
			 VALUES (@IDProceso, @NVNumero, @Producto, @Proceso, @TiempoMinutos, @Usuario, GETDATE())`,
			sql.Named("IDProceso", payload.IDProceso),
			sql.Named("NVNumero", nvNumeroStr),
			sql.Named("Producto", payload.Producto),
			sql.Named("Proceso", payload.Proceso),
			sql.Named("TiempoMinutos", payload.TiempoMinutos),
			sql.Named("Usuario", payload.Usuario),
		)

		if err != nil {
			log.Printf("Failed to insert into SETUP_TIMES: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":      "Error al guardar el tiempo de seteo",
				"message":    "No se pudo completar la operación. Por favor, haga clic en Reintentar.",
				"reintentar": true,
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":        "Setup time saved successfully",
			"tiempo_minutos": payload.TiempoMinutos,
			"proceso":        payload.Proceso,
			"usuario":        payload.Usuario,
		})
	})

	// GET endpoint to retrieve setup times (optional, for viewing history)
	r.GET("/app/setup-times/:id_proceso", func(c *gin.Context) {
		idProceso := c.Param("id_proceso")

		rows, err := db.Query(
			`SELECT ID, ID_PROCESO, NVNUMERO, PRODUCTO, PROCESO, TIEMPO_MINUTOS, USUARIO, FECHA 
			 FROM dbo.SETUP_TIMES WHERE ID_PROCESO = @IDProceso ORDER BY FECHA DESC`,
			sql.Named("IDProceso", idProceso),
		)

		if err != nil {
			log.Printf("Failed to query SETUP_TIMES: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve setup times",
			})
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		for rows.Next() {
			var id, idProceso int
			var tiempoMinutos float64
			var nvNumero, producto, proceso, usuario, fecha string

			err := rows.Scan(&id, &idProceso, &nvNumero, &producto, &proceso, &tiempoMinutos, &usuario, &fecha)
			if err != nil {
				log.Printf("Failed to scan row: %v", err)
				continue
			}

			results = append(results, map[string]interface{}{
				"id":             id,
				"id_proceso":     idProceso,
				"nv_numero":      nvNumero,
				"producto":       producto,
				"proceso":        proceso,
				"tiempo_minutos": tiempoMinutos,
				"usuario":        usuario,
				"fecha":          fecha,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"data": results,
		})
	})
}
