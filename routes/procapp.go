package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

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
}

func queryDatabase(c *gin.Context, query string) {
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	results := []map[string]interface{}{}

	for rows.Next() {
		row := make([]interface{}, len(columns))
		rowPointers := make([]interface{}, len(columns))
		for i := range row {
			rowPointers[i] = &row[i]
		}

		if err := rows.Scan(rowPointers...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := map[string]interface{}{}
		for i, col := range columns {
			result[col] = row[i]
		}
		results = append(results, result)
	}

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

func SetupPostRoutes(r *gin.Engine) {
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r.POST("/app/update-troquelado", func(c *gin.Context) {
		var payload struct {
			Items []struct {
				ID                int      `json:"ID"`
				CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
				TransformedPlacas []string `json:"transformedPlacas"`
				PlacasUsadas      []int    `json:"placasUsadas"`
			} `json:"items"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(400, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to begin transaction"})
			return
		}

		defer func() {
			if p := recover(); p != nil {
				log.Printf("Panic during transaction: %v", p)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Transaction rolled back due to panic"})
			}
		}()

		_, err = tx.Exec("DELETE FROM TROQUELADO")
		if err != nil {
			log.Printf("Failed to delete from TROQUELADO: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from TROQUELADO"})
			return
		}

		priority := 1
		for _, item := range payload.Items {
			if item.ID == 0 || item.CANT_A_FABRICAR == 0 {
				log.Printf("Invalid data structure: %+v", item)
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid data structure", "item": item})
				return
			}

			if len(item.TransformedPlacas) == 0 && len(item.PlacasUsadas) == 0 {
				log.Printf("Warning: Empty TransformedPlacas and PlacasUsadas for item ID %d", item.ID)
			}

			_, err = tx.Exec(
				`INSERT INTO TROQUELADO (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into TROQUELADO: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into TROQUELADO"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into TROQUELADO"})
	})

	r.POST("/app/update-troquelado2", func(c *gin.Context) {
		var payload struct {
			Items []struct {
				ID                int      `json:"ID"`
				CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
				TransformedPlacas []string `json:"transformedPlacas"`
				PlacasUsadas      []int    `json:"placasUsadas"`
			} `json:"items"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(400, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to begin transaction"})
			return
		}

		defer func() {
			if p := recover(); p != nil {
				log.Printf("Panic during transaction: %v", p)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Transaction rolled back due to panic"})
			}
		}()

		_, err = tx.Exec("DELETE FROM TROQUELADO2")
		if err != nil {
			log.Printf("Failed to delete from TROQUELADO2: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from TROQUELADO2"})
			return
		}

		priority := 1
		for _, item := range payload.Items {
			if item.ID == 0 || item.CANT_A_FABRICAR == 0 {
				log.Printf("Invalid data structure: %+v", item)
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid data structure", "item": item})
				return
			}

			if len(item.TransformedPlacas) == 0 && len(item.PlacasUsadas) == 0 {
				log.Printf("Warning: Empty TransformedPlacas and PlacasUsadas for item ID %d", item.ID)
			}

			_, err = tx.Exec(
				`INSERT INTO TROQUELADO2 (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into TROQUELADO2: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into TROQUELADO2"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into TROQUELADO2"})
	})

	r.POST("/app/update-encolado", func(c *gin.Context) {
		var payload struct {
			Items []struct {
				ID                int      `json:"ID"`
				CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
				TransformedPlacas []string `json:"transformedPlacas"`
				PlacasUsadas      []int    `json:"placasUsadas"`
			} `json:"items"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(400, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to begin transaction"})
			return
		}

		defer func() {
			if p := recover(); p != nil {
				log.Printf("Panic during transaction: %v", p)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Transaction rolled back due to panic"})
			}
		}()

		_, err = tx.Exec("DELETE FROM ENCOLADO")
		if err != nil {
			log.Printf("Failed to delete from ENCOLADO: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from ENCOLADO"})
			return
		}

		priority := 1
		for _, item := range payload.Items {
			if item.ID == 0 || item.CANT_A_FABRICAR == 0 {
				log.Printf("Invalid data structure: %+v", item)
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid data structure", "item": item})
				return
			}

			if len(item.TransformedPlacas) == 0 && len(item.PlacasUsadas) == 0 {
				log.Printf("Warning: Empty TransformedPlacas and PlacasUsadas for item ID %d", item.ID)
			}

			_, err = tx.Exec(
				`INSERT INTO ENCOLADO (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into ENCOLADO: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into ENCOLADO"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into ENCOLADO"})
	})

	r.POST("/app/update-encolado2", func(c *gin.Context) {
		var payload struct {
			Items []struct {
				ID                int      `json:"ID"`
				CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
				TransformedPlacas []string `json:"transformedPlacas"`
				PlacasUsadas      []int    `json:"placasUsadas"`
			} `json:"items"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(400, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to begin transaction"})
			return
		}

		defer func() {
			if p := recover(); p != nil {
				log.Printf("Panic during transaction: %v", p)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Transaction rolled back due to panic"})
			}
		}()

		_, err = tx.Exec("DELETE FROM ENCOLADO2")
		if err != nil {
			log.Printf("Failed to delete from ENCOLADO2: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from ENCOLADO2"})
			return
		}

		priority := 1
		for _, item := range payload.Items {
			if item.ID == 0 || item.CANT_A_FABRICAR == 0 {
				log.Printf("Invalid data structure: %+v", item)
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid data structure", "item": item})
				return
			}

			if len(item.TransformedPlacas) == 0 && len(item.PlacasUsadas) == 0 {
				log.Printf("Warning: Empty TransformedPlacas and PlacasUsadas for item ID %d", item.ID)
			}

			_, err = tx.Exec(
				`INSERT INTO ENCOLADO2 (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into ENCOLADO2: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into ENCOLADO2"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into ENCOLADO2"})
	})

	r.POST("/app/update-multiple", func(c *gin.Context) {
		var payload struct {
			Items []struct {
				ID                int      `json:"ID"`
				CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
				TransformedPlacas []string `json:"transformedPlacas"`
				PlacasUsadas      []int    `json:"placasUsadas"`
			} `json:"items"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(400, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to begin transaction"})
			return
		}

		defer func() {
			if p := recover(); p != nil {
				log.Printf("Panic during transaction: %v", p)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Transaction rolled back due to panic"})
			}
		}()

		_, err = tx.Exec("DELETE FROM MULTIPLE")
		if err != nil {
			log.Printf("Failed to delete from MULTIPLE: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from MULTIPLE"})
			return
		}

		priority := 1
		for _, item := range payload.Items {
			if item.ID == 0 || item.CANT_A_FABRICAR == 0 {
				log.Printf("Invalid data structure: %+v", item)
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid data structure", "item": item})
				return
			}

			if len(item.TransformedPlacas) == 0 && len(item.PlacasUsadas) == 0 {
				log.Printf("Warning: Empty TransformedPlacas and PlacasUsadas for item ID %d", item.ID)
			}

			_, err = tx.Exec(
				`INSERT INTO MULTIPLE (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into MULTIPLE: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into MULTIPLE"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into MULTIPLE"})
	})

	r.POST("/app/update-multiple2", func(c *gin.Context) {
		var payload struct {
			Items []struct {
				ID                int      `json:"ID"`
				CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
				TransformedPlacas []string `json:"transformedPlacas"`
				PlacasUsadas      []int    `json:"placasUsadas"`
			} `json:"items"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(400, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to begin transaction"})
			return
		}

		defer func() {
			if p := recover(); p != nil {
				log.Printf("Panic during transaction: %v", p)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Transaction rolled back due to panic"})
			}
		}()

		_, err = tx.Exec("DELETE FROM MULTIPLE2")
		if err != nil {
			log.Printf("Failed to delete from MULTIPLE2: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from MULTIPLE2"})
			return
		}

		priority := 1
		for _, item := range payload.Items {
			if item.ID == 0 || item.CANT_A_FABRICAR == 0 {
				log.Printf("Invalid data structure: %+v", item)
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid data structure", "item": item})
				return
			}

			if len(item.TransformedPlacas) == 0 && len(item.PlacasUsadas) == 0 {
				log.Printf("Warning: Empty TransformedPlacas and PlacasUsadas for item ID %d", item.ID)
			}

			_, err = tx.Exec(
				`INSERT INTO MULTIPLE2 (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into MULTIPLE2: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into MULTIPLE2"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into MULTIPLE2"})
	})
	r.POST("/app/update-pegado", func(c *gin.Context) {
		var payload struct {
			Items []struct {
				ID                int      `json:"ID"`
				CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
				TransformedPlacas []string `json:"transformedPlacas"`
				PlacasUsadas      []int    `json:"placasUsadas"`
			} `json:"items"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(400, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to begin transaction"})
			return
		}

		defer func() {
			if p := recover(); p != nil {
				log.Printf("Panic during transaction: %v", p)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Transaction rolled back due to panic"})
			}
		}()

		_, err = tx.Exec("DELETE FROM PEGADO")
		if err != nil {
			log.Printf("Failed to delete from PEGADO: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from PEGADO"})
			return
		}

		priority := 1
		for _, item := range payload.Items {
			if item.ID == 0 || item.CANT_A_FABRICAR == 0 {
				log.Printf("Invalid data structure: %+v", item)
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid data structure", "item": item})
				return
			}

			if len(item.TransformedPlacas) == 0 && len(item.PlacasUsadas) == 0 {
				log.Printf("Warning: Empty TransformedPlacas and PlacasUsadas for item ID %d", item.ID)
			}

			_, err = tx.Exec(
				`INSERT INTO PEGADO (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into PEGADO: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into PEGADO"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into PEGADO"})
	})
}

func toJSON(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return ""
	}
	return string(jsonData)
}
