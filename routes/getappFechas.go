package routes

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func SetupProcAppFechaRoutes(r *gin.Engine) {
	r.GET("/app/troquelado/fecha", getTroqueladoFecha)
	r.GET("/app/troquelado2/fecha", getTroquelado2Fecha)
	r.GET("/app/encolado/fecha", getEncoladoFecha)
	r.GET("/app/encolado2/fecha", getEncolado2Fecha)
	r.GET("/app/multiple/fecha", getMultipleFecha)
	r.GET("/app/multiple2/fecha", getMultiple2Fecha)
	r.GET("/app/pegado/fecha", getPegadoFecha)
	r.GET("/app/trozado/fecha", getTrozadoFecha)
	r.GET("/app/impresion/fecha", getImpresionFecha)
	r.GET("/app/calado/fecha", getCaladoFecha)
	r.GET("/app/plizado/fecha", getPlizadoFecha)
	r.GET("/app/emplacado/fecha", getEmplacadoFecha)

}

func getTroqueladoFecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN TROQUELADO p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func getTroquelado2Fecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN TROQUELADO2 p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func getEncoladoFecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN REPORTES.dbo.ENCOLADO p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func getEncolado2Fecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN ENCOLADO2 p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func getMultipleFecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN MULTIPLE p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}

func getMultiple2Fecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN MULTIPLE2 p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getPegadoFecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN PEGADO p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getTrozadoFecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN TROZADO p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getImpresionFecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN IMPRESION p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getCaladoFecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN CALADO p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getPlizadoFecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN PLIZADO p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
func getEmplacadoFecha(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, 
		       p2.CANT_A_FABRICAR, p2.PLACAS_A_USAR, p2.CANTIDAD_PLACAS
		FROM procesos p
		JOIN EMPLACADO p2 ON p.ID = p2.ID
		WHERE CAST(p2.createdat AS DATE) = CAST(GETDATE() AS DATE)
		ORDER BY p2.PRIORITY
	`
	queryDatabase(c, query)
}
