package routes

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

type alertItem struct {
	IDProceso        int     `json:"id_proceso"`
	Proceso          string  `json:"proceso"`
	NVNumero         string  `json:"nv_numero"`
	Cliente          string  `json:"cliente"`
	Producto         string  `json:"producto"`
	Placa            string  `json:"placa"`
	PlacasOrden      int     `json:"placas_orden"`
	PlacasUsadas     int     `json:"placas_usadas"`
	Diferencia       int     `json:"diferencia"`
	PorcentajeExceso float64 `json:"porcentaje_exceso"`
	Mensaje          string  `json:"mensaje"`
}

type alertRequest struct {
	Items []alertItem `json:"items"`
}

func SetupAlertsRoutes(r *gin.Engine) {
	db, err := sql.Open("sqlserver", alertConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database for alerts: %v", err)
	}

	if err := ensureAlertasTable(db); err != nil {
		log.Printf("Warning: could not ensure ALERTAS table exists: %v", err)
	}

	r.POST("/app/alertas", func(c *gin.Context) {
		var payload alertRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to bind JSON for alerts: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
			return
		}

		if len(payload.Items) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No alerts to insert", "inserted": 0})
			return
		}

		if err := ensureAlertasTable(db); err != nil {
			log.Printf("Failed to ensure ALERTAS table exists: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ensure alerts table"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin alerts transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
			return
		}

		defer func() {
			if p := recover(); p != nil {
				log.Printf("Panic during alerts transaction: %v", p)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction rolled back due to panic"})
			}
		}()

		for _, item := range payload.Items {
			if item.Placa == "" || item.Mensaje == "" {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Each alert must include placa and mensaje"})
				return
			}

			_, err = tx.Exec(
				`INSERT INTO dbo.ALERTAS
					(ID_PROCESO, PROCESO, NVNUMERO, CLIENTE, PRODUCTO, PLACA, PLACAS_ORDEN, PLACAS_USADAS, DIFERENCIA, PORCENTAJE_EXCESO, MENSAJE, FECHA)
				 VALUES
					(@ID_PROCESO, @PROCESO, @NVNUMERO, @CLIENTE, @PRODUCTO, @PLACA, @PLACAS_ORDEN, @PLACAS_USADAS, @DIFERENCIA, @PORCENTAJE_EXCESO, @MENSAJE, GETDATE())`,
				sql.Named("ID_PROCESO", item.IDProceso),
				sql.Named("PROCESO", item.Proceso),
				sql.Named("NVNUMERO", item.NVNumero),
				sql.Named("CLIENTE", item.Cliente),
				sql.Named("PRODUCTO", item.Producto),
				sql.Named("PLACA", item.Placa),
				sql.Named("PLACAS_ORDEN", item.PlacasOrden),
				sql.Named("PLACAS_USADAS", item.PlacasUsadas),
				sql.Named("DIFERENCIA", item.Diferencia),
				sql.Named("PORCENTAJE_EXCESO", item.PorcentajeExceso),
				sql.Named("MENSAJE", item.Mensaje),
			)
			if err != nil {
				log.Printf("Failed to insert alert: %v", err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert alert"})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit alerts transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Alerts inserted", "inserted": len(payload.Items)})
	})

	r.GET("/app/alertas", func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT ID, ID_PROCESO, PROCESO, NVNUMERO, CLIENTE, PRODUCTO, PLACA, PLACAS_ORDEN, PLACAS_USADAS, DIFERENCIA, PORCENTAJE_EXCESO, MENSAJE, FECHA
			FROM dbo.ALERTAS
			ORDER BY FECHA DESC, ID DESC`)
		if err != nil {
			log.Printf("Failed to query alerts: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve alerts"})
			return
		}
		defer rows.Close()

		var results []gin.H
		for rows.Next() {
			var (
				id               int
				idProceso        sql.NullInt64
				proceso          sql.NullString
				nvNumero         sql.NullString
				cliente          sql.NullString
				producto         sql.NullString
				placa            sql.NullString
				placasOrden      sql.NullInt64
				placasUsadas     sql.NullInt64
				diferencia       sql.NullInt64
				porcentajeExceso sql.NullFloat64
				mensaje          sql.NullString
				fecha            sql.NullTime
			)

			if err := rows.Scan(&id, &idProceso, &proceso, &nvNumero, &cliente, &producto, &placa, &placasOrden, &placasUsadas, &diferencia, &porcentajeExceso, &mensaje, &fecha); err != nil {
				log.Printf("Failed to scan alert row: %v", err)
				continue
			}

			results = append(results, gin.H{
				"id":                id,
				"id_proceso":        nullInt64ToInt(idProceso),
				"proceso":           nullStringToString(proceso),
				"nv_numero":         nullStringToString(nvNumero),
				"cliente":           nullStringToString(cliente),
				"producto":          nullStringToString(producto),
				"placa":             nullStringToString(placa),
				"placas_orden":      nullInt64ToInt(placasOrden),
				"placas_usadas":     nullInt64ToInt(placasUsadas),
				"diferencia":        nullInt64ToInt(diferencia),
				"porcentaje_exceso": nullFloat64ToFloat64(porcentajeExceso),
				"mensaje":           nullStringToString(mensaje),
				"fecha":             nullTimeToString(fecha),
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": results})
	})
}

func alertConnectionString() string {
	return "Server=" + os.Getenv("SQL_SERVER") + "\\" + os.Getenv("SQL_INSTANCE") + ";Database=" + os.Getenv("SQL_DATABASE2") + ";User Id=" + os.Getenv("SQL_USER") + ";Password=" + os.Getenv("SQL_PASSWORD") + ";Encrypt=disable"
}

func ensureAlertasTable(db *sql.DB) error {
	_, err := db.Exec(`
		IF NOT EXISTS (
			SELECT * FROM INFORMATION_SCHEMA.TABLES
			WHERE TABLE_SCHEMA = 'dbo' AND TABLE_NAME = 'ALERTAS'
		)
		BEGIN
			CREATE TABLE dbo.ALERTAS (
				ID INT IDENTITY(1,1) PRIMARY KEY,
				ID_PROCESO INT NULL,
				PROCESO NVARCHAR(MAX) NULL,
				NVNUMERO NVARCHAR(MAX) NULL,
				CLIENTE NVARCHAR(MAX) NULL,
				PRODUCTO NVARCHAR(MAX) NULL,
				PLACA NVARCHAR(MAX) NOT NULL,
				PLACAS_ORDEN INT NULL,
				PLACAS_USADAS INT NULL,
				DIFERENCIA INT NULL,
				PORCENTAJE_EXCESO DECIMAL(10,2) NULL,
				MENSAJE NVARCHAR(MAX) NOT NULL,
				FECHA DATETIME DEFAULT GETDATE()
			);
		END
	`)
	return err
}

func nullStringToString(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

func nullInt64ToInt(v sql.NullInt64) int {
	if v.Valid {
		return int(v.Int64)
	}
	return 0
}

func nullFloat64ToFloat64(v sql.NullFloat64) float64 {
	if v.Valid {
		return v.Float64
	}
	return 0
}

func nullTimeToString(v sql.NullTime) string {
	if v.Valid {
		return v.Time.Format("2006-01-02 15:04:05")
	}
	return ""
}
