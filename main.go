package main

import (
	"database/sql"
	"example.com/go-rest/handlers"
	"example.com/go-rest/routes"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main(){
	dsn := "root:@tcp(127.0.0.1:3306)/go_notes_api?parseTime=true"

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	fmt.Println("Berhasil terhubung ke database!")

	noteHandler := &handlers.NoteHandler{DB: db}
	authHandler := &handlers.AuthHandler{DB:db}
}