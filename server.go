package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
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

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	infoLog.Printf("Соединение с БД")
	initDB()
	r := mux.NewRouter()
	r.Handle("/api/get-token", GetTokensHandler).Methods("POST")
	r.Handle("/api/refresh-token", RefreshTokenHandler).Methods("PUT")
	r.Handle("/api/check-token", CheckTokenHandler).Methods("POST")
	addr := flag.String("addr", ":4000", "localhost")
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  r,
	}
	infoLog.Printf("Запуск сервера на %s", *addr)
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

var GetTokensHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Host", "localhost")
	body, _ := io.ReadAll(r.Body)
	var guid Guid
	_ = json.Unmarshal(body, &guid)
	if token, err := ReadRefreshToken(guid.Guid); err != mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusInternalServerError)
		message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Произошла ошибка на стороне сервера"})
		_, _ = w.Write(message)
		return
	} else if token != nil {
		w.WriteHeader(http.StatusForbidden)
		message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Пользователь с указаным guid уже существует"})
		_, err = w.Write(message)
		return
	}
	WriteInResponseNewTokensJson(guid.Guid, &w)
})

var CheckTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Host", "localhost")
	var token Tokens
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &token)
	_, err = ParseVerifiedAccessToken(token.Access)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Токен не валидный!"})
		_, _ = w.Write(message)
		return
	}
	message, _ := json.Marshal(MsgJson{Status: 1, Msg: "Валидация прошла успешно!"})
	_, err = w.Write(message)
})

var RefreshTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Host", "localhost")
	body, err := io.ReadAll(r.Body)
	token, err := DecodingJsonToken(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Произошла ошибка на стороне сервера"})
		_, _ = w.Write(message)
		return
	}
	if claims, _ := ParseVerifiedAccessToken(token.Access); claims == nil {
		w.WriteHeader(http.StatusBadRequest)
		message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Произошла ошибка при валдицаии access токена"})
		_, _ = w.Write(message)
		return
	} else {
		if err := RefreshTokenValidate(claims.Guid, token.Refresh); err == nil {
			err = DeleteRefreshToken(claims.Guid)
			WriteInResponseNewTokensJson(claims.Guid, &w)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Произошла ошибка при валдицаии refresh токена"})
			_, _ = w.Write(message)
		}
	}
})

func WriteInResponseNewTokensJson(guid string, w *http.ResponseWriter) {
	writer := *w
	access, err := GetNewAccessToken(guid)
	refresh, err := CreateRefreshToken(guid, InsertRefreshToken)
	response, err := TokenEncodingJson(Tokens{Status: 1, Access: access, Refresh: refresh, Guid: guid})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		message, _ := json.Marshal(MsgJson{Status: 0, Msg: "Произошла ошибка на стороне сервера"})
		_, _ = writer.Write(message)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	_, err = writer.Write(response)
}
