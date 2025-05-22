package routes

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupEditsapp(r *gin.Engine) {
	r.POST("/edit-troquelado", editTroqueladoapp)
	r.POST("/edit-emplacado", editEmplacadoapp)
	r.POST("/edit-trozado", editTrozadoapp)
	r.POST("/edit-encolado", editEncoladoapp)
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

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Actualizar CANT_A_FABRICAR y CANTIDAD_PRODUCIDA como antes
	query := `UPDATE PROCESOS SET CANT_A_FABRICAR = CASE WHEN CANT_A_FABRICAR - @cantidad < 0 THEN 0 ELSE CANT_A_FABRICAR - @cantidad END, CANTIDAD_PRODUCIDA = ISNULL(CANTIDAD_PRODUCIDA,0) + @cantidad WHERE ID = @id`
	_, err = db.Exec(query, sql.Named("cantidad", req.Cantidad), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update process"})
		return
	}

	// Actualizar CANTIDAD_PLACAS por cada placa recibida
	for i, placa := range req.Placas {
		if i >= len(req.PlacasBuenas) {
			continue
		}
		cantidadBuena := req.PlacasBuenas[i]
		// Solo restar si la placa existe
		updatePlacaQuery := `UPDATE PROCESOS SET CANTIDAD_PLACAS = CASE WHEN CANTIDAD_PLACAS - @cantidadBuena < 0 THEN 0 ELSE CANTIDAD_PLACAS - @cantidadBuena END WHERE PLACA = @placa`
		res, err := db.Exec(updatePlacaQuery, sql.Named("cantidadBuena", cantidadBuena), sql.Named("placa", placa))
		if err != nil {
			continue // Si hay error, seguir con la siguiente
		}
		nRows, _ := res.RowsAffected()
		if nRows == 0 {
			continue // Si no existe la placa, no hacer nada
		}
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

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	query := `UPDATE PROCESOS SET CANT_A_FABRICAR = CASE WHEN CANT_A_FABRICAR - @cantidad < 0 THEN 0 ELSE CANT_A_FABRICAR - @cantidad END, CANTIDAD_PRODUCIDA = ISNULL(CANTIDAD_PRODUCIDA,0) + @cantidad WHERE ID = @id`
	_, err = db.Exec(query, sql.Named("cantidad", req.Cantidad), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update process"})
		return
	}

	for i, placa := range req.Placas {
		if i >= len(req.PlacasBuenas) {
			continue
		}
		cantidadBuena := req.PlacasBuenas[i]
		updatePlacaQuery := `UPDATE PROCESOS SET CANTIDAD_PLACAS = CASE WHEN CANTIDAD_PLACAS - @cantidadBuena < 0 THEN 0 ELSE CANTIDAD_PLACAS - @cantidadBuena END WHERE PLACA = @placa`
		res, err := db.Exec(updatePlacaQuery, sql.Named("cantidadBuena", cantidadBuena), sql.Named("placa", placa))
		if err != nil {
			continue
		}
		nRows, _ := res.RowsAffected()
		if nRows == 0 {
			continue
		}
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

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	query := `UPDATE PROCESOS SET CANT_A_FABRICAR = CASE WHEN CANT_A_FABRICAR - @cantidad < 0 THEN 0 ELSE CANT_A_FABRICAR - @cantidad END, CANTIDAD_PRODUCIDA = ISNULL(CANTIDAD_PRODUCIDA,0) + @cantidad WHERE ID = @id`
	_, err = db.Exec(query, sql.Named("cantidad", req.Cantidad), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update process"})
		return
	}

	for i, placa := range req.Placas {
		if i >= len(req.PlacasBuenas) {
			continue
		}
		cantidadBuena := req.PlacasBuenas[i]
		updatePlacaQuery := `UPDATE PROCESOS SET CANTIDAD_PLACAS = CASE WHEN CANTIDAD_PLACAS - @cantidadBuena < 0 THEN 0 ELSE CANTIDAD_PLACAS - @cantidadBuena END WHERE PLACA = @placa`
		res, err := db.Exec(updatePlacaQuery, sql.Named("cantidadBuena", cantidadBuena), sql.Named("placa", placa))
		if err != nil {
			continue
		}
		nRows, _ := res.RowsAffected()
		if nRows == 0 {
			continue
		}
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

	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	query := `UPDATE PROCESOS SET CANT_A_FABRICAR = CASE WHEN CANT_A_FABRICAR - @cantidad < 0 THEN 0 ELSE CANT_A_FABRICAR - @cantidad END, CANTIDAD_PRODUCIDA = ISNULL(CANTIDAD_PRODUCIDA,0) + @cantidad WHERE ID = @id`
	_, err = db.Exec(query, sql.Named("cantidad", req.Cantidad), sql.Named("id", req.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update process"})
		return
	}

	for i, placa := range req.Placas {
		if i >= len(req.PlacasBuenas) {
			continue
		}
		cantidadBuena := req.PlacasBuenas[i]
		updatePlacaQuery := `UPDATE PROCESOS SET CANTIDAD_PLACAS = CASE WHEN CANTIDAD_PLACAS - @cantidadBuena < 0 THEN 0 ELSE CANTIDAD_PLACAS - @cantidadBuena END WHERE PLACA = @placa`
		res, err := db.Exec(updatePlacaQuery, sql.Named("cantidadBuena", cantidadBuena), sql.Named("placa", placa))
		if err != nil {
			continue
		}
		nRows, _ := res.RowsAffected()
		if nRows == 0 {
			continue
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process updated"})
}
