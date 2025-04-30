package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupEdits(r *gin.Engine) {
	r.POST("/edit-troquelado", editTroquelado)
	r.POST("/edit-emplacado", editEmplacado)
	r.POST("/edit-trozado", editTrozado)
	r.POST("/edit-encolado", editEncolado)
}

func editTroquelado(c *gin.Context) {
	editTable(c, "TROQUELADO")
}

func editEmplacado(c *gin.Context) {
	editTable(c, "EMPLACADO")
}

func editTrozado(c *gin.Context) {
	editTable(c, "TROZADO")
}

func editEncolado(c *gin.Context) {
	editTable(c, "ENCOLADO")
}

func editTable(c *gin.Context, tableName string) {
	var requestBody struct {
		ID              int             `json:"ID"`
		CANT_A_FABRICAR int             `json:"CANT_A_FABRICAR"`
		PLACAS_A_USAR   json.RawMessage `json:"PLACAS_A_USAR"`
		CANTIDAD_PLACAS json.RawMessage `json:"CANTIDAD_PLACAS"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate JSON structure
	if !json.Valid(requestBody.PLACAS_A_USAR) || !json.Valid(requestBody.CANTIDAD_PLACAS) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON structure for PLACAS_A_USAR or CANTIDAD_PLACAS"})
		return
	}

	// Establish database connection
	db, err := sql.Open("sqlserver", "Server="+os.Getenv("SQL_SERVER")+"\\"+os.Getenv("SQL_INSTANCE")+";Database="+os.Getenv("SQL_DATABASE2")+";User Id="+os.Getenv("SQL_USER")+";Password="+os.Getenv("SQL_PASSWORD")+";Encrypt=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Execute the update query
	query := `
		UPDATE ` + tableName + `
		SET CANT_A_FABRICAR = @CANT_A_FABRICAR,
			PLACAS_A_USAR = @PLACAS_A_USAR,
			CANTIDAD_PLACAS = @CANTIDAD_PLACAS
		WHERE ID = @ID
	`
	_, err = db.Exec(query,
		sql.Named("ID", requestBody.ID),
		sql.Named("CANT_A_FABRICAR", requestBody.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", string(requestBody.PLACAS_A_USAR)),
		sql.Named("CANTIDAD_PLACAS", string(requestBody.CANTIDAD_PLACAS)),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute update query"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated " + tableName})
}
