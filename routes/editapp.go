package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupEditsApp(r *gin.Engine) {
	r.POST("/edit-app-troquelado", editTroqueladoapp)
	r.POST("/edit-app-emplacado", editEmplacadoapp)
	r.POST("/edit-app-trozado", editTrozadoapp)
	r.POST("/edit-app-encolado", editEncoladoapp)
}

func editTroqueladoapp(c *gin.Context) {
	type request struct {
		ID           int      `json:"id"`
		Cantidad     int      `json:"cantidad"`
		Placas       []string `json:"placas"`
		PlacasBuenas []int    `json:"placasBuenas"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Log de los datos recibidos
	fmt.Printf("Datos recibidos en %s: ID=%d, Cantidad=%d, Placas=%v, PlacasBuenas=%v\n", c.FullPath(), req.ID, req.Cantidad, req.Placas, req.PlacasBuenas)

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	query := `UPDATE TROQUELADO SET CANT_A_FABRICAR = CASE WHEN CANT_A_FABRICAR - @cantidad < 0 THEN 0 ELSE CANT_A_FABRICAR - @cantidad END, CANTIDAD_PRODUCIDA = ISNULL(CANTIDAD_PRODUCIDA,0) + @cantidad WHERE ID = @id`
	_, err = db.Exec(query, sql.Named("cantidad", req.Cantidad), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update process"})
		return
	}

	// Leer el valor actual de CANTIDAD_PLACAS para el ID recibido
	var cantidadPlacasJSON string
	err = db.QueryRow("SELECT CANTIDAD_PLACAS FROM TROQUELADO WHERE ID = @id", sql.Named("id", req.ID)).Scan(&cantidadPlacasJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo leer CANTIDAD_PLACAS"})
		return
	}

	// Parsear el JSON a un array de ints
	var cantidadPlacasArr []int
	err = json.Unmarshal([]byte(cantidadPlacasJSON), &cantidadPlacasArr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CANTIDAD_PLACAS no es un array JSON válido"})
		return
	}

	// Restar placas buenas según índice
	for i := range req.Placas {
		if i >= len(req.PlacasBuenas) || i >= len(cantidadPlacasArr) {
			continue
		}
		cantidadPlacasArr[i] -= req.PlacasBuenas[i]
		if cantidadPlacasArr[i] < 0 {
			cantidadPlacasArr[i] = 0
		}
	}

	// Guardar el array actualizado como JSON
	nuevaCantidadPlacasJSON, err := json.Marshal(cantidadPlacasArr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo serializar el array actualizado"})
		return
	}

	_, err = db.Exec("UPDATE TROQUELADO SET CANTIDAD_PLACAS = @nuevo WHERE ID = @id", sql.Named("nuevo", string(nuevaCantidadPlacasJSON)), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar CANTIDAD_PLACAS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process updated"})
}

func editEmplacadoapp(c *gin.Context) {
	type request struct {
		ID           int      `json:"id"`
		Cantidad     int      `json:"cantidad"`
		Placas       []string `json:"placas"`
		PlacasBuenas []int    `json:"placasBuenas"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Log de los datos recibidos
	fmt.Printf("Datos recibidos en %s: ID=%d, Cantidad=%d, Placas=%v, PlacasBuenas=%v\n", c.FullPath(), req.ID, req.Cantidad, req.Placas, req.PlacasBuenas)

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	query := `UPDATE EMPLACADO SET CANT_A_FABRICAR = CASE WHEN CANT_A_FABRICAR - @cantidad < 0 THEN 0 ELSE CANT_A_FABRICAR - @cantidad END, CANTIDAD_PRODUCIDA = ISNULL(CANTIDAD_PRODUCIDA,0) + @cantidad WHERE ID = @id`
	_, err = db.Exec(query, sql.Named("cantidad", req.Cantidad), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update process"})
		return
	}

	// Leer el valor actual de CANTIDAD_PLACAS para el ID recibido
	var cantidadPlacasJSON string
	err = db.QueryRow("SELECT CANTIDAD_PLACAS FROM EMPLACADO WHERE ID = @id", sql.Named("id", req.ID)).Scan(&cantidadPlacasJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo leer CANTIDAD_PLACAS"})
		return
	}

	// Parsear el JSON a un array de ints
	var cantidadPlacasArr []int
	err = json.Unmarshal([]byte(cantidadPlacasJSON), &cantidadPlacasArr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CANTIDAD_PLACAS no es un array JSON válido"})
		return
	}

	// Restar placas buenas según índice
	for i := range req.Placas {
		if i >= len(req.PlacasBuenas) || i >= len(cantidadPlacasArr) {
			continue
		}
		cantidadPlacasArr[i] -= req.PlacasBuenas[i]
		if cantidadPlacasArr[i] < 0 {
			cantidadPlacasArr[i] = 0
		}
	}

	// Guardar el array actualizado como JSON
	nuevaCantidadPlacasJSON, err := json.Marshal(cantidadPlacasArr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo serializar el array actualizado"})
		return
	}

	_, err = db.Exec("UPDATE EMPLACADO SET CANTIDAD_PLACAS = @nuevo WHERE ID = @id", sql.Named("nuevo", string(nuevaCantidadPlacasJSON)), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar CANTIDAD_PLACAS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process updated"})
}

func editTrozadoapp(c *gin.Context) {
	type request struct {
		ID           int      `json:"id"`
		Cantidad     int      `json:"cantidad"`
		Placas       []string `json:"placas"`
		PlacasBuenas []int    `json:"placasBuenas"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Log de los datos recibidos
	fmt.Printf("Datos recibidos en %s: ID=%d, Cantidad=%d, Placas=%v, PlacasBuenas=%v\n", c.FullPath(), req.ID, req.Cantidad, req.Placas, req.PlacasBuenas)

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	query := `UPDATE TROZADO SET CANT_A_FABRICAR = CASE WHEN CANT_A_FABRICAR - @cantidad < 0 THEN 0 ELSE CANT_A_FABRICAR - @cantidad END, CANTIDAD_PRODUCIDA = ISNULL(CANTIDAD_PRODUCIDA,0) + @cantidad WHERE ID = @id`
	_, err = db.Exec(query, sql.Named("cantidad", req.Cantidad), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update process"})
		return
	}

	// Leer el valor actual de CANTIDAD_PLACAS para el ID recibido
	var cantidadPlacasJSON string
	err = db.QueryRow("SELECT CANTIDAD_PLACAS FROM TROZADO WHERE ID = @id", sql.Named("id", req.ID)).Scan(&cantidadPlacasJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo leer CANTIDAD_PLACAS"})
		return
	}

	// Parsear el JSON a un array de ints
	var cantidadPlacasArr []int
	err = json.Unmarshal([]byte(cantidadPlacasJSON), &cantidadPlacasArr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CANTIDAD_PLACAS no es un array JSON válido"})
		return
	}

	// Restar placas buenas según índice
	for i := range req.Placas {
		if i >= len(req.PlacasBuenas) || i >= len(cantidadPlacasArr) {
			continue
		}
		cantidadPlacasArr[i] -= req.PlacasBuenas[i]
		if cantidadPlacasArr[i] < 0 {
			cantidadPlacasArr[i] = 0
		}
	}

	// Guardar el array actualizado como JSON
	nuevaCantidadPlacasJSON, err := json.Marshal(cantidadPlacasArr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo serializar el array actualizado"})
		return
	}

	_, err = db.Exec("UPDATE TROZADO SET CANTIDAD_PLACAS = @nuevo WHERE ID = @id", sql.Named("nuevo", string(nuevaCantidadPlacasJSON)), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar CANTIDAD_PLACAS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process updated"})
}

func editEncoladoapp(c *gin.Context) {
	type request struct {
		ID           int      `json:"id"`
		Cantidad     int      `json:"cantidad"`
		Placas       []string `json:"placas"`
		PlacasBuenas []int    `json:"placasBuenas"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Log de los datos recibidos
	fmt.Printf("Datos recibidos en %s: ID=%d, Cantidad=%d, Placas=%v, PlacasBuenas=%v\n", c.FullPath(), req.ID, req.Cantidad, req.Placas, req.PlacasBuenas)

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	query := `UPDATE ENCOLADO SET CANT_A_FABRICAR = CASE WHEN CANT_A_FABRICAR - @cantidad < 0 THEN 0 ELSE CANT_A_FABRICAR - @cantidad END, CANTIDAD_PRODUCIDA = ISNULL(CANTIDAD_PRODUCIDA,0) + @cantidad WHERE ID = @id`
	_, err = db.Exec(query, sql.Named("cantidad", req.Cantidad), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update process"})
		return
	}

	// Leer el valor actual de CANTIDAD_PLACAS para el ID recibido
	var cantidadPlacasJSON string
	err = db.QueryRow("SELECT CANTIDAD_PLACAS FROM ENCOLADO WHERE ID = @id", sql.Named("id", req.ID)).Scan(&cantidadPlacasJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo leer CANTIDAD_PLACAS"})
		return
	}

	// Parsear el JSON a un array de ints
	var cantidadPlacasArr []int
	err = json.Unmarshal([]byte(cantidadPlacasJSON), &cantidadPlacasArr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CANTIDAD_PLACAS no es un array JSON válido"})
		return
	}

	// Restar placas buenas según índice
	for i := range req.Placas {
		if i >= len(req.PlacasBuenas) || i >= len(cantidadPlacasArr) {
			continue
		}
		cantidadPlacasArr[i] -= req.PlacasBuenas[i]
		if cantidadPlacasArr[i] < 0 {
			cantidadPlacasArr[i] = 0
		}
	}

	// Guardar el array actualizado como JSON
	nuevaCantidadPlacasJSON, err := json.Marshal(cantidadPlacasArr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo serializar el array actualizado"})
		return
	}

	_, err = db.Exec("UPDATE ENCOLADO SET CANTIDAD_PLACAS = @nuevo WHERE ID = @id", sql.Named("nuevo", string(nuevaCantidadPlacasJSON)), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar CANTIDAD_PLACAS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process updated"})
}
