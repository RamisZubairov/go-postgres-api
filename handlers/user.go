package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gopostgres/database"
	"gopostgres/models"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUser(w, r)
	case http.MethodGet:
		if r.URL.Query().Get("id") != "" {
			getOneUser(w, r)
		} else {
			getAllUsers(w, r)
		}
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u models.User 
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (name, id) VALUES ($1, $2)", u.Name, u.ID)
	if err != nil {
		fmt.Println("Ошибка БД:", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name FROM users")
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]models.User, 0) 
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name)
		if err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		list = append(list, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id must be a number", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		fmt.Println("Ошибка БД:", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getOneUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var u models.User
	err = database.DB.QueryRow("SELECT id, name FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name)

	if err == sql.ErrNoRows {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}