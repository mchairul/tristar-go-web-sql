package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"websql/handlers"
	"websql/middlewares"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/latihan_go")

	if err != nil {
		fmt.Println("Gagal koneksi ke DB : " + err.Error())
	}

	defer db.Close()

	http.Handle("/assets/",
		http.StripPrefix("/assets/",
			http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", handlers.HandleFormLogin())
	http.HandleFunc("/login", handlers.HandlePostLogin(db))
	http.HandleFunc("/logout", handlers.HandleLogout())

	http.HandleFunc("/listkaryawan", middlewares.SessionMiddleware(
		handlers.HandleListKaryawan(db)))
	http.HandleFunc("/tambahkaryawan", middlewares.SessionMiddleware(
		handlers.HandleTambahKaryawan(db)))
	http.HandleFunc("/postkaryawan", middlewares.SessionMiddleware(
		handlers.HandlePostTambahKaryawan(db)))
	http.HandleFunc("/editkaryawan", middlewares.SessionMiddleware(
		handlers.HandleEditKaryawan(db)))
	http.HandleFunc("/posteditkaryawan", middlewares.SessionMiddleware(
		handlers.HandlePostEditKaryawan(db)))
	http.HandleFunc("/deletekaryawan", middlewares.SessionMiddleware(
		handlers.HandleDeleteKaryawan(db)))

	fmt.Println("Server Started ...")
	http.ListenAndServe("localhost:8080", nil)
}
