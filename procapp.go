package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func SetupRoutes(r *gin.Engine) {
	// Use environment variables for database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r.GET("/troquelado", queryHandler(db, `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN TROQUELADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`))

	r.GET("/troquelado2", queryHandler(db, `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN TROQUELADO2 p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`))

	r.GET("/encolado", queryHandler(db, `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN ENCOLADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`))

	r.GET("/encolado2", queryHandler(db, `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN ENCOLADO2 p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`))

	r.GET("/multiple", queryHandler(db, `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN MULTIPLE p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`))

	r.GET("/multiple2", queryHandler(db, `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN MULTIPLE2 p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`))
}

func SetupPostRoutes(r *gin.Engine) {
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r.POST("/update-troquelado", func(c *gin.Context) {
		var items []struct {
			ID                int           `json:"ID"`
			CANT_A_FABRICAR   int           `json:"CANT_A_FABRICAR"`
			TransformedPlacas []interface{} `json:"transformedPlacas"`
			PlacasUsadas      []interface{} `json:"placasUsadas"`
		}
		if err := c.ShouldBindJSON(&items); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to begin transaction"})
			return
		}

		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Transaction rolled back due to panic"})
			}
		}()

		_, err = tx.Exec("DELETE FROM TROQUELADO")
		if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from TROQUELADO"})
			return
		}

		priority := 1
		for _, item := range items {
			if item.ID == 0 || item.CANT_A_FABRICAR == 0 || len(item.TransformedPlacas) == 0 || len(item.PlacasUsadas) == 0 {
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid data structure"})
				return
			}

			_, err = tx.Exec(
				"INSERT INTO TROQUELADO (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) VALUES (?, ?, ?, ?, ?)",
				item.ID, item.CANT_A_FABRICAR, priority, toJSON(item.TransformedPlacas), toJSON(item.PlacasUsadas),
			)
			if err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into TROQUELADO"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into TROQUELADO"})
	})
	r.POST("/update-troquelado2", func(c *gin.Context) {
		var items []struct {
			ID                int           `json:"ID"`
			CANT_A_FABRICAR   int           `json:"CANT_A_FABRICAR"`
			TransformedPlacas []interface{} `json:"transformedPlacas"`
			PlacasUsadas      []interface{} `json:"placasUsadas"`
		}
		if err := c.ShouldBindJSON(&items); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to begin transaction"})
			return
		}

		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Transaction rolled back due to panic"})
			}
		}()

		_, err = tx.Exec("DELETE FROM TROQUELADO2")
		if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from TROQUELADO2"})
			return
		}

		priority := 1
		for _, item := range items {
			if item.ID == 0 || item.CANT_A_FABRICAR == 0 || len(item.TransformedPlacas) == 0 || len(item.PlacasUsadas) == 0 {
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid data structure"})
				return
			}

			_, err = tx.Exec(
				"INSERT INTO TROQUELADO2 (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) VALUES (?, ?, ?, ?, ?)",
				item.ID, item.CANT_A_FABRICAR, priority, toJSON(item.TransformedPlacas), toJSON(item.PlacasUsadas),
			)
			if err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into TROQUELADO2"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}
		c.JSON(201, gin.H{"message": "Inserted into TROQUELADO2"})
	})
}

func queryHandler(db *sql.DB, query string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type Response struct {
			ID              int    `json:"ID"`
			NVNUMERO        int    `json:"NVNUMERO"`
			NOMAUX          string `json:"NOMAUX"`
			FECHA_ENTREGA   string `json:"FECHA_ENTREGA"`
			PROCESO         string `json:"PROCESO"`
			DETPROD         string `json:"DETPROD"`
			CANTPROD        int    `json:"CANTPROD"`
			CANT_A_FABRICAR int    `json:"CANT_A_FABRICAR"`
			PLACAS_A_USAR   string `json:"PLACAS_A_USAR"`
			CANTIDAD_PLACAS string `json:"CANTIDAD_PLACAS"`
		}

		results := []Response{}
		for rows.Next() {
			var res Response
			err := rows.Scan(
				&res.ID, &res.NVNUMERO, &res.NOMAUX, &res.FECHA_ENTREGA, &res.PROCESO,
				&res.DETPROD, &res.CANTPROD, &res.CANT_A_FABRICAR, &res.PLACAS_A_USAR, &res.CANTIDAD_PLACAS,
			)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			results = append(results, res)
		}
		c.JSON(200, results)
	}
}

func toJSON(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return ""
	}
	return string(jsonData)
}
