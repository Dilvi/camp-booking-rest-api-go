package main

import (
	"fmt"
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/app"
	"github.com/dilvi/camp-booking-rest-api-go/internal/config"
)

func main() {
	cfg := config.Load()
	router := app.NewRouter()
	addr := ":" + cfg.AppPort
	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}
}