package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Предупреждение: файл .env не найден")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		panic("Переменная DATABASE_URL не задана!")
	}

	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Успешно подключились к БД из пакета database!")
}