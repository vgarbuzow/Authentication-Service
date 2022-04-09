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
	initDB()
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	r.Handle("/api/get-token", GetToken).Methods("GET")
	r.Handle("/api/refresh-token", NotImplemented).Methods("GET")
	r.Handle("/api", CheckToken).Methods("GET")
	r.Handle("/", http.FileServer(http.Dir("./views/")))

	err := http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Fatal(err)
	}

}

var GetToken = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if guid := query.Get("guid"); guid != "" {
		token, refresh, err := createTokens(guid)
		cookieToken := http.Cookie{Name: "Token", Value: token, Expires: time.Now().Add(365 * 24 * time.Hour), MaxAge: 0,
			Secure: false, HttpOnly: true, Domain: "localhost", Path: "/"}
		cookieRefresh := http.Cookie{Name: "Refresh", Value: refresh, Expires: time.Now().Add(365 * 24 * time.Hour), MaxAge: 0,
			Secure: false, HttpOnly: true, Domain: "localhost", Path: "/api/refresh-token"}
		http.SetCookie(w, &cookieToken)
		http.SetCookie(w, &cookieRefresh)

		_, err = w.Write([]byte(token))
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := w.Write([]byte("Необходимо указать GUID в параметре запроса"))
		if err != nil {
			log.Fatal(err)
		}
	}
})

var CheckToken = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if token, err := r.Cookie("Refresh"); err != nil {
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
