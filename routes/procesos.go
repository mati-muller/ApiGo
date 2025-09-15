package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// Asegúrate de que la variable db esté disponible
var db *sql.DB // Debe ser inicializada en el main o en otro archivo y ser accesible aquí

func SetupProcesosRoutes(r *gin.Engine) {
	r.GET("/procesos/data", getData)
	r.GET("/procesos/pendientes-encolado", getPendientesEncolado)
	r.GET("/procesos/pendientes-encolado2", getPendientesEncolado)
	r.GET("/procesos/pendientes-emplacado", getPendientesEmplacado)
	r.GET("/procesos/pendientes-troquelado", getPendientesTroquelado)
	r.GET("/procesos/pendientes-troquelado2", getPendientesTroquelado)
	r.GET("/procesos/pendientes-calado", getPendientesCalado)
	r.GET("/procesos/pendientes-pegado", getPendientesPegado)
	r.GET("/procesos/pendientes-plizado", getPendientesPlizado)
	r.GET("/procesos/pendientes-trozado", getPendientesTrozado)
	r.GET("/procesos/pendientes-impresion", getPendientesImpresion)
	r.GET("/procesos/pendientes-multiple", getPendientesMultiple)
	r.GET("/procesos/pendientes-multiple2", getPendientesMultiple)
	r.GET("/procesos/pendientes-otro", getPendientesOtro)
	r.GET("/procesos/nv", getNV)
	r.GET("/procesosapp/encolado", getEncoladoProcesos)
}

// queryDatabase is imported from procapp.go

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
	log.Println("[DEBUG] Entrando a getPendientesEncolado")
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.NVCANT, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'ENCOLADO'
		ORDER BY p.FECHA_ENTREGA ASC
	`
	log.Println("[DEBUG] Ejecutando queryDatabase en getPendientesEncolado")
	queryDatabase(c, query)
	log.Println("[DEBUG] Saliendo de getPendientesEncolado")
}

func getPendientesEmplacado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.NVCANT, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'EMPLACADO'
		ORDER BY p.FECHA_ENTREGA ASC
	`
	queryDatabase(c, query)
}

func getPendientesTroquelado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.NVCANT, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'TROQUELADO'
		ORDER BY p.FECHA_ENTREGA ASC
	`
	queryDatabase(c, query)
}

func getPendientesCalado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.NVCANT, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'CALADO'
		ORDER BY p.FECHA_ENTREGA ASC
	`
	queryDatabase(c, query)
}

func getPendientesPegado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.NVCANT, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'PEGADO'
		ORDER BY p.FECHA_ENTREGA ASC
	`
	queryDatabase(c, query)
}

func getPendientesPlizado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.NVCANT, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'PLIZADO'
		ORDER BY p.FECHA_ENTREGA ASC
	`
	queryDatabase(c, query)
}

func getPendientesTrozado(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.NVCANT, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'TROZADO'
		ORDER BY p.FECHA_ENTREGA ASC
	`
	queryDatabase(c, query)
}

func getPendientesImpresion(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.NVCANT, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'IMPRESION'
		ORDER BY p.FECHA_ENTREGA ASC
	`
	queryDatabase(c, query)
}

