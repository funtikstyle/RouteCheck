package service

import (
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func SendRequest(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	status := resp.StatusCode

	return status
}

func GetMD5SumString(f *os.File) (string, error) {
	file1Sum := md5.New()
	_, err := io.Copy(file1Sum, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", file1Sum.Sum(nil)), nil
}

func LogFailRoute(url string, status int) {
	file, err := os.OpenFile("routesOut.csv", os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		log.Println(err)
	}

	w := csv.NewWriter(file)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	timeNow := time.Now()

	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		timeNow.Year(), timeNow.Month(), timeNow.Day(),
		timeNow.Hour(), timeNow.Minute(), timeNow.Second())
	log.Println(status)

	row := []string{url, strconv.Itoa(status), formatted}
	err = w.Write(row)

	if err != nil {
		log.Println(err)
	}
	w.Flush()

}
