package main

import (
	"checklist/internal/handlers"
	"checklist/internal/models"
	"database/sql"
	"flag"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	port := flag.String("port", ":8000", "HTTP port")
	dsn := flag.String("dsn", "root:password@tcp(mysql_checklist:3306)/checklist", "MySQL data source name")
	jwtKey := flag.String("key", "secret", "jwt key")
	flag.Parse()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	db, err := openDB(*dsn)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	h := handlers.TaskHandler{TaskModel: models.TaskModel{DB: db}}
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(*jwtKey),
	}))
	e.GET("/:user_id", h.GetAll)
	e.POST("/", h.Create)
	e.PUT("/:id", h.Update)
	e.PATCH("/:id", h.Complete)
	e.DELETE("/:id", h.Delete)
	e.Logger.Fatal(e.Start(*port))
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