func getPendientesMultiple(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.NVCANT, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.DesProd, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'MULTIPLE'
		ORDER BY p.FECHA_ENTREGA ASC
	`
	queryDatabase(c, query)
}

func getPendientesOtro(c *gin.Context) {
	query := `
		SELECT p.ID, p.NVNUMERO, p.NOMAUX, p.FECHA_ENTREGA, p.PROCESO, p.DETPROD, p.NVCANT, p2.CANT_A_PROD,
			   JSON_QUERY((
				   SELECT r.CodMat, r.CantMat
				   FROM REPORTES.dbo.recetas r
				   WHERE r.CodProd = p.CODPROD
				   FOR JSON PATH
			   )) AS Placas
		FROM REPORTES.dbo.procesos p
		JOIN REPORTES.dbo.procesos2 p2 ON p.ID = p2.ID
		WHERE p2.ESTADO_PROC = 'PENDIENTE' AND p.PROCESO = 'OTRO'
		ORDER BY p.FECHA_ENTREGA ASC
	`
	queryDatabase(c, query)
}
func getEncoladoProcesos(c *gin.Context) {
	query := `
		SELECT * FROM ENCOLADO
	`
	queryDatabase(c, query)
}

func getNV(c *gin.Context) {
	// Get query parameters for date filtering
	fechaInicio := c.Query("fecha_inicio")
	fechaFin := c.Query("fecha_fin")

	query := `
		SELECT 
			p.NVNUMERO,
			p.DetProd,
			p.NOMAUX,
			p.PROCESO, 
			p2.ESTADO_PROC, 
			p.NVCANT,
			p2.CANT_A_PROD,
			p.FECHA_ENTREGA,
			(p.NVCANT - p2.CANT_A_PROD) AS cantidad_producida
		FROM procesos p
		JOIN procesos2 p2 ON p.ID = p2.ID
		WHERE p.PROCESO != 'OTRO'`

	var args []interface{}

	if fechaInicio != "" {
		query += " AND p.FECHA_ENTREGA >= ?"
		args = append(args, fechaInicio)
	}

	if fechaFin != "" {
		query += " AND p.FECHA_ENTREGA <= ?"
		args = append(args, fechaFin)
	}

	// Add ordering by fecha_entrega (oldest to newest)
	query += " ORDER BY p.NVNUMERO ASC, p.DetProd ASC, p.FECHA_ENTREGA ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type Proceso struct {
		PROCESO           string `json:"PROCESO"`
		ESTADO_PROC       string `json:"ESTADO_PROC"`
		NVCANT            int    `json:"NVCANT"`
		CANT_A_PROD       int    `json:"CANT_A_PROD"`
		FECHA_ENTREGA     string `json:"FECHA_ENTREGA"`
		CantidadProducida int    `json:"cantidad_producida"`
	}

	type NV struct {
		NVNUMERO string    `json:"NVNUMERO"`
		DetProd  string    `json:"DetProd"`
		NOMAUX   string    `json:"NOMAUX"`
		Procesos []Proceso `json:"procesos"`
	}

	// Cambiar la clave para incluir tanto NVNUMERO como DetProd
	groupedData := map[string]*NV{}

	for rows.Next() {
		var nvnumero, detprod, nomaux, proceso, estadoProc, fechaEntrega string
		var cantprod, cantAProd, cantidadProducida int

		if err := rows.Scan(&nvnumero, &detprod, &nomaux, &proceso, &estadoProc, &cantprod, &cantAProd, &fechaEntrega, &cantidadProducida); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Crear clave única combinando NVNUMERO y DetProd
		key := nvnumero + "|" + detprod

		if _, exists := groupedData[key]; !exists {
			groupedData[key] = &NV{
				NVNUMERO: nvnumero,
				DetProd:  detprod,
				NOMAUX:   nomaux,
				Procesos: []Proceso{},
			}
		}

		groupedData[key].Procesos = append(groupedData[key].Procesos, Proceso{
			PROCESO:           proceso,
			ESTADO_PROC:       estadoProc,
			NVCANT:            cantprod,
			CANT_A_PROD:       cantAProd,
			FECHA_ENTREGA:     fechaEntrega,
			CantidadProducida: cantidadProducida,
		})
	}

	result := []NV{}
	for _, nv := range groupedData {
		result = append(result, *nv)
	}

	// Ordenar el resultado por NVNUMERO de menor a mayor, luego por DetProd
	slices.SortFunc(result, func(a, b NV) int {
		if a.NVNUMERO < b.NVNUMERO {
			return -1
		}
		if a.NVNUMERO > b.NVNUMERO {
			return 1
		}
		if a.DetProd < b.DetProd {
			return -1
		}
		if a.DetProd > b.DetProd {
			return 1
		}
		return 0
	})

	c.JSON(http.StatusOK, result)
}
