# Сервис аутентификации

---

**Используемые технологии:**

- Go
- JWT
- MongoDB
- Postman

---

Три REST маршрута:

- Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
- Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов
- Третий маршрут проверят валидность Access токена

---

**Требования:**

Access токен: тип JWT, алгоритм SHA512, в базе не хранится, в теле присутствует идентификатор (GUID) .

Refresh токен: формат передачи base64, хранится в базе в виде bcrypt хеша.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним.

---

### Описание

db.go - представляет собой контроллер для обращения (записи, чтения и удаления) к базе данных MongoDB с помощью "go.mongodb.org/mongo-driver/mongo"

json_token.go - включает в себя структуру для представления токенов в формате json и их дальнейшей сериализации/десериализации

token.go - используется для верификации и генерации Refresh и Access токенов

server.go - представляет собой сервер (localhost:4000) с роутингом на /api/ с помощью "github.com/gorilla/mux",
ответ на запросы формируется в формате json

---

### Принцип работы сервиса
Для проверки работоспособности сервера в качесте клиента использовался Postman

1. Для получения Access, Refresh токенов необходимо сформировать запрос в формате json на адрес localhost:4000/api/get-token.
```yaml
{
    "guid":"6F9619FF-8B86-D011-B42D-00CF4FC964FF"
}  
```
В результате запроса будут сформированы токены, предварительно Refresh токен помещается в базу данных,
где id записи выступает guid. Ответ формируется в формате Json </br>

```yaml
{
"status": 1, </br>
"access":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.
eyJndWlkIjoiNkY5NjE5RkYtOEI4Ni1EMDExLUI0MkQtMDBDRjRGQzk2NEZGIiwiZXhwIjoxNjQ5NjE4Mzc2fQ.
CEoTyIHRJsXREF1SOGP-2UxqSoYGXwvF-Yj-6nMvUAo6LgIY9L4gArLQeIxdPvdgHzzIk8YoPo0MdxbMTWoCzw",
"refresh": "V1dXV1dXV1doQQ==",
"guid": "6F9619FF-8B86-D011-B42D-00CF4FC964FF"
} 
```
где status = 0..1 и в случае ошибки равен 0, в случае успешной генерации равен 1

2. Для обновления токенов необходимо сформировать запрос в формате json на адрес localhost:4000/api/refresh-token </br> 
```yaml
{
"access":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.
eyJndWlkIjoiNkY5NjE5RkYtOEI4Ni1EMDExLUI0MkQtMDBDRjRGQzk2NEZGIiwiZXhwIjoxNjQ5NjE4Mzc2fQ.
CEoTyIHRJsXREF1SOGP-2UxqSoYGXwvF-Yj-6nMvUAo6LgIY9L4gArLQeIxdPvdgHzzIk8YoPo0MdxbMTWoCzw",
"refresh": "V1dXV1dXV1doQQ=="
} 
```
В результате такого запроса будут получены токены в аналогичном формате п. 1.

3. Для проверки валидности Access токена необходимо сформировать запрос в формате json на адрес localhost:4000/api/check-token </br> 
```yaml
{
"access":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.
eyJndWlkIjoiNkY5NjE5RkYtOEI4Ni1EMDExLUI0MkQtMDBDRjRGQzk2NEZGIiwiZXhwIjoxNjQ5NjE4Mzc2fQ.
CEoTyIHRJsXREF1SOGP-2UxqSoYGXwvF-Yj-6nMvUAo6LgIY9L4gArLQeIxdPvdgHzzIk8YoPo0MdxbMTWoCzw"
} 
```
В результате запроса будет получен ответ, где status = 1 - токен валиден, status = 0 - ошибка при валидции
```yaml
{
"status": 1,
"message": "Валидация прошла успешно!"
}
```

В случае ошибок как со стороны сервиса, так и со стороны клиента (некорректный запрос), возвращается ошибка с message, в котором указана предварительная причина ошибки,
например при попытке валидировать некорректный токен, будет получена следующая ошибка:
```yaml
{
"status": 0,
"message": "Ошибка валидации access токена"
}
```
