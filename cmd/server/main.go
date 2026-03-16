package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/app"
	"github.com/dilvi/camp-booking-rest-api-go/internal/config"
)

func main() {
	cfg := config.Load()
	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	addr := ":" + cfg.AppPort
	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(addr, application.Router); err != nil {
		panic(err)
	}
}