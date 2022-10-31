package main

import (
	"checklist/internal/handlers"
	"checklist/internal/models"
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	dsn := "root:password@tcp(mysql_checklist:3306)/checklist"
	db, err := openDB(dsn)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	h := handlers.TaskHandler{TaskModel: models.TaskModel{DB: db}}
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handlers.SigningKey),
	}))
	e.GET("/:user_id", h.GetAll)
	e.POST("/", h.Create)
	e.PUT("/:id", h.Update)
	e.PATCH("/:id", h.Complete)
	e.DELETE(":id", h.Delete)
	e.Logger.Fatal(e.Start(":8000"))
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
