package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

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
				ID                int       `json:"ID"`
				CANT_A_FABRICAR   int       `json:"CANT_A_FABRICAR"`
				TransformedPlacas []string  `json:"transformedPlacas"`
				PlacasUsadas      []float64 `json:"placasUsadas"`
			} `json:"items"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(400, gin.H{"error": "Invalid JSON payload", "details": err.Error()})
			return
		}

		for i, item := range payload.Items {
			for j, placa := range item.PlacasUsadas {
				payload.Items[i].PlacasUsadas[j] = math.Ceil(placa)
			}
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
	r.POST("/app/update-trozado", func(c *gin.Context) {
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

		_, err = tx.Exec("DELETE FROM TROZADO")
		if err != nil {
			log.Printf("Failed to delete from TROZADO: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from TROZADO"})
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
				`INSERT INTO TROZADO (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into TROZADO: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into TROZADO"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into TROZADO"})
	})
	r.POST("/app/update-calado", func(c *gin.Context) {
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

		_, err = tx.Exec("DELETE FROM CALADO")
		if err != nil {
			log.Printf("Failed to delete from CALADO: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from CALADO"})
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
				`INSERT INTO CALADO (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into CALADO: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into CALADO"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into CALADO"})
	})
	r.POST("/app/update-plizado", func(c *gin.Context) {
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

		_, err = tx.Exec("DELETE FROM PLIZADO")
		if err != nil {
			log.Printf("Failed to delete from PLIZADO: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from PLIZADO"})
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
				`INSERT INTO PLIZADO (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into PLIZADO: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into PLIZADO"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into TROZADO"})
	})
	r.POST("/app/update-impresion", func(c *gin.Context) {
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

		_, err = tx.Exec("DELETE FROM IMPRESION")
		if err != nil {
			log.Printf("Failed to delete from IMPRESION: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from IMPRESION"})
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
				`INSERT INTO IMPRESION (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into IMPRESION: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into IMPRESION"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into TROZADO"})
	})

	r.POST("/app/update-emplacado", func(c *gin.Context) {
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

		_, err = tx.Exec("DELETE FROM EMPLACADO")
		if err != nil {
			log.Printf("Failed to delete from EMPLACADO: %v", err)
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to delete from EMPLACADO"})
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
				`INSERT INTO EMPLACADO (ID, CANT_A_FABRICAR, PRIORITY, PLACAS_A_USAR, CANTIDAD_PLACAS) 
				VALUES (@ID, @CANT_A_FABRICAR, @PRIORITY, @PLACAS_A_USAR, @CANTIDAD_PLACAS)`,
				sql.Named("ID", item.ID),
				sql.Named("CANT_A_FABRICAR", item.CANT_A_FABRICAR),
				sql.Named("PRIORITY", priority),
				sql.Named("PLACAS_A_USAR", toJSON(item.TransformedPlacas)),
				sql.Named("CANTIDAD_PLACAS", toJSON(item.PlacasUsadas)),
			)
			if err != nil {
				log.Printf("Failed to insert into EMPLACADO: %v", err)
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to insert into EMPLACADO"})
				return
			}
			priority++
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			c.JSON(500, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(201, gin.H{"message": "Inserted into TROZADO"})
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
