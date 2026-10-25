package main

import (
	"fmt"
	"net/http"

	"gopostgres/database"
	"gopostgres/handlers"
)

func main() {
	database.Connect()
	defer database.DB.Close()

	http.HandleFunc("/api/users", handlers.UsersHandler)

	fmt.Println("Сервер запущен на http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}