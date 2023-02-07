package main

import (
	"faceit/db"
	"faceit/helpers"
	"faceit/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			helpers.HandlePanic(r)
		}
	}()

	appConf, err := helpers.GetConfig()
	helpers.ErrChan <- err

	helpers.ErrChan <- db.SetupDatabase(appConf.DB)
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger(), helpers.TooManyRequestMiddleware) // middleware.Recover() egyelőre debugolni jobb ha nincs

	webService, err := services.NewWebService(e, db.GetManager(), services.NewTestListener())
	helpers.ErrChan <- err

	helpers.ErrChan <- webService.Setup()

	// Start server
	helpers.ErrChan <- e.Start(fmt.Sprintf(":%d", appConf.AppPort))
}
