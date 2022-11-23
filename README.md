# Test Task Golang

Запуск nats-streaming локально:
```
docker compose up -d
```

Запуск двух реплик в разных терминалах:
```
REPLICA_ID=1 NEXT_REPLICA_URL=http://localhost:8282 SETTER_PORT=8181 GETTER_PORT=8080 go run ./cmd/main.go 
```
```
REPLICA_ID=2 NEXT_REPLICA_URL=http://localhost:8080 SETTER_PORT=8383 GETTER_PORT=8282 go run ./cmd/main.go
```

В /config/config.go заданы некоторые данные для натс стриминга и также логин/пароль для Basic Auth для закрытия методы Set.
Логин и пароль для Basic Auth: test/test.

## API
### Метод SET
#### Запрос
`POST /set/`
```
curl --location --request POST 'http://localhost:8181/set' \
--header 'Authorization: Basic dGVzdDp0ZXN0' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id": "3333",
    "postpaid_limit": 1111,
    "spp": 23123,
    "shipping_fee": 333,
    "return_fee": 444
}'
```
#### Ответ
```
{
    "result": "User grade created successfully"
}
```
### Метод GET
#### Запрос
`GET /get/user_id=:id`
```
curl --location --request GET 'http://localhost:8080/get?user_id=3333'
```
#### Ответ
```
{
    "result": {
        "user_id": "3333",
        "postpaid_limit": 3333,
        "spp": 23123,
        "shipping_fee": 333,
        "return_fee": 444
    }
}
```
### Метод BACKUP
#### Запрос
`GET /backup`
```
curl --location --request GET 'http://localhost:8080/backup'
```
#### Ответ
```
Сжатый контент csv. Имя файла в заголовке Content-Disposition.
```

### Принцип работы синхронизации данных:
```
1. Отправляем POST запрос на http://localhost:8181/set с json, например таким:
{
    "user_id": "3333",
    "postpaid_limit": 3333,
    "spp": 23123,
    "shipping_fee": 333,
    "return_fee": 444
}
2. Проверяем, что на первой реплике сохранился наш объект http://localhost:8080/get?user_id=3333
3. Проверяем, что на второй реплике он тоже появился http://localhost:8282/get?user_id=3333, так как
данные ушли в натс канал, и вторая реплика их получила.
```
### Принцип заполнения данных с другой реплики:
```
1. Запускаем первую реплику REPLICA_ID=1 NEXT_REPLICA_URL=http://localhost:8282 SETTER_PORT=8181 GETTER_PORT=8080 go run ./cmd/main.go 
2. Отправляем POST запрос на http://localhost:8181/set с json, например таким:
{
    "user_id": "3333",
    "postpaid_limit": 3333,
    "spp": 23123,
    "shipping_fee": 333,
    "return_fee": 444
}
3. Запускаем вторую реплику REPLICA_ID=2 NEXT_REPLICA_URL=http://localhost:8080 SETTER_PORT=8383 GETTER_PORT=8282 go run ./cmd/main.go
4. Получаем со второй реплики данные http://localhost:8282/get?user_id=3333 и видим, что вторая реплика
дернула данные с первой реплики и заполнила ими сторадж.
```