package routes

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupEdits(r *gin.Engine) {
	r.POST("/edit-troquelado", editTroquelado)
	r.POST("/edit-troquelado2", editTroquelado)
	r.POST("/edit-emplacado", editEmplacado)
	r.POST("/edit-trozado", editTrozado)
	r.POST("/edit-encolado", editEncolado)
	r.POST("/edit-encolado2", editEncolado2)
	r.POST("/edit-multiple", editMultiple)
	r.POST("/edit-multiple2", editMultiple2)
	r.POST("/edit-pegado", editPegado)
	r.POST("/edit-impresion", editImpresion)
	r.POST("/edit-calado", editCalado)
	r.POST("/edit-plizado", editPlizado)
}

func editTroquelado(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE TROQUELADO SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update TROQUELADO"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TROQUELADO updated"})
}

func editEmplacado(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE EMPLACADO SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update EMPLACADO"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "EMPLACADO updated"})
}

func editTrozado(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE TROZADO SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update TROZADO"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TROZADO updated"})
}

func editEncolado(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE ENCOLADO SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ENCOLADO"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ENCOLADO updated"})
}
func editEncolado2(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE ENCOLADO2 SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ENCOLADO2"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ENCOLADO2 updated"})
}
func editPlizado(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE PLIZADO SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update PLIZADO"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PLIZADO updated"})
}
func editPegado(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE PEGADO SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update PEGADO"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PEGADO updated"})
}
func editImpresion(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE IMPRESION SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update IMPRESION"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "IMPRESION updated"})
}
func editCalado(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE CALADO SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update CALADO"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "CALADO updated"})
}
func editMultiple(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE MULTIPLE SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update MULTIPLE"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "MULTIPLE updated"})
}
func editMultiple2(c *gin.Context) {
	type request struct {
		ID                int      `json:"ID"`
		CANT_A_FABRICAR   int      `json:"CANT_A_FABRICAR"`
		TransformedPlacas []string `json:"transformedPlacas"`
		PlacasUsadas      []int    `json:"placasUsadas"`
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

	_, err = db.Exec(`UPDATE MULTIPLE2 SET CANT_A_FABRICAR = @CANT_A_FABRICAR, PLACAS_A_USAR = @PLACAS_A_USAR, CANTIDAD_PLACAS = @CANTIDAD_PLACAS WHERE ID = @ID`,
		sql.Named("CANT_A_FABRICAR", req.CANT_A_FABRICAR),
		sql.Named("PLACAS_A_USAR", toJSON(req.TransformedPlacas)),
		sql.Named("CANTIDAD_PLACAS", toJSON(req.PlacasUsadas)),
		sql.Named("ID", req.ID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update MULTIPLE2"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "MULTIPLE2 updated"})
}


// Utilidad para serializar a JSON (solo si no existe ya en el paquete)
// func toJSON(v interface{}) string {
// 	b, err := json.Marshal(v)
// 	if err != nil {
// 		return "[]"
// 	}
// 	return string(b)
// }
