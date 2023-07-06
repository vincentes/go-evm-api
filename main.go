package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"web3pro.com/blockchain/internal"
	"web3pro.com/blockchain/internal/gas"
)

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal("Could not load env file.")
		return
	}
	internal.LoadEnvironmentVariables()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/gas/estimate", gas.Estimate)

	e.Logger.Fatal(e.Start(":1323"))
}
