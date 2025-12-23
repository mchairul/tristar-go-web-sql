package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"websql/helpers"
	"websql/typecustom"

	"github.com/google/uuid"
)

func HandleListKaryawan(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.SetHeaders(w, r)
		rows, err := db.Query("SELECT * FROM karyawan")
		if err != nil {
			fmt.Println(err.Error())
		}

		defer rows.Close()

		var result []typecustom.Karyawan

		for rows.Next() {
			var each = typecustom.Karyawan{}
			var err = rows.Scan(&each.Id, &each.Nik, &each.Nama, &each.Alamat, &each.TglLahir, &each.Jk, &each.Foto)

			if err != nil {
				w.Write([]byte("error : " + err.Error()))
				return
			}

			result = append(result, each)
		}

		if err := rows.Err(); err != nil {
			w.Write([]byte("error: " + err.Error()))
			return
		}

		session, err := helpers.GetSessionStore(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		username, ok := session.Values["Username"]
		fmt.Println("userid", username)
		if !ok {
			http.Error(w, "gagal mendapatkan session Username", http.StatusInternalServerError)
			return
		}

		var data = typecustom.WebData{
			"Title":    "List Karyawan",
			"Karyawan": result,
			"Username": username,
		}

		var tmpl = template.Must(template.ParseFiles(
			"views/_header.html",
			"views/_footer.html",
			"views/_nav.html",
			"views/list_karyawan.html",
		))

		err = tmpl.ExecuteTemplate(w, "list_karyawan", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleTambahKaryawan(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.SetHeaders(w, r)
		var data = typecustom.WebData{
			"Title": "Tambah Karyawan",
		}

		var tmpl = template.Must(template.ParseFiles(
			"views/_header.html",
			"views/_footer.html",
			"views/_nav.html",
			"views/form_karyawan.html",
		))

		err := tmpl.ExecuteTemplate(w, "form_karyawan", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandlePostTambahKaryawan(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.SetHeaders(w, r)
		nik := r.FormValue("nik")
		nama := r.FormValue("nama")
		alamat := r.FormValue("alamat")
		tanggalLahir := r.FormValue("tanggal_lahir")
		jk := r.FormValue("jk")

		fileFoto, handlerFile, err := r.FormFile("foto")

		if err != nil {
			fmt.Println("error form file")
		}

		defer fileFoto.Close()

		fmt.Println("Size foto ", handlerFile.Size)

		fileName := handlerFile.Filename
		extension := filepath.Ext(fileName)
		uniqueFileName := uuid.New().String() + extension

		pathFile, err := os.Create("./assets/images/uploads/" + uniqueFileName)
		if err != nil {
			fmt.Println(err)
		}
		defer pathFile.Close()

		if _, err := io.Copy(pathFile, fileFoto); err != nil {
			fmt.Println("error upload")
		}

		strQuery := "INSERT INTO karyawan (nik, nama, alamat, tanggal_lahir, jenis_kelamin, foto) VALUES (?,?,?,?,?,?)"
		statementInsert, err := db.Prepare(strQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer statementInsert.Close()

		pathFoto := "assets/images/uploads/" + uniqueFileName
		_, err = statementInsert.Exec(nik, nama, alamat, tanggalLahir, jk, pathFoto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// redirec
		http.Redirect(w, r, "/listkaryawan", http.StatusSeeOther)
	}
}

func HandleEditKaryawan(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.SetHeaders(w, r)
		queryParam := r.URL.Query()

		id := queryParam.Get("id")

		if id == "" {
			http.Redirect(w, r, "/listuser", http.StatusMovedPermanently)
			return
		}

		var result = typecustom.Karyawan{}

		err := db.QueryRow("SELECT * FROM karyawan WHERE id = ?", id).
			Scan(&result.Id, &result.Nik, &result.Nama, &result.Alamat, &result.TglLahir, &result.Jk)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var data = typecustom.WebData{
			"Title": "Edit Karyawan",
			"Data":  result,
		}

		var tmpl = template.Must(template.ParseFiles(
			"views/_header.html",
			"views/_footer.html",
			"views/_nav.html",
			"views/form_edit_karyawan.html",
		))

		err = tmpl.ExecuteTemplate(w, "form_edit_karyawan", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandlePostEditKaryawan(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.SetHeaders(w, r)
		id := r.FormValue("id")
		nik := r.FormValue("nik")
		nama := r.FormValue("nama")
		alamat := r.FormValue("alamat")
		tanggalLahir := r.FormValue("tanggal_lahir")
		jk := r.FormValue("jk")

		strQuery := "UPDATE karyawan SET nik = ?, nama = ?, alamat = ?, "
		strQuery = strQuery + "tanggal_lahir = ?, jenis_kelamin = ? "
		strQuery = strQuery + "WHERE id = ?"
		statementUpdate, err := db.Prepare(strQuery)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer statementUpdate.Close()

		_, err = statementUpdate.Exec(nik, nama, alamat, tanggalLahir, jk, id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/listuser", http.StatusSeeOther)
	}
}

func HandleDeleteKaryawan(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.SetHeaders(w, r)
		queryParam := r.URL.Query()

		id := queryParam.Get("id")

		statemenDelete, err := db.Prepare("DELETE FROM karyawan WHERE id = ?")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer statemenDelete.Close()

		_, err = statemenDelete.Exec(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/listuser", http.StatusSeeOther)
	}
}
