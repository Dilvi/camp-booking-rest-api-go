package main

import (
	"fmt"
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/app"
)

func main() {
	router := app.NewRouter()
	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}