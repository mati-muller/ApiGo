package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func Reportes(r *gin.Engine) {
	r.POST("/reportes/update", updateHandler)
	r.GET(("/reportes/historial"), getHistorialHandler)
}

func updateHandler(c *gin.Context) {
	var reqBody struct {
		ID              int      `json:"id"`
		SubtractValue   int      `json:"subtractValue"`
		Placas          []string `json:"placas"`
		PlacasUsadas    []int    `json:"placasUsadas"`
		PlacasBuenas    []int    `json:"placasBuenas"`
		PlacasMalas     []int    `json:"placasMalas"`
		TiempoTotal     float64  `json:"tiempoTotal"`
		User            string   `json:"user"`
		StockCant       int      `json:"stockCant"`
		NumeroPersonas  int      `json:"numeroPersonas"`
		AddToStock      bool     `json:"addToStock"`
		RemoveFromStock bool     `json:"removeFromStock"`
		RemoveStockCant int      `json:"removeStockCant"`
	}

	// Parse and validate request body
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Println("Error parsing request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Validate required fields
	if reqBody.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field 'id' is required and must be greater than 0"})
		return
	}
	if reqBody.SubtractValue < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field 'subtractValue' must be non-negative"})
		return
	}
	// Eliminada la validación de placas vacías
	if len(reqBody.PlacasUsadas) != len(reqBody.Placas) || len(reqBody.PlacasBuenas) != len(reqBody.Placas) || len(reqBody.PlacasMalas) != len(reqBody.Placas) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mismatch in the number of Placas, PlacasUsadas, PlacasBuenas, and PlacasMalas"})
		return
	}

	// Use environment variables for database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		log.Println("Database connection error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println("Transaction begin error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction begin error"})
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Println("Transaction rollback due to panic:", p)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		}
	}()

	// Fetch current data
	var currentPlacas, currentPlacasUsadas, currentPlacasBuenas, currentPlacasMalas, currentUser sql.NullString
	err = tx.QueryRow(`
		SELECT PLACA, PLACAS_USADAS, PLACAS_BUENAS, PLACAS_MALAS, [USER]
		FROM procesos2
		WHERE ID = @p1
	`, reqBody.ID).Scan(&currentPlacas, &currentPlacasUsadas, &currentPlacasBuenas, &currentPlacasMalas, &currentUser)
	if err == sql.ErrNoRows {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	} else if err != nil {
		tx.Rollback()
		log.Println("Query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error"})
		return
	}

	// Use empty array if NULL
	placasStr := "[]"
	if currentPlacas.Valid {
		placasStr = currentPlacas.String
	}
	placasUsadasStr := "[]"
	if currentPlacasUsadas.Valid {
		placasUsadasStr = currentPlacasUsadas.String
	}
	placasBuenasStr := "[]"
	if currentPlacasBuenas.Valid {
		placasBuenasStr = currentPlacasBuenas.String
	}
	placasMalasStr := "[]"
	if currentPlacasMalas.Valid {
		placasMalasStr = currentPlacasMalas.String
	}
	userStr := ""
	if currentUser.Valid {
		userStr = currentUser.String
	}

	// Merge user if different
	if userStr != reqBody.User && userStr != "" {
		userStr = userStr + ", " + reqBody.User
	} else if userStr == "" {
		userStr = reqBody.User
	}

	// Parse JSON fields
	var placas []string
	if err := json.Unmarshal([]byte(placasStr), &placas); err != nil {
		log.Println("Error parsing PLACA JSON:", err)
	}
	// Use integer slices for counts
	var placasUsadasArr, placasBuenasArr, placasMalasArr []int
	if err := json.Unmarshal([]byte(placasUsadasStr), &placasUsadasArr); err != nil {
		log.Println("Error parsing PLACAS_USADAS JSON:", err)
		placasUsadasArr = make([]int, len(placas))
	}
	if err := json.Unmarshal([]byte(placasBuenasStr), &placasBuenasArr); err != nil {
		log.Println("Error parsing PLACAS_BUENAS JSON:", err)
		placasBuenasArr = make([]int, len(placas))
	}
	if err := json.Unmarshal([]byte(placasMalasStr), &placasMalasArr); err != nil {
		log.Println("Error parsing PLACAS_MALAS JSON:", err)
		placasMalasArr = make([]int, len(placas))
	}

	// Merge new data
	for i, placa := range reqBody.Placas {
		index := -1
		for j, existingPlaca := range placas {
			if existingPlaca == placa {
				index = j
				break
			}
		}
		if index != -1 {
			// sum counts for existing placa
			placasUsadasArr[index] += reqBody.PlacasUsadas[i]
			placasBuenasArr[index] += reqBody.PlacasBuenas[i]
			placasMalasArr[index] += reqBody.PlacasMalas[i]
		} else {
			// append new placa and counts
			placas = append(placas, placa)
			placasUsadasArr = append(placasUsadasArr, reqBody.PlacasUsadas[i])
			placasBuenasArr = append(placasBuenasArr, reqBody.PlacasBuenas[i])
			placasMalasArr = append(placasMalasArr, reqBody.PlacasMalas[i])
		}
	}

	// Update database with JSON arrays of strings and ints
	placasBytes, _ := json.Marshal(placas)
	placasStr = string(placasBytes)
	placasUsadasBytes, _ := json.Marshal(placasUsadasArr)
	placasUsadasStr = string(placasUsadasBytes)
	placasBuenasBytes, _ := json.Marshal(placasBuenasArr)
	placasBuenasStr = string(placasBuenasBytes)
	placasMalasBytes, _ := json.Marshal(placasMalasArr)
	placasMalasStr = string(placasMalasBytes)

	_, err = tx.Exec(`
		UPDATE procesos2
		SET CANT_A_PROD = CASE WHEN CANT_A_PROD - @p1 <= 0 THEN 0 ELSE CANT_A_PROD - @p1 END,
			CANT_PROD = CANT_PROD + @p1,
			ESTADO_PROC = CASE WHEN CANT_A_PROD - @p1 <= 0 THEN 'LISTO' ELSE ESTADO_PROC END,
			PLACA = @p2,
			PLACAS_USADAS = @p3,
			PLACAS_BUENAS = @p4,
			PLACAS_MALAS = @p5,
			TIEMPO_TOTAL = TIEMPO_TOTAL + @p6,
			[USER] = @p7,
			STOCK_CANT = STOCK_CANT + @p8,
			NUMERO_PERSONAS = @p9
		WHERE ID = @p10
	`, reqBody.SubtractValue, placasStr, placasUsadasStr, placasBuenasStr, placasMalasStr,
		reqBody.TiempoTotal, userStr, reqBody.StockCant, reqBody.NumeroPersonas, reqBody.ID)
	if err != nil {
		tx.Rollback()
		log.Println("Update error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update error"})
		return
	}

	// Insert data into HISTORIAL table
	placasBytes, _ = json.Marshal(reqBody.Placas)
	placasStr = string(placasBytes)
	placasUsadasBytes, _ = json.Marshal(reqBody.PlacasUsadas)
	placasUsadasStr = string(placasUsadasBytes)
	placasBuenasBytes, _ = json.Marshal(reqBody.PlacasBuenas)
	placasBuenasStr = string(placasBuenasBytes)
	placasMalasBytes, _ = json.Marshal(reqBody.PlacasMalas)
	placasMalasStr = string(placasMalasBytes)

	_, err = tx.Exec(`
		INSERT INTO HISTORIAL (
			ID_PROCESO, CANTIDAD, PLACA, PLACAS_USADAS, PLACAS_BUENAS, PLACAS_MALAS, 
			TIEMPO_TOTAL, NUMERO_PERSONAS, STOCK, [USER], STOCK_CANT
		) VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11
		)
	`, reqBody.ID, reqBody.SubtractValue, placasStr, placasUsadasStr,
		placasBuenasStr, placasMalasStr, reqBody.TiempoTotal,
		reqBody.NumeroPersonas, "", reqBody.User, reqBody.StockCant)
	if err != nil {
		tx.Rollback()
		log.Println("Insert into HISTORIAL error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert into HISTORIAL error"})
		return
	}

	// Handle Add to Stock
	if reqBody.AddToStock {
		if reqBody.StockCant <= 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock quantity to add."})
			return
		}

		if err := updateStock(tx, reqBody.ID, reqBody.StockCant, "Add"); err != nil {
			tx.Rollback()
			log.Println("Add to stock error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Add to stock error"})
			return
		}
	}

	// Handle Remove from Stock
	if reqBody.RemoveFromStock {
		if reqBody.RemoveStockCant <= 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock quantity to remove."})
			return
		}

		if err := updateStock(tx, reqBody.ID, reqBody.RemoveStockCant, "Remove"); err != nil {
			tx.Rollback()
			log.Println("Remove from stock error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Remove from stock error"})
			return
		}
	}

	// Deduct inventory for each placa
	for i, placa := range reqBody.Placas {
		if err := deductInventory(tx, placa, reqBody.PlacasUsadas[i]); err != nil {
			tx.Rollback()
			log.Println("Inventory deduction error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Println("Transaction commit error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
}

func mergeStringValues(existing, new string) string {
	existingValue, _ := strconv.Atoi(existing)
	newValue, _ := strconv.Atoi(new)
	return strconv.Itoa(existingValue + newValue)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// updateStock updates the stock quantity in procesos2 and desprod in procesos.
func updateStock(tx *sql.Tx, processID, quantity int, action string) error {
	// Fetch DETPROD from procesos using the provided ID
	var detProd string
	err := tx.QueryRow(`
		SELECT DETPROD
		FROM procesos
		WHERE ID = @p1
	`, processID).Scan(&detProd)
	if err != nil {
		return fmt.Errorf("failed to fetch DETPROD for ID %d: %v", processID, err)
	}

	if action == "Add" {
		// Update STOCK_CANT in procesos2
		_, err := tx.Exec(`
			UPDATE procesos2
			SET STOCK = 'Add',
				STOCK_CANT = STOCK_CANT + @p1
			WHERE ID = @p2
		`, quantity, processID)
		if err != nil {
			return fmt.Errorf("failed to add stock in procesos2: %v", err)
		}

		// Check if DETPROD exists in inventory
		var currentCantidad int
		err = tx.QueryRow(`
			SELECT cantidad
			FROM inventario
			WHERE placa = @p1
		`, detProd).Scan(&currentCantidad)
		if err == sql.ErrNoRows {
			// Insert DETPROD into inventory letting SQL Server auto-generate the ID
			_, err = tx.Exec(`
				INSERT INTO inventario (placa, cantidad)
				VALUES (@p1, @p2)
			`, detProd, quantity)
			if err != nil {
				return fmt.Errorf("failed to insert DETPROD into inventory: %v", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to fetch inventory for DETPROD: %v", err)
		} else {
			// Update existing inventory record for DETPROD
			_, err = tx.Exec(`
				UPDATE inventario
				SET cantidad = cantidad + @p1
				WHERE placa = @p2
			`, quantity, detProd)
			if err != nil {
				return fmt.Errorf("failed to update inventory for DETPROD: %v", err)
			}
		}
	} else if action == "Remove" {
		// Update STOCK_CANT in procesos2
		_, err := tx.Exec(`
			UPDATE procesos2
			SET STOCK = 'Remove',
				STOCK_CANT = CASE WHEN STOCK_CANT - @p1 < 0 THEN 0 ELSE STOCK_CANT - @p1 END
			WHERE ID = @p2
		`, quantity, processID)
		if err != nil {
			return fmt.Errorf("failed to remove stock in procesos2: %v", err)
		}

		// Deduct DETPROD from inventory
		_, err = tx.Exec(`
			UPDATE inventario
			SET cantidad = CASE WHEN cantidad - @p1 < 0 THEN 0 ELSE cantidad - @p1 END
			WHERE placa = @p2
		`, quantity, detProd)
		if err != nil {
			return fmt.Errorf("failed to deduct DETPROD from inventory: %v", err)
		}
	} else {
		return fmt.Errorf("invalid stock action: %s", action)
	}
	return nil
}

// deductInventory deducts the specified quantity from the inventory for a given placa.
func deductInventory(tx *sql.Tx, placa string, quantityToDeduct int) error {
	rows, err := tx.Query(`
		SELECT ID, cantidad
		FROM inventario
		WHERE placa = @p1
		ORDER BY ID ASC
	`, placa)
	if err != nil {
		return fmt.Errorf("inventory query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var inventoryID, currentCantidad int
		if err := rows.Scan(&inventoryID, &currentCantidad); err != nil {
			return fmt.Errorf("row scan error: %v", err)
		}

		if quantityToDeduct <= 0 {
			break
		}

		subtractFromCurrent := min(currentCantidad, quantityToDeduct)
		_, err := tx.Exec(`
			UPDATE inventario
			SET cantidad_total_usada = cantidad_total_usada + @p1,
				cantidad = CASE WHEN cantidad - @p1 <= 0 THEN 0 ELSE cantidad - @p1 END,
				precio_total = (CASE WHEN cantidad - @p1 <= 0 THEN 0 ELSE cantidad - @p1 END) * precio_pp
			WHERE ID = @p2
		`, subtractFromCurrent, inventoryID)
		if err != nil {
			return fmt.Errorf("inventory update error: %v", err)
		}

		quantityToDeduct -= subtractFromCurrent
	}

	if quantityToDeduct > 0 {
		return fmt.Errorf("not enough inventory for placa: %s", placa)
	}

	return nil
}

func getHistorialHandler(c *gin.Context) {
	// Query the HISTORIAL table with a JOIN on procesos
	rows, err := db.Query(`
        SELECT h.ID, h.ID_PROCESO, h.CANTIDAD, h.PLACA, h.PLACAS_USADAS, h.PLACAS_BUENAS, h.PLACAS_MALAS, 
               h.TIEMPO_TOTAL, h.NUMERO_PERSONAS, h.STOCK, h.[USER], h.STOCK_CANT,
               p.NVNUMERO, p.FECHA_ENTREGA, p.NOMAUX, p.NVCANT, p.DETPROD, p.PROCESO, h.FECHA
        FROM HISTORIAL h
        JOIN procesos p ON h.ID_PROCESO = p.ID
    `)
	if err != nil {
		log.Println("Error querying HISTORIAL table with JOIN:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying HISTORIAL table with JOIN"})
		return
	}
	defer rows.Close()

	// Parse rows into a slice of maps
	var historial []map[string]interface{}
	for rows.Next() {
		var (
			id, idProceso, cantidad, numeroPersonas, stockCant, nvnumero, nvcant                                 int
			placa, placasUsadas, placasBuenas, placasMalas, stock, user, fechaEntrega, codprod, detprod, proceso string
			tiempoTotal                                                                                          float64
			fecha                                                                                                sql.NullString
		)
		if err := rows.Scan(&id, &idProceso, &cantidad, &placa, &placasUsadas, &placasBuenas, &placasMalas, &tiempoTotal, &numeroPersonas, &stock, &user, &stockCant, &nvnumero, &fechaEntrega, &codprod, &nvcant, &detprod, &proceso, &fecha); err != nil {
			log.Println("Error scanning row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row"})
			return
		}

		historial = append(historial, map[string]interface{}{
			"ID":              id,
			"ID_PROCESO":      idProceso,
			"CANTIDAD":        cantidad,
			"PLACA":           placa,
			"PLACAS_USADAS":   placasUsadas,
			"PLACAS_BUENAS":   placasBuenas,
			"PLACAS_MALAS":    placasMalas,
			"TIEMPO_TOTAL":    tiempoTotal,
			"NUMERO_PERSONAS": numeroPersonas,
			"STOCK":           stock,
			"USER":            user,
			"STOCK_CANT":      stockCant,
			"NVNUMERO":        nvnumero,
			"FECHA_ENTREGA":   fechaEntrega,
			"CODPROD":         codprod,
			"NVCANT":          nvcant,
			"DETPROD":         detprod,
			"PROCESO":         proceso,
			"FECHA":           nil,
		})
		if fecha.Valid {
			historial[len(historial)-1]["FECHA"] = fecha.String
		}
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating rows"})
		return
	}

	c.JSON(http.StatusOK, historial)
}
