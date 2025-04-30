package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupProcesosRoutes(r *gin.Engine) {
	r.GET("/data", getData)
	r.GET("/pendientes-encolado", getPendientesEncolado)
	r.GET("/pendientes-emplacado", getPendientesEmplacado)
	r.GET("/pendientes-troquelado", getPendientesTroquelado)
	r.GET("/pendientes-calado", getPendientesCalado)
	r.GET("/pendientes-pegado", getPendientesPegado)
	r.GET("/pendientes-plizado", getPendientesPlizado)
	r.GET("/pendientes-trozado", getPendientesTrozado)
	r.GET("/pendientes-impresion", getPendientesImpresion)
	r.GET("/pendientes-multiple", getPendientesMultiple)
	r.GET("/pendientes-otro", getPendientesOtro)
	r.GET("/nv", getNV)
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

func getData(c *gin.Context) {
	query := `
		SELECT * FROM (
			-- ...existing SQL query from /data route...
		) AS subc
		WHERE dif_fact > 0
		ORDER BY NVNUMERO ASC
	`
	queryDatabase(c, query)
}

func getPendientesEncolado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'ENCOLADO'
	`
	queryDatabase(c, query)
}

func getPendientesEmplacado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'EMPLACADO'
	`
	queryDatabase(c, query)
}

func getPendientesTroquelado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'TROQUELADO'
	`
	queryDatabase(c, query)
}

func getPendientesCalado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'CALADO'
	`
	queryDatabase(c, query)
}

func getPendientesPegado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'PEGADO'
	`
	queryDatabase(c, query)
}

func getPendientesPlizado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'PLIZADO'
	`
	queryDatabase(c, query)
}

func getPendientesTrozado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'TROZADO'
	`
	queryDatabase(c, query)
}

func getPendientesImpresion(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'IMPRESION'
	`
	queryDatabase(c, query)
}

func getPendientesMultiple(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'MULTIPLE'
	`
	queryDatabase(c, query)
}

func getPendientesOtro(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.CANTPROD, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.CodMat, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'OTRO'
	`
	queryDatabase(c, query)
}

func getNV(c *gin.Context) {
	query := `
		SELECT 
			p.NVNUMERO,
			p.DetProd,
			p.NOMAUX,
			p.PROCESO, 
			p2.ESTADO_PROC, 
			(p.CANTPROD - p2.CANT_A_PROD) AS cantidad_producida
		FROM procesos p
		JOIN procesos2 p2 ON p.ID = p2.ID
		WHERE p.PROCESO != 'OTRO'
	`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type Proceso struct {
		PROCESO           string `json:"PROCESO"`
		ESTADO_PROC       string `json:"ESTADO_PROC"`
		CantidadProducida int    `json:"cantidad_producida"`
	}

	type NV struct {
		NVNUMERO string    `json:"NVNUMERO"`
		DetProd  string    `json:"DetProd"`
		NOMAUX   string    `json:"NOMAUX"`
		Procesos []Proceso `json:"procesos"`
	}

	groupedData := map[string]*NV{}

	for rows.Next() {
		var nvnumero, detprod, nomaux, proceso, estadoProc string
		var cantidadProducida int

		if err := rows.Scan(&nvnumero, &detprod, &nomaux, &proceso, &estadoProc, &cantidadProducida); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if _, exists := groupedData[nvnumero]; !exists {
			groupedData[nvnumero] = &NV{
				NVNUMERO: nvnumero,
				DetProd:  detprod,
				NOMAUX:   nomaux,
				Procesos: []Proceso{},
			}
		}

		groupedData[nvnumero].Procesos = append(groupedData[nvnumero].Procesos, Proceso{
			PROCESO:           proceso,
			ESTADO_PROC:       estadoProc,
			CantidadProducida: cantidadProducida,
		})
	}

	result := []NV{}
	for _, nv := range groupedData {
		result = append(result, *nv)
	}

	c.JSON(http.StatusOK, result)
}
