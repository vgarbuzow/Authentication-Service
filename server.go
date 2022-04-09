package main

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

type Guid struct {
	Guid string `json:"guid"`
}

type MsgJson struct {
	Status int    `json:"status"`
	Msg    string `json:"message"`
}

func main() {
	initDB()
	r := mux.NewRouter()

	r.Handle("/api/get-token", GetTokensHandler).Methods("POST")
	r.Handle("/api/refresh-token", RefreshTokenHandler).Methods("GET")
	r.Handle("/api/check-token", CheckTokenHandler).Methods("POST")

	err := http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Fatal(err)
	}

}

var GetTokensHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Host", "localhost")
	body, err := io.ReadAll(r.Body)
	var guid Guid
	err = json.Unmarshal(body, &guid)
	token, err := readRefreshToken(guid.Guid)
	if token != nil {
		w.WriteHeader(http.StatusForbidden)
		message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Пользователь с указаным guid уже существует"})
		_, err = w.Write(message)
		return
	}
	access, err := GetNewAccessToken(guid.Guid)
	refresh, err := GetNewRefreshToken(guid.Guid)
	response, err := TokenEncodingJson(Tokens{Status: 1, Access: access, Refresh: refresh, Guid: guid.Guid})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Произошла ошибка на стороне сервера"})
		_, _ = w.Write(message)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
})

var CheckTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Host", "localhost")
	var token Tokens
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &token)
	_, err = AccessTokenParse(token.Access)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Токен не валидный!"})
		_, _ = w.Write(message)
		return
	}
	message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Валидация прошла успешно!"})
	_, err = w.Write(message)
	/*if access, err := r.Cookie("Access"); err != nil {
		_, err = w.Write([]byte("А токен то где, мужик?"))
	} else {
		result, err := ParseAccessToken(access.Value)
		_, err = w.Write([]byte(result))
		if err != nil {
			log.Fatal(err)
		}
	}*/

})

var RefreshTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	/*if access, err := r.Cookie("Access"); err == nil {
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
	}*/
})
