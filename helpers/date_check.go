package helpers

import "time"

func CheckTanggal(tanggal string) bool {
	format := "2006-01-02"

	_, err := time.Parse(format, tanggal)

	if err != nil {
		return false
	}
	return true
}
