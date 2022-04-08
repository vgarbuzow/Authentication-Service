package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/getToken", getAccess)
	http.HandleFunc("/api/refreshToken", refreshToken)
	http.HandleFunc("/api/checkToken", checkToken) // each request calls handler
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func getAccess(w http.ResponseWriter, r *http.Request) {
	tkn, _ := createJWT()
	fmt.Fprintf(w, "Получение токена..."+tkn)

}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Обновление токена")
}

func checkToken(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Проверка токена")
}
