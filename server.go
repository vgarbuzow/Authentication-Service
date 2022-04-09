package main

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

	r.Handle("/api/get-token", GetTokensHandler).Methods("GET")
	r.Handle("/api/refresh-token", RefreshTokenHandler).Methods("GET")
	r.Handle("/api", CheckTokenHandler).Methods("GET")
	r.Handle("/", http.FileServer(http.Dir("./views/")))

	err := http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Fatal(err)
	}

}

var GetTokensHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if guid := query.Get("guid"); guid != "" {
		access, refresh, err := CreateTokens(guid)
		cookieAccess, cookieRefresh := BuildCookiesTokens(access, refresh)
		http.SetCookie(w, &cookieAccess)
		http.SetCookie(w, &cookieRefresh)

		_, err = w.Write([]byte("Токены успешно созданы"))
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

var CheckTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if access, err := r.Cookie("Access"); err != nil {
		_, err = w.Write([]byte("А токен то где, мужик?"))
	} else {
		result, err := ParseAccessToken(access.Value)
		_, err = w.Write([]byte(result))
		if err != nil {
			log.Fatal(err)
		}
	}

})

var RefreshTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if access, err := r.Cookie("Access"); err == nil {
		if refresh, err := r.Cookie("Refresh"); err == nil {
			if isValid, guid := IsValidTokens(access.Value, refresh.Value); isValid {
				deleteRefreshToken(guid)
				access, refresh, err := CreateTokens(guid)
				cookieAccess, cookieRefresh := BuildCookiesTokens(access, refresh)
				http.SetCookie(w, &cookieAccess)
				http.SetCookie(w, &cookieRefresh)
				_, err = w.Write([]byte("Токены успешно обновлены!"))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	} else {
		_, err = w.Write([]byte("Токены не обновлены!"))
		if err != nil {
			log.Fatal(err)
		}
	}
})

func BuildCookiesTokens(access, refresh string) (http.Cookie, http.Cookie) {
	cookieAccess := http.Cookie{Name: "Access", Value: access, Expires: time.Now().Add(365 * 24 * time.Hour), MaxAge: 0,
		Secure: false, HttpOnly: true, Domain: "localhost", Path: "/"}
	cookieRefresh := http.Cookie{Name: "Refresh", Value: refresh, Expires: time.Now().Add(365 * 24 * time.Hour), MaxAge: 0,
		Secure: false, HttpOnly: true, Domain: "localhost", Path: "/api/refresh-token"}
	return cookieAccess, cookieRefresh
}
