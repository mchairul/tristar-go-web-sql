package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"websql/handlers"

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

	http.HandleFunc("/listuser", handlers.HandleListKaryawan(db))
	http.HandleFunc("/tambahkaryawan", handlers.HandleTambahKaryawan(db))
	http.HandleFunc("/postkaryawan", handlers.HandlePostTambahKaryawan(db))
	http.HandleFunc("/editkaryawan", handlers.HandleEditKaryawan(db))
	http.HandleFunc("/posteditkaryawan", handlers.HandlePostEditKaryawan(db))
	http.HandleFunc("/deletekaryawan", handlers.HandleDeleteKaryawan(db))

	fmt.Println("Server Started ...")
	http.ListenAndServe("localhost:8080", nil)
}
