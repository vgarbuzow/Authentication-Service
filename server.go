package main

// Импортируем необходимые зависимости. Мы будем использовать
// пакет из стандартной библиотеки и пакет от gorilla

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	r.Handle("/api/get-token", GetToken).Methods("GET")
	r.Handle("/api/refresh-token", NotImplemented).Methods("GET")
	r.Handle("/api", CheckToken).Methods("GET")
	r.Handle("/", http.FileServer(http.Dir("./views/")))

	err := http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		return
	}
}

var GetToken = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	token, err := createTokens()
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "Token", Value: token, Expires: expiration, MaxAge: 0, Secure: false, HttpOnly: true, Domain: "localhost", Path: "/"}
	http.SetCookie(w, &cookie)

	_, err = w.Write([]byte(token))
	if err != nil {
		log.Println(err)
	}
})

var CheckToken = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if token, err := r.Cookie("Token"); err != nil {
		_, err = w.Write([]byte("А токен то где, мужик?"))
	} else {
		result, err := parseAccessToken(token.Value)
		_, err = w.Write([]byte(result))
		if err != nil {
			log.Fatal(err)
		}
	}

})

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Not Implemented"))
	if err != nil {
		log.Println(err)
	}
})
