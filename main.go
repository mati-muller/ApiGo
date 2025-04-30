package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var db *sql.DB

func init() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	// Configurar conexión a la base de datos
	dbServer := os.Getenv("SQL_SERVER")
	dbUser := os.Getenv("SQL_USER")
	dbPassword := os.Getenv("SQL_PASSWORD")
	dbName := os.Getenv("SQL_DATABASE2")

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", dbUser, dbPassword, dbServer, dbName)
	var err error
	db, err = sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatalf("Error abriendo la conexión a la base de datos: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error conectándose a la base de datos: %v", err)
	}

	log.Println("Conexión a la base de datos exitosa")
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Buena ctm",
		})
	})

	r.Run() // Por defecto en localhost:8080
}
