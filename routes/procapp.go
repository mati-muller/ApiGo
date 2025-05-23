package routes

import (

	"net/http"


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
	r.GET("/app/pegado", getPegado)
	r.GET("/app/trozado", getTrozado)
	r.GET("/app/impresion", getImpresion)
	r.GET("/app/calado", getCalado)
	r.GET("/app/plizado", getPlizado)

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
func getPegado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN PEGADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getTrozado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN TROZADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getImpresion(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN IMPRESION p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getCalado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN CALADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getPlizado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN PLIZADO p2 ON p.ID = p2.ID
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}